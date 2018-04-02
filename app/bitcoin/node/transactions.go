package node

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/btcd/txscript"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func onTx(n *Node, msg *wire.MsgTx) {
	n.CheckedTxns++
	scriptAddresses := getScriptAddresses(n)
	var found *wallet.Address
	var txnInfo string
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
	if found != nil {
		transaction := db.ConvertMsgToTransaction(msg)
		transaction.Block = findHashBlock(n, transaction.GetChainHash())
		if transaction.Block != nil {
			transaction.BlockId = transaction.Block.Id
		}
		existingTransactions, err := db.GetTransactionsForKey(transaction.KeyId)
		if err != nil {
			jerr.Get("error getting transactions for key", err).Print()
			return
		}
		var updateOldOutput struct {
			txOut *db.TransactionOut
			txIn  *db.TransactionIn
		}
		for _, existingTransaction := range existingTransactions {
			for _, in := range transaction.TxIn {
				if bytes.Equal(in.PreviousOutPointHash, existingTransaction.Hash) {
					var txOut *db.TransactionOut
					for _, existingTxOut := range existingTransaction.TxOut {
						if uint32(existingTxOut.Index) == in.PreviousOutPointIndex {
							txOut = existingTxOut
						}
					}
					if txOut == nil {
						jerr.New("error finding matching txOut!").Print()
						return
					}
					in.TxnOutId = txOut.Id
					txnInfo += fmt.Sprintf("matched existing transaction: %s\n", existingTransaction.GetChainHash().String())
					updateOldOutput.txOut = txOut
					updateOldOutput.txIn = in
				}
			}
		}

		existingTransaction, err := db.GetTransactionByHash(transaction.Hash)
		if err != nil && !db.IsRecordNotFoundError(err) {
			fmt.Println(jerr.Get("error getting transaction from db", err))
			return
		}
		if existingTransaction != nil {
			if existingTransaction.BlockId != 0 || transaction.BlockId == 0 {
				// Only thing that should update is potentially block height
				return
			}
			fmt.Printf("Updating existing transaction...\n" + txnInfo)
			existingTransaction.BlockId = transaction.BlockId
			existingTransaction.Block = transaction.Block
			transaction = existingTransaction
		} else {
			fmt.Printf("Found new transaction...\n" + txnInfo)
			transaction.Key = getKeyFromScriptAddress(n, found)
			transaction.KeyId = transaction.Key.Id
		}
		err = transaction.Save()
		if err != nil {
			jerr.Get("error saving transaction", err).Print()
			return
		}
		if updateOldOutput.txOut != nil {
			updateOldOutput.txOut.TxnInId = updateOldOutput.txIn.Id
			err := updateOldOutput.txOut.Save()
			if err != nil {
				jerr.Get("error updating old transaction output", err).Print()
				return
			}
		}
	}
}

func getTransaction(n *Node, txId chainhash.Hash) {
	msgGetData := wire.NewMsgGetData()
	err := msgGetData.AddInvVect(&wire.InvVect{
		Type: wire.InvTypeTx,
		Hash: txId,
	})
	if err != nil {
		jerr.Get("error adding invVect: %s\n", err).Print()
		return
	}
	n.Peer.QueueMessage(msgGetData, nil)
}
