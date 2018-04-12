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

	txn, err = db.ConvertMsgToTransaction(msg)
	if err != nil {
		return jerr.Get("error converting message to transaction", err)
	}

	fmt.Printf("Found new txn: %s\n", txn.GetChainHash().String())

	if block != nil {
		txn.BlockId = block.Id
	}

	for _, in := range txn.TxIn {
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
	}

	for _, out := range txn.TxOut {
		existingInputs, err := db.GetTransactionInputsForPkHash(out.KeyPkHash)
		if err != nil {
			return jerr.Get("error getting existing inputs", err)
		}
		for _, existingTxIn := range existingInputs {
			if bytes.Equal(txn.Hash, existingTxIn.PreviousOutPointHash) && out.Index == existingTxIn.PreviousOutPointIndex {
				out.TxnInHashString = existingTxIn.HashString
			}
		}
	}

	err = txn.Save()
	if err != nil {
		return jerr.Get("error saving txn", err)
	}

	err = FindAndSaveMemos(txn, block)
	if err != nil {
		return jerr.Get("error finding and saving memos", err)
	}

	return nil
}
