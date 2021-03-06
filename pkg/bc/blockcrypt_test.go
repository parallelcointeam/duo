package bc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"testing"
)

func TestGenerate(t *testing.T) {
	p := []byte("testingpassword123!")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	nc := New()
	nc.Ciphertext = nil
	nc.Generate(pass)
	// nc.Generate(pass)
	nc.Generate(nil)
}

func TestBlockCrypt(t *testing.T) {
	bc := New()
	p := []byte("testingpassword123!")
	pass := buf.NewSecure()
	pass.Copy(&p)
	bc = bc.Generate(pass)
	fmt.Println(bc.Crypt.Bytes())
	fmt.Println(bc.Ciphertext.Bytes())
	fmt.Println(bc.Arm())
	fmt.Println(p)
	pp := bc.Encrypt(&p)
	fmt.Println(*pp)
	ppp := bc.Decrypt(pp)
	fmt.Println(*ppp)
	if bytes.Compare(p, *ppp) != 0 {
		t.Fatal("Did not correctly encrypt and decrypt")
	}
	IV := bc.IV
	Crypt := bc.Crypt
	Iterations := bc.Iterations
	nc := New()
	nc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), Iterations)
	nc.Unlock(pass)
	nc.Arm()
	fmt.Println(p)
	pp = nc.Encrypt(&p)
	fmt.Println(*pp)
	ppp = nc.Decrypt(pp)
	fmt.Println(*ppp)
	if bytes.Compare(p, *ppp) != 0 {
		t.Fatal("Did not correctly encrypt and decrypt")
	}
	var ec *BlockCrypt
	err := errors.New("")
	ec.SetStatus("test")
	ec.SetStatusIf(err)
	ec.OK()
	ec.UnsetStatus()
	_ = ec.Error()
	nc.OK()
	_ = nc.Error()
	nc.Disarm()
	nc.Arm()
	ec.Generate(nil)
	ec.LoadCrypt(nil, nil, 0)
	ec.Unlock(nil)
	ec.Lock()
	ec.Arm()
	ec.Encrypt(nil)
	ec.Decrypt(nil)
	IV = bc.IV
	Crypt = bc.Crypt
	Iterations = bc.Iterations
	sc := New()
	sc.LoadCrypt(Crypt.Bytes(), nil, Iterations)
	sc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), -1)
	wrongIV := make([]byte, 11)
	sc.LoadCrypt(Crypt.Bytes(), &wrongIV, Iterations)
	sc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), 0)
	sc.Ciphertext = buf.NewSecure()
	sc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), Iterations)
	sc.Password = buf.NewSecure()
	sc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), Iterations)
	sc.Password = buf.NewSecure()
	sc.Unlock(pass)
	sc.Password = buf.NewSecure()
	sc.Lock()
	sc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), -1).Arm()
	sc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), Iterations)
	sc.Crypt = nil
	sc.Arm()
	sc.Password.Free()
	sc.Password = nil
	sc.LoadCrypt(Crypt.Bytes(), IV.Bytes(), Iterations).Arm()
	bc.Encrypt(nil)
	bc.Decrypt(nil)
	bc.Encrypt(&[]byte{})
	bc.Decrypt(&[]byte{})
	bc.GCM = nil
	bc.Encrypt(&[]byte{1, 2, 3, 4, 5, 6, 7, 8})
	bc.Decrypt(&[]byte{1, 2, 3, 4, 5, 6, 7, 8})
	sc.Iterations = 0
	sc.Arm()
}
