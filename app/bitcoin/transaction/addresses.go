package transaction

import (
	"bytes"
	"github.com/memocash/memo/app/db"
)

func GetPkHashesFromTxn(dbTxn *db.Transaction) [][]byte {
	var pkHashes [][]byte
	for _, in := range dbTxn.TxIn {
		if len(in.KeyPkHash) > 0 {
			pkHashes = append(pkHashes, in.KeyPkHash)
		}
	}
	for _, out := range dbTxn.TxOut {
		if len(out.KeyPkHash) > 0 {
			pkHashes = append(pkHashes, out.KeyPkHash)
		}
	}
	for i := 0; i < len(pkHashes); i++ {
		for g := 0; g < len(pkHashes); g++ {
			if i == g {
				continue
			}
			if bytes.Equal(pkHashes[i], pkHashes[g]) {
				pkHashes = append(pkHashes[:g], pkHashes[g+1:]...)
				g--
			}
		}
	}
	return pkHashes
}
