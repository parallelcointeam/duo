package wallet

import (
	"testing"

	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/wallet/db"
)

func TestNewKeyPool(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	wdb := db.NewWalletDB().WithBC(BC)
	W := New(wdb)
	W.NewKeyPool()
	W.DB.Dump()
	W.DB.DeleteAll()
	// W.DB.Dump()
	W.DB.Close()
}
