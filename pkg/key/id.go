package key

import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
)

// ID is the RIPEMD160 hash of a key
type ID struct {
	Uint.U160
}

// Get the value of an ID
func (i *ID) Get() Uint.U160 {
	return i.U160
}

// Set the value of an ID
func (i *ID) Set(d Uint.U160) {
	i.U160 = d
}
