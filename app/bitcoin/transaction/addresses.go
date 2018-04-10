package transaction

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/cpacia/btcd/txscript"
	"github.com/cpacia/btcd/wire"
)

func GetAddressInTx(msg *wire.MsgTx, scriptAddresses []*wallet.Address) *wallet.Address {
	var found *wallet.Address
	var txnInfo = fmt.Sprintf("Txn: %s\n", msg.TxHash().String())
	for _, in := range msg.TxIn {
		for _, scriptAddress := range scriptAddresses {
			if bytes.Equal(in.SignatureScript, scriptAddress.GetScriptAddress()) {
				found = scriptAddress
			}
		}
		unlockScript, err := txscript.DisasmString(in.SignatureScript)
		if err != nil {
			txnInfo = txnInfo + fmt.Sprintf("Error disassembling unlockScript: %s\n", err.Error())
		}
		txnInfo = txnInfo + fmt.Sprintf("  TxIn - Sequence: %d\n"+
			"    prevOut: %s\n"+
			"    unlockScript: %s\n",
			in.Sequence, in.PreviousOutPoint.String(), unlockScript)
	}
	for _, out := range msg.TxOut {
		lockScript, err := txscript.DisasmString(out.PkScript)
		if err != nil {
			txnInfo = txnInfo + fmt.Sprintf("Error disassembling lockScript: %s\n", err.Error())
			continue
		}
		scriptClass, addresses, sigCount, err := txscript.ExtractPkScriptAddrs(out.PkScript, &wallet.MainNetParamsOld)
		for _, address := range addresses {
			for _, scriptAddress := range scriptAddresses {
				if bytes.Equal(address.ScriptAddress(), scriptAddress.GetScriptAddress()) {
					found = scriptAddress
				}
			}
			txnInfo = txnInfo + fmt.Sprintf("  TxOut - value: %d\n"+
				"    lockScript: %s\n"+
				"    scriptClass: %s\n"+
				"    requiredSigs: %d\n",
				out.Value, lockScript, scriptClass, sigCount)
			txnInfo = txnInfo + fmt.Sprintf("    address: %s\n", address.String())
		}
	}
	return found
}
