package wallet

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
)

type KeyDB struct {
	Keys map[string]*btcec.PrivateKey
}


func (k KeyDB) GetKey(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
	for address, key  := range k.Keys {
		if addr.String() == address {
			return key, true, nil
		}
	}
	return nil, true, jerr.New("key not found")
}

type ScriptDb struct {}
func (s ScriptDb) GetScript(btcutil.Address) ([]byte, error) {
	return []byte{}, jerr.New("no scripts")
}
