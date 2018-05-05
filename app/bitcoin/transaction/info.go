package transaction

import (
	"fmt"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/cpacia/btcd/txscript"
	"github.com/cpacia/btcd/wire"
)

func GetTxInfo(msg *wire.MsgTx) string {
	var txnInfo = fmt.Sprintf("Txn: %s\n", msg.TxHash().String())
	for _, in := range msg.TxIn {
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
		txnInfo = txnInfo + fmt.Sprintf("  TxOut - value: %d\n"+
			"    lockScript: %s\n", out.Value, lockScript)
		scriptClass, addresses, sigCount, err := txscript.ExtractPkScriptAddrs(out.PkScript, &wallet.MainNetParamsOld)
		for _, address := range addresses {
			txnInfo = txnInfo + fmt.Sprintf("    address: %s\n"+
				"    scriptClass: %s\n"+
				"    requiredSigs: %d\n",
				address.String(), scriptClass, sigCount)
		}
	}
	return txnInfo
}
