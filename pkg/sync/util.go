package sync

import (
	"encoding/binary"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/golang/snappy"
	"github.com/parallelcointeam/duo/pkg/core"
)

func removeTrailingZeroes(in []byte) []byte {
	length := 0
	for i := len(in) - 1; i >= 0; i-- {
		if in[i] != 0 {
			length = i
			break
		}
	}
	return in[:length+1]
}

func removeLeadingZeroes(in []byte) []byte {
	nonzerostart := 0
	for i := range in {
		if in[i] != 0 {
			nonzerostart = i
			break
		}
	}
	return in[nonzerostart:]
}

func (r *Node) getLatest() (h uint32) {
	var latestB []byte

	r.SetStatusIf(r.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("latest"))
		if err == nil {
			latestB, err = item.Value()
		}
		return err
	}))
	if latestB != nil {
		heightB := latestB[:4]
		core.BytesToInt(&h, &heightB)
	}
	r.UnsetStatus()
	return
}

func decodeAddressRecord(addr []byte) (out []Location, length uint32) {
	var dec []byte
	var err error
	dec, err = snappy.Decode(nil, addr)
	if err == nil {
		addr = dec
	}
	cursor, length := 0, 0
	for cursor < len(addr) {
		var h, n uint64
		var height uint32
		var txnum uint16
		h, step := binary.Uvarint(addr[cursor:])
		cursor += step
		n, step = binary.Uvarint(addr[cursor:])
		cursor += step
		l := len(out)
		if l > 1 {
			height = out[l-1].Height + uint32(h)
		} else {
			height, txnum = uint32(h), uint16(n)
		}
		out = append(out, Location{Height: uint32(height), TxNum: uint16(txnum)})
		length++
	}
	return
}

func encodeAddressRecord(existing []byte, loc Location) (out []byte) {
	ex, _ := decodeAddressRecord(existing)
	// fmt.Println("ex", ex)
	ex = append(ex, loc)
	var h, n []byte
	for i := range ex {
		if i == 0 {
			h = make([]byte, 5)
			l := binary.PutUvarint(h, uint64(ex[0].Height))
			h = h[:l]
			n = make([]byte, 5)
			l = binary.PutUvarint(n, uint64(ex[0].TxNum))
			n = n[:l]
			out = append(h, n...)
		} else {
			nh := ex[i].Height - ex[i-1].Height
			// fmt.Println("\n", ex[i].Height, ex[i-1].Height, nh)
			h = make([]byte, 5)
			l := binary.PutUvarint(h, uint64(nh))
			h = h[:l]
			n = make([]byte, 5)
			l = binary.PutUvarint(n, uint64(ex[i].TxNum))
			n = n[:l]
			out = append(out, append(h, n...)...)
		}
		// fmt.Println("\n", uint64(ex[i].Height), h, uint64(ex[0].TxNum), n)
	}
	enc := snappy.Encode(nil, out)
	out = enc
	if len(enc) > len(out) {
		fmt.Println("\ncompressor expanded!")
	}
	// fmt.Print(len(out)) //, hex.EncodeToString(out))
	return
}

// AppendVarint takes any type of integer and returns the Varint. This is 7 bits per byte with a continuation marker on the 8th bit
func AppendVarint(to []byte, in interface{}) (out []byte) {
	var U uint64
	var I int64
	var signed bool
	switch in.(type) {
	case uint:
		U = uint64(in.(uint))
	case byte:
		U = uint64(in.(byte))
	case uint16:
		U = uint64(in.(uint16))
	case uint32:
		U = uint64(in.(uint32))
	case uint64:
		U = uint64(in.(uint64))
	case int:
		signed = true
		I = int64(in.(int))
	case int8:
		signed = true
		I = int64(in.(int8))
	case int16:
		signed = true
		I = int64(in.(int16))
	case int32:
		signed = true
		I = int64(in.(int32))
	case int64:
		signed = true
		I = int64(in.(int64))
	default:
		return []byte{}
	}
	out = make([]byte, 9)
	if signed {
		l := binary.PutVarint(out, I)
		out = out[:l]
		return append(to, out...)
	}
	l := binary.PutUvarint(out, U)
	out = out[:l]
	return append(to, out...)

}

// ExtractVarint reads the first varint contained in a given byte slice and returns the value according to the type of the typ parameter, and slices the input bytes removing the extracted integer
func ExtractVarint(typ interface{}, in []byte) (outbytes []byte, outint interface{}) {
	switch typ.(type) {
	case uint:
		o, l := binary.Uvarint(in)
		outbytes = in[l:]
		outint = uint(o)
	case byte:
		o, l := binary.Uvarint(in)
		outbytes = in[l:]
		outint = byte(o)
	case uint16:
		o, l := binary.Uvarint(in)
		outbytes = in[l:]
		outint = uint16(o)
	case uint32:
		o, l := binary.Uvarint(in)
		outbytes = in[l:]
		outint = uint32(o)
	case uint64:
		o, l := binary.Uvarint(in)
		outbytes = in[l:]
		outint = uint64(o)
	case int:
		o, l := binary.Varint(in)
		outbytes = in[l:]
		outint = int(o)
	case int8:
		o, l := binary.Varint(in)
		outbytes = in[l:]
		outint = int8(o)
	case int16:
		o, l := binary.Varint(in)
		outbytes = in[l:]
		outint = int16(o)
	case int32:
		o, l := binary.Varint(in)
		outbytes = in[l:]
		outint = int32(o)
	case int64:
		o, l := binary.Varint(in)
		outbytes = in[l:]
		outint = int64(o)
	default:
		return nil, []byte{}
	}
	return
}
