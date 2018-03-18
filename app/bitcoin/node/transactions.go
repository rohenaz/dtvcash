package node

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func onTx(n *Node, msg *wire.MsgTx) {
	n.CheckedTxns++
	//fmt.Printf("Transaction - version: %d, locktime: %d, inputs: %d, outputs: %d\n", msg.Version, msg.LockTime, len(msg.TxIn), len(msg.TxOut))
	scriptAddresses := getScriptAddresses(n)
	var found bool
	var txnInfo string
	for _, in := range msg.TxIn {
		for _, scriptAddress := range scriptAddresses {
			if bytes.Equal(in.SignatureScript, scriptAddress.GetScriptAddress()) {
				found = true
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
		scriptClass, addresses, sigCount, err := txscript.ExtractPkScriptAddrs(out.PkScript, &wallet.MainNetParams)
		for _, address := range addresses {
			for _, scriptAddress := range scriptAddresses {
				if bytes.Equal(address.ScriptAddress(), scriptAddress.GetScriptAddress()) {
					found = true
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
	if found {
		txnInfo = "Found transaction...\n" + txnInfo
		fmt.Printf(txnInfo)

		/*var transaction = db.Transaction{
			Address: scriptAddress,
		}
		err := transaction.Save()
		if err != nil {
			fmt.Println(jerr.Get("error saving transaction", err))
			return
		}*/
	}
}

func getTransaction(n *Node, txId chainhash.Hash) {
	msgGetData := wire.NewMsgGetData()
	err := msgGetData.AddInvVect(&wire.InvVect{
		Type: wire.InvTypeTx,
		Hash: txId,
	})
	if err != nil {
		fmt.Printf("error adding invVect: %s\n", err)
		return
	}
	n.Peer.QueueMessage(msgGetData, nil)
}
