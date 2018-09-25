package walletdb

import (
	"github.com/dgraph-io/badger"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/proto"
)

var er = proto.Errors

// DB is the central data repository for the wallet database
type DB struct {
	Path     string
	BaseDir  string
	ValueDir string
	Options  *badger.Options
	DB       *badger.DB
	BC       *bc.BlockCrypt
	proto.State
}
