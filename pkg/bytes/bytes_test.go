package bytes

import (
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	a := new(Bytes)
	A := []byte("test")
	a.Load(&A)
	fmt.Println("a", a.Buf())
	b := new(Bytes)
	b.Copy(a)
	fmt.Println("copy a to b", b.Buf())
	fmt.Println("before move a", *a, "b", *b)
	b.Move(a)
	fmt.Println("after move a", *a, "b", *b)
	a.Link(b)
	fmt.Println("link emptied b to a", a.Buf(), *b.Buf())
	c := a.buf
	(*c)[0] = 1
	fmt.Println("now both the same memory (changed byte zero of first only)", a.Buf(), b.Buf())
	fmt.Println("Struct literal with Rand", struct{ *Bytes }{}.Rand(32).Buf())
	fmt.Println("Struct literal with Null", struct{ *Bytes }{}.Null().Buf())
	fmt.Println("Struct literal with Len()", struct{ *Bytes }{}.Size())
	fmt.Println("Struct literal with Null().Len()", struct{ *Bytes }{}.Null().Size())
	fmt.Println("Struct literal with Null().New(32)", struct{ *Bytes }{}.Null().New(32))
	var d *Bytes
	fmt.Println("nil pointer with Buf()", d.Buf())
	d = nil
	fmt.Println("nil pointer with Load()", d.Load(&A).Buf())
	d = nil
	fmt.Println("nil pointer with Copy()", d.Copy(a).Buf())
	d = nil
	fmt.Println("nil pointer with Copy(empty)", d.Copy(&Bytes{nil, false, 0, nil}))
	fmt.Println("nil pointer with Copy(Buf zero len)", d.Copy(&Bytes{&[]byte{}, false, 0, nil}))
	fmt.Println("Struct pointer with Copy(<nil>)", a.Load(&A).Copy(nil))
	d = nil
	A = []byte("this is longer")
	fmt.Println(A)
	a.Load(&A)
	fmt.Println(a.Buf())
	fmt.Println("nil pointer with Link()", a.Buf(), d.Link(a).Buf())
	f := NewBytes().Rand(13)
	fmt.Println("NewBytes().Rand(13)", f, f.Buf())
	fmt.Println("NewBytes().Move(NewBytes().New(13)).Error()", NewBytes().Move(NewBytes().New(13)).Error())
	d = nil
	fmt.Println("nil pointer with Move(empty)", d.Move(&Bytes{nil, false, 0, nil}))
	d = nil
	fmt.Println("nil pointer with Error()", d.Error())
	d = nil
	fmt.Println("nil pointer with Error().SetError()", d.SetError("testing").Error())
	fmt.Println("copying self", f.Copy(f))
	fmt.Println("nil IsSet()", d.IsSet())
	fmt.Println("non nil IsSet()", f.IsSet())
	fmt.Println("nil Load(nil)", d.Load(nil))
	fmt.Println("nil Move(nil)", d.Move(nil))
	fmt.Println("JSON UTF8", f.Load(&A).SetCoding("utf8").Coding())
	B := []byte("this is longer    ")
	fmt.Println("JSON hex", f.Load(&B).SetCoding("binary").Coding())
	fmt.Println("JSON nil val", f.Load(nil).String())
	f.Elem(0)
	f.SetElem(0, NewBytes().Load(&[]byte{100}))
	d.Purge()
	var x *Bytes
	x.SetCoding("binary").UnsetError()
	x.Elem(0)
	x.SetElem(0, NewBytes().Load(&[]byte{100}))
	x.Load(&B).SetElem(0, NewBytes().Load(&[]byte{100}))
	x.Load(&B).Elem(0)
	var y *Bytes
	y.UnsetError()
}
