package wallet

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/anaskhan96/base58check"
	"github.com/awnumar/memguard"
	"github.com/mitchellh/go-homedir"
	"time"
)

// Stores an address book entry in a wallet.dat
type BName struct {
	Addr []byte
	Name []byte
}

// Stores key metadata in a wallet.dat
type BMetadata struct {
	Pub        []byte
	Version    uint32
	CreateTime time.Time
}

// Stores unencrypted keys in a wallet.dat
type BKey struct {
	Pub  []byte
	Priv []byte
}

// An unencrypted key pair with extra metadata for managing expiry in a wallet.dat
type BWKey struct {
	Pub         []byte
	Priv        []byte
	TimeCreated time.Time
	TimeExpires time.Time
	Comment     string
}

// A key pair with plaintext public key and AES-256-CBC encrypted private key
type BCKey struct {
	Pub  []byte
	Priv []byte
}

// Stores the default key that will appear in a wallet interface when creating a payment request
type BDefaultKey struct {
	Key []byte
}

// A collection of tables from a wallet.dat file with optional en/decryptors
type Imports struct {
	*Serializable
	Names      []BName
	Metadata   []BMetadata
	Keys       []BKey
	WKeys      []BWKey
	CKeys      []BCKey
	DefaultKey BDefaultKey
}
type imports interface {
	ToEncryptedStore() (es EncryptedStore)
	EncryptData(dst *memguard.LockedBuffer, src []byte)
}

var ES *EncryptedStore

// Import reads an existing wallet.dat and returns all the keys and address data in it. If a password is given, the private keys in the CKeys array are decrypted and the encrypter/decrypter functions are armed.
func Import(pass *memguard.LockedBuffer, filename ...string) (es *EncryptedStore, err error) {
	es = NewEncryptedStore()
	ES = es
	if pass == nil {
		err = errors.New("for reasons of security, it is not possible to have an unencrypted wallet")
		return
	}
	var db BDB
	var cursor bdb.Cursor
	if len(filename) == 0 {
		home, _ := homedir.Dir()
		db.SetFilename(home + "/.parallelcoin/wallet.dat")
	} else {
		db.SetFilename(filename[0])
	}
	if err = db.Open(); err != nil {
		return
	} else if cursor, err = db.Cursor(bdb.NoTransaction); err != nil {
		return
	} else {
		rec := [2][]byte{}
		if err = cursor.First(&rec); err != nil {
			return
		} else {
			for {
				idLen := rec[0][0] + 1
				rec[0] = []byte(string(rec[0]))
				rec[1] = []byte(string(rec[1]))
				id := string(rec[0][1:idLen])
				switch id {
				case "mkey":
					I := len(es.masterKey)
					es.masterKey = append(es.masterKey, new(MasterKey))
					es.masterKey[I].MKeyID = int64(binary.LittleEndian.Uint32(rec[0][idLen : idLen+4]))
					ekLen := rec[1][0] + 1
					es.masterKey[I].EncryptedKey = rec[1][1:ekLen]
					sLen := rec[1][ekLen]
					es.masterKey[I].Salt = rec[1][ekLen+1 : sLen+ekLen+1]
					es.masterKey[I].Method = binary.LittleEndian.Uint32(rec[1][sLen+ekLen+1 : sLen+ekLen+5])
					es.masterKey[I].Iterations = binary.LittleEndian.Uint32(rec[1][sLen+ekLen+5 : sLen+ekLen+9])
					es.masterKey[I].Other = rec[1][sLen+ekLen+9:]
					es.MasterKey = append(es.MasterKey, new(MasterKey))
					es.MasterKey[I].MKeyID = es.masterKey[I].MKeyID
					es.MasterKey[I].EncryptedKey = es.masterKey[I].EncryptedKey
					es.MasterKey[I].Salt = es.masterKey[I].Salt
					es.MasterKey[I].Method = es.masterKey[I].Method
					es.MasterKey[I].Iterations = es.masterKey[I].Iterations
					es.MasterKey[I].Other = es.masterKey[I].Other
				}
				if err = cursor.Next(&rec); err != nil {
					err = cursor.Close()
					break
				}
			}
			if pass != nil {
				es.armed, es.ckey, es.iv, _ = es.DeriveCipher(pass)
				es.LastLocked = time.Now()
			}
			if err = cursor.First(&rec); err != nil {
				return
			} else {
				for {
					idLen := rec[0][0] + 1
					rec[0] = []byte(string(rec[0]))
					rec[1] = []byte(string(rec[1]))
					id := string(rec[0][1:idLen])
					switch id {
					case "name":
						nameLen := rec[1][0] + 1
						if nameLen == 1 {
							break
						} else {
							pubS := rec[0][idLen+1 : idLen+1+rec[0][idLen]]
							pubH, _ := base58check.Decode(string(pubS))
							pubB, _ := hex.DecodeString(pubH)
							pub, _ := NewBufferFromBytes(pubB)
							label, _ := NewBufferFromBytes(rec[1][1:nameLen])
							I := len(es.AddressBook)
							es.AddressBook = append(es.AddressBook, NewAddressBook(es.Serializable))
							var r []*memguard.LockedBuffer
							r, _ = es.armed.Encrypt(es.ckey, es.iv, pub, label)
							es.AddressBook[I].Pub = append([]byte{}, r[0].Buffer()...)
							es.AddressBook[I].Label = append([]byte{}, r[1].Buffer()...)
							DeleteBuffers(r...)
						}
					case "key":
						pub, _ := NewBufferFromBytes(rec[0][idLen+1 : idLen+1+rec[0][idLen]])
						priv, _ := NewBufferFromBytes(rec[1][1 : 1+rec[1][0]])
						I := len(es.Key)
						es.Key = append(es.Key, NewKey(es.Serializable))
						var r []*memguard.LockedBuffer
						r, _ = es.armed.Encrypt(es.ckey, es.iv, pub, priv)
						es.Key[I].Pub = append([]byte{}, r[0].Buffer()...)
						es.Key[I].Priv = append([]byte{}, r[1].Buffer()...)
						DeleteBuffers(r...)
					case "wkey":
						pub, _ := NewBufferFromBytes(rec[0][idLen+1 : idLen+1+rec[0][idLen]])
						pLen := rec[1][0] + 1
						priv, _ := NewBufferFromBytes(rec[1][1:pLen])
						Len := len(es.Key)
						es.Key = append(es.Key, NewKey(es.Serializable))
						var r []*memguard.LockedBuffer
						if es.armed != nil {
							r, _ = es.armed.Encrypt(es.ckey, es.iv, pub, priv)
						} else {
							r = []*memguard.LockedBuffer{pub, priv}
						}
						DeleteBuffers(pub, priv)
						es.Key[Len].Pub = append([]byte{}, r[0].Buffer()...)
						es.Key[Len].Priv = append([]byte{}, r[1].Buffer()...)
						DeleteBuffers(r...)
						tc, _ := NewBufferFromBytes(rec[1][pLen : pLen+8])
						te, _ := NewBufferFromBytes(rec[1][pLen+8 : pLen+16])
						cLen := rec[1][pLen+16]
						comment, _ := NewBufferFromBytes(rec[1][pLen+16 : pLen+cLen+16])
						r, _ = es.armed.Encrypt(es.ckey, es.iv, tc, te, comment)
						DeleteBuffers(r...)
						I := len(es.Wdata)
						es.Wdata = append(es.Wdata, NewWdata(es.Serializable))
						es.Wdata[I].Pub = es.Key[Len].Pub
						es.Wdata[I].Created = append([]byte{}, r[0].Buffer()...)
						es.Wdata[I].Expires = append([]byte{}, r[1].Buffer()...)
						es.Wdata[I].Comment = append([]byte{}, r[2].Buffer()...)
						DeleteBuffers(r...)
					case "ckey":
						pub, _ := NewBufferFromBytes(rec[0][idLen+1 : idLen+1+rec[0][idLen]])
						privLen := rec[1][0] + 1
						I := len(es.Key)
						es.Key = append(es.Key, NewKey(es.Serializable))
						var r []*memguard.LockedBuffer
						r, _ = es.armed.Encrypt(es.ckey, es.iv, pub)
						es.Key[I].Pub = append([]byte{}, r[0].Buffer()...)
						es.Key[I].Priv = rec[1][1:privLen]
						DeleteBuffers(r...)
					case "keymeta":
						pubLen := rec[0][idLen]
						pub := rec[0][idLen+1 : pubLen+idLen+1]
						I := len(es.Metadata)
						es.Metadata = append(es.Metadata, NewMetadata(es.Serializable))
						p, _ := NewBufferFromBytes(pub)
						ct, _ := NewBufferFromBytes(rec[1][4:12])
						var r []*memguard.LockedBuffer
						r, _ = es.armed.Encrypt(es.ckey, es.iv, p, ct)
						es.Metadata[I].Pub = append([]byte{}, r[0].Buffer()...)
						es.Metadata[I].Version = binary.LittleEndian.Uint32(rec[1][:4])
						es.Metadata[I].CreateTime = append([]byte{}, r[1].Buffer()...)
						DeleteBuffers(r...)
					case "defaultkey":
						l := rec[1][0] + 1
						k, _ := NewBufferFromBytes(rec[1][1:l])
						var r []*memguard.LockedBuffer
						r, _ = es.armed.Encrypt(es.ckey, es.iv, k)
						es.DefaultKey = append([]byte{}, r[0].Buffer()...)
						DeleteBuffers(r...)
					case "pool":
						id := binary.LittleEndian.Uint64(rec[0][5:13])
						version := binary.LittleEndian.Uint32(rec[1][:4])
						timestamp := rec[1][4:12]
						pub := rec[1][13:46]
						I := len(es.Pool)
						es.Pool = append(es.Pool, NewPool(es.Serializable))
						es.Pool[I].Index = id
						es.Pool[I].Version = version
						T, _ := NewBufferFromBytes(timestamp)
						P, _ := NewBufferFromBytes(pub)
						var r []*memguard.LockedBuffer
						r, _ = es.armed.Encrypt(es.ckey, es.iv, T, P)
						es.Pool[I].Time = append([]byte{}, r[0].Buffer()...)
						es.Pool[I].Pub = append([]byte{}, r[1].Buffer()...)
						DeleteBuffers(r...)
						// abc = AllocatedBufferCount
					}
					if err = cursor.Next(&rec); err != nil {
						err = cursor.Close()
						break
					}
				}
			}
		}
	}
	return
}
