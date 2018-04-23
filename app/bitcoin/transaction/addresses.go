package transaction

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/cpacia/btcd/txscript"
	"github.com/cpacia/btcd/wire"
)

func GetAddressesFromTxn(txn *wire.MsgTx) []*wallet.PublicKey {
	var addresses []*wallet.PublicKey
	for _, in := range txn.TxIn {
		address := GetPubKeyInInput(in)
		if address.GetSerializedString() != "" {
			addresses = append(addresses, &address)
		}
	}
	for _, out := range txn.TxOut {
		address := GetPubKeyInOutput(out)
		if address.GetSerializedString() != "" {
			addresses = append(addresses, &address)
		}
	}
	for i := range addresses {
		for g := range addresses {
			if i == g {
				continue
			}
			if bytes.Equal(addresses[i].GetSerialized(), addresses[g].GetSerialized()) {
				addresses = append(addresses[:g], addresses[g+1:]...)
				g--
			}
		}
	}
	return addresses
}

func GetPubKeyInInput(in *wire.TxIn) wallet.PublicKey {
	return wallet.GetPublicKey(in.SignatureScript)
}

func GetPubKeyInOutput(out *wire.TxOut) wallet.PublicKey {
	_, addresses, _, _ := txscript.ExtractPkScriptAddrs(out.PkScript, &wallet.MainNetParamsOld)
	if len(addresses) != 1 {
		return wallet.PublicKey{}
	}
	return wallet.GetPublicKey(addresses[0].ScriptAddress())
}
