package transaction

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func SaveTransaction(msg *wire.MsgTx, block *db.Block) error {
	hash := msg.TxHash()
	txn, err := db.GetTransactionByHash(hash.CloneBytes())
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting transaction from db", err)
	}
	if txn != nil {
		err = updateTxn(txn, block)
		if err != nil {
			return jerr.Get("error updating transaction", err)
		}
	} else {
		txn, err = db.ConvertMsgToTransaction(msg)
		if err != nil {
			return jerr.Get("error converting message to transaction", err)
		}
		err = newTxn(txn, block)
		if err != nil {
			return jerr.Get("error saving new transaction", err)
		}
	}
	memoOutput, err := GetMemoOutputIfExists(txn)
	if err != nil {
		return jerr.Get("error getting memo output", err)
	}
	if memoOutput != nil {
		err = SaveMemo(txn, memoOutput, block)
		if err != nil {
			return jerr.Get("error saving memos", err)
		}
	}
	return nil
}

func updateTxn(txn *db.Transaction, block *db.Block) error {
	if block == nil {
		// Nothing to update
		return nil
	}
	txn.BlockId = block.Id
	err := txn.Save()
	if err != nil {
		return jerr.Get("error saving updated transaction", err)
	}
	return nil
}

func newTxn(txn *db.Transaction, block *db.Block) error {
	var blockId uint
	if block != nil {
		blockId = block.Id
	}
	txn.BlockId = blockId
	fmt.Printf("Found new txn: %s, block id: %d\n", txn.GetChainHash().String(), blockId)
	for _, in := range txn.TxIn {
		err := updateExistingOutputs(in)
		if err != nil {
			return jerr.Get("error updating existing outputs", err)
		}
	}

	for _, out := range txn.TxOut {
		err := updateExistingInputs(out, txn.Hash)
		if err != nil {
			return jerr.Get("error updating existing inputs", err)
		}
	}

	err := txn.Save()
	if err != nil {
		return jerr.Get("error saving txn", err)
	}

	return nil
}

func updateExistingOutputs(in *db.TransactionIn) error {
	existingOutputs, err := db.GetTransactionOutputsForPkHash(in.KeyPkHash)
	if err != nil {
		return jerr.Get("error getting existing outputs", err)
	}
	for _, existingTxOut := range existingOutputs {
		if ! bytes.Equal(in.PreviousOutPointHash, existingTxOut.Transaction.Hash) {
			continue
		}
		if uint32(existingTxOut.Index) == in.PreviousOutPointIndex {
			existingTxOut.TxnInHashString = in.HashString
			err := existingTxOut.Save()
			if err != nil {
				return jerr.Get("error saving existing transaction output", err)
			}
		}
	}
	return nil
}

func updateExistingInputs(out *db.TransactionOut, txHash []byte) error {
	existingInputs, err := db.GetTransactionInputsForPkHash(out.KeyPkHash)
	if err != nil {
		return jerr.Get("error getting existing inputs", err)
	}
	for _, existingTxIn := range existingInputs {
		if bytes.Equal(txHash, existingTxIn.PreviousOutPointHash) && out.Index == existingTxIn.PreviousOutPointIndex {
			out.TxnInHashString = existingTxIn.HashString
		}
	}
	return nil
}
