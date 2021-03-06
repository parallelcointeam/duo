package key

import (
	"github.com/parallelcointeam/duo/pkg/bc"
	"github.com/parallelcointeam/duo/pkg/core"
)

// NewStore creates a new Store
func NewStore() *Store {
	r := new(Store)
	r.privs = make(map[core.Address]*Priv)
	r.pubs = make(map[core.Address]*Pub)
	return r
}

// AddPriv adds a new private key to the store
func (r *Store) AddPriv(priv *Priv) *Store {
	switch {
	case r == nil:
		r = NewStore()
		r.SetStatus(er.NilRec)
	case priv == nil:
		r.SetStatus(er.NilParam)
	default:
		id := priv.GetID()
		if _, ok := r.privs[id]; ok {
			r.SetStatus("priv already in store")
		} else {
			r.privs[id] = priv
			r.UnsetStatus()
		}
	}
	return r
}

// AddPub adds a new public key to the store
func (r *Store) AddPub(pub *Pub) *Store {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
	}
	id := NewID(pub.Bytes())
	if _, ok := r.privs[id]; ok {
		r.SetStatus("pub already in with priv")
	} else if _, ok := r.pubs[id]; ok {
		r.SetStatus("pub already in store")
	} else {
		r.pubs[id] = pub
		r.UnsetStatus()
	}
	return r
}

// Remove a key from the store by ID (address)
func (r *Store) Remove(id core.Address) *Store {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
		return r
	}
	if _, ok := r.privs[id]; ok {
		delete(r.privs, id)
		return r.UnsetStatus().(*Store)
	} else if _, ok := r.pubs[id]; ok {
		delete(r.pubs, id)
		return r.UnsetStatus().(*Store)
	}
	r.SetStatus("id not found")
	return r
}

// Encrypt sets the store to encrypt private keys
func (r *Store) Encrypt(blockCrypt *bc.BlockCrypt) *Store {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
	} else {
		for i := range r.privs {
			r.privs[i].BC = blockCrypt
			r.privs[i].Copy(r.privs[i].Bytes())
		}
	}
	return r
}

// Decrypt sets the store to not encrypt private keys
func (r *Store) Decrypt() *Store {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
	} else {
		for i := range r.privs {
			tmp := r.privs[i].Bytes()
			r.privs[i].BC = nil
			r.privs[i].Copy(tmp)
		}
		r.UnsetStatus()
	}
	return r
}

// Find returns the key with matching ID as requested. The return type is Priv but if there is no private key the field is empty
func (r *Store) Find(id core.Address) (out *Priv) {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
		return &Priv{}
	}
	out = new(Priv)
	I := core.Address(id)
	if _, ok := r.privs[I]; ok {
		return r.privs[I]
	}
	if _, ok := r.pubs[I]; ok {
		out.pub = r.pubs[I]
		return
	}
	return
}

// SetStatus is a
func (r *Store) SetStatus(s string) core.Status {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
	} else {
		r.Status = s
	}
	return r
}

// SetStatusIf is a
func (r *Store) SetStatusIf(err error) core.Status {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
	} else {
		if err != nil {
			r.Status = err.Error()
		}
	}
	return r
}

// UnsetStatus is a
func (r *Store) UnsetStatus() core.Status {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
	} else {
		r.Status = ""
	}
	return r
}

// OK returns true if there is no error
func (r *Store) OK() bool {
	if r == nil {
		r = NewStore()
		r.SetStatus(er.NilRec)
		return false
	}
	return r.Status == ""
}
