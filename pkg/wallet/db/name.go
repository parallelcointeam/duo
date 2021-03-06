package db

import (
	"encoding/hex"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// ReadName reads a name entry out of the database
func (r *DB) ReadName(id *[]byte) (out *rec.Name) {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	out = new(rec.Name)
	k := []byte(rec.Tables["Name"])
	idx := core.Hash64(id)
	k = append(k, *idx...)
	if r.BC != nil {
		id = r.BC.Encrypt(id)
	}
	k = append(k, *id...)
	fmt.Println("search", hex.EncodeToString(k))
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	var V []byte
	var meta byte
	err := r.DB.View(func(txn *badger.Txn) error {
		item, er := txn.Get(k)
		if er != nil {
			return er
		}
		meta = item.UserMeta()
		V, er = item.Value()
		if er != nil {
			return er
		}
		return nil
	})
	if r.SetStatusIf(err).OK() {
		switch {
		case r.BC != nil && meta&1 == 1:
			out.Address = *r.BC.Decrypt(id)
			out.Label = string(*r.BC.Decrypt(&V))
		case meta&1 != 1:
			out.Address = *id
			out.Label = string(V)
		default:
			r.SetStatus("record marked encrypted but no BC to decrypt with")
			fmt.Println(r.Error())
			out = new(rec.Name)
		}
		out.Idx = *idx
	}
	return
}

// WriteName writes a name entry to the database
func (r *DB) WriteName(address, label *[]byte) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	if address == nil || label == nil {
		r.SetStatus(er.NilParam)
	}
	idx := core.Hash64(address)
	var meta byte
	if r.BC != nil {
		meta = 1
		address = r.BC.Encrypt(address)
		label = r.BC.Encrypt(label)
	}
	k := []byte(rec.Tables["Name"])
	k = append(k, *idx...)
	k = append(k, *address...)
	v := *label
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.SetWithMeta(k, v, meta)
	}))
	return r
}

// EraseName removes a name entry from the database
func (r *DB) EraseName(address *[]byte) *DB {
	r = r.NewIf()
	if !r.OK() {
		return nil
	}
	opt := badger.DefaultIteratorOptions
	opt.PrefetchValues = false
	idx := core.Hash64(address)
	if r.BC != nil {
		address = r.BC.Encrypt(address)
	}
	k := []byte(rec.Tables["Name"])
	k = append(k, *idx...)
	k = append(k, *address...)
	r.SetStatusIf(r.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(k)
	}))
	return r
}
