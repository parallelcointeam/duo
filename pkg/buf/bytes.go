package buf

import (
	"github.com/parallelcointeam/duo/pkg/proto"
)

// Bytes is a simple single byte
type Bytes proto.Bytes

// Bytes is a
func (r *Bytes) Bytes() *[]byte {
	panic("not implemented")
}

// Copy is a
func (r *Bytes) Copy(*[]byte) proto.Buffer {
	panic("not implemented")
}

// Zero is a
func (r *Bytes) Zero() proto.Buffer {
	panic("not implemented")
}

// Free is a
func (r *Bytes) Free() proto.Buffer {
	panic("not implemented")
}

// GetCoding is a
func (r *Bytes) GetCoding() string {
	panic("not implemented")
}

// SetCoding is a
func (r *Bytes) SetCoding(string) proto.Coder {
	panic("not implemented")
}

// ListCodings is a
func (r *Bytes) ListCodings() []string {
	panic("not implemented")
}

// Freeze is a
func (r *Bytes) Freeze() *[]byte {
	panic("not implemented")
}

// Thaw is a
func (r *Bytes) Thaw(*[]byte) interface{} {
	panic("not implemented")
}

// SetStatus is a
func (r *Bytes) SetStatus(string) proto.Status {
	panic("not implemented")
}

// SetStatusIf is a
func (r *Bytes) SetStatusIf(error) proto.Status {
	panic("not implemented")
}

// UnsetStatus is a
func (r *Bytes) UnsetStatus() proto.Status {
	panic("not implemented")
}

// SetElem is a
func (r *Bytes) SetElem(int, interface{}) proto.Array {
	panic("not implemented")
}

// GetElem is a
func (r *Bytes) GetElem(int) interface{} {
	panic("not implemented")
}

// Len is a
func (r *Bytes) Len() int {
	panic("not implemented")
}
