package transaction

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func SaveTransaction(msg *wire.MsgTx, key *db.Key, block *db.Block) error {
	txn := db.ConvertMsgToTransaction(msg)
	FindAndSaveMemos(txn, block)
	if block != nil {
		txn.BlockId = block.Id
		txn.Block = block
	}
	existingTransactions, err := db.GetTransactionsForKey(txn.KeyId)
	if err != nil {
		return jerr.Get("error getting transactions for key", err)
	}
	var updateOldOutput struct {
		txOut *db.TransactionOut
		txIn  *db.TransactionIn
	}
	var updateOldInput struct {
		txOut *db.TransactionOut
		txIn  *db.TransactionIn
	}
	for _, existingTransaction := range existingTransactions {
		// Check if inputs come from known outputs
		for _, in := range txn.TxIn {
			if ! bytes.Equal(in.PreviousOutPointHash, existingTransaction.Hash) {
				continue
			}
			var txOut *db.TransactionOut
			for _, existingTxOut := range existingTransaction.TxOut {
				if uint32(existingTxOut.Index) == in.PreviousOutPointIndex {
					txOut = existingTxOut
				}
			}
			if txOut == nil {
				return jerr.New("error finding matching txOut!")
			}
			in.TxnOutId = txOut.Id
			fmt.Printf("matched existing txn: %s\n", existingTransaction.GetChainHash().String())
			updateOldOutput.txOut = txOut
			updateOldOutput.txIn = in
		}
		// Check if outputs have been used in known inputs
		for _, in := range existingTransaction.TxIn {
			if ! bytes.Equal(in.PreviousOutPointHash, txn.Hash) {
				continue
			}
			for _, out := range txn.TxOut {
				if out.Index != in.PreviousOutPointIndex {
					continue
				}
				in.TxnOutId = out.Id
				out.TxnInId = in.Id
				out.TxnIn = in
				updateOldInput.txOut = out
				updateOldInput.txIn = in
			}
		}
	}

	existingTxn, err := db.GetTransactionByHash(txn.Hash)
	if err != nil && !db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting txn from db", err)
	}
	if existingTxn != nil {
		var updated bool
		if existingTxn.BlockId == 0 && txn.BlockId != 0 {
			existingTxn.BlockId = txn.BlockId
			existingTxn.Block = txn.Block
			updated = true
		}
		if existingTxn.KeyId == 0 && txn.KeyId != 0 {
			existingTxn.KeyId = txn.KeyId
			existingTxn.Key = txn.Key
			updated = true
		}
		if !updated {
			return nil
		}
		fmt.Println("Updating existing txn...")
		txn = existingTxn
	} else {
		fmt.Println("Found new txn...")
		txn.Key = key
		txn.KeyId = txn.Key.Id
	}
	err = txn.Save()
	if err != nil {
		return jerr.Get("error saving txn", err)
	}
	if updateOldOutput.txOut != nil && updateOldOutput.txIn != nil {
		updateOldOutput.txOut.TxnInId = updateOldOutput.txIn.Id
		err := updateOldOutput.txOut.Save()
		if err != nil {
			return jerr.Get("error updating old txn output", err)
		}
	}
	if updateOldInput.txOut != nil && updateOldInput.txIn != nil {
		updateOldInput.txIn.TxnOutId = updateOldInput.txOut.Id
		err := updateOldInput.txIn.Save()
		if err != nil {
			return jerr.Get("error updating old txn input", err)
		}
	}
	return nil
}
