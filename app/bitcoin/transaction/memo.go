package transaction

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
	"html"
)

func GetMemoOutputIfExists(txn *db.Transaction) (*db.TransactionOut, error) {
	var out *db.TransactionOut
	for _, txOut := range txn.TxOut {
		if len(txOut.PkScript) < 5 || ! bytes.Equal(txOut.PkScript[0:3], []byte{
			txscript.OP_RETURN,
			txscript.OP_DATA_2,
			memo.CodePrefix,
		}) {
			continue
		}
		if out != nil {
			return nil, jerr.New("UNEXPECTED ERROR: found more than one memo in transaction")
		}
		out = txOut
	}
	return out, nil
}

func SaveMemo(txn *db.Transaction, out *db.TransactionOut, block *db.Block) error {
	_, err := db.GetMemoTest(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_test", err)
	}
	if err == nil {
		err = updateMemo(txn, out, block)
		if err != nil {
			return jerr.Get("error updating memo", err)
		}
	} else {
		err = newMemo(txn, out, block)
		if err != nil {
			return jerr.Get("error saving new memo", err)
		}
	}
	return nil
}

func getInputPkHash(txn *db.Transaction) (*btcutil.AddressPubKeyHash, error) {
	var pkHash []byte
	for _, in := range txn.TxIn {
		tmpPkHash := in.GetAddress().GetScriptAddress()
		if len(tmpPkHash) > 0 {
			if len(pkHash) != 0 && ! bytes.Equal(tmpPkHash, pkHash) {
				return nil, jerr.New("error found multiple addresses in inputs")
			}
			pkHash = tmpPkHash
		}
	}
	if len(pkHash) == 0 {
		// Unknown script type
		return nil, jerr.New("error no pk hash found")
	}
	addressPkHash, err := btcutil.NewAddressPubKeyHash(pkHash, &wallet.MainNetParamsOld)
	if err != nil {
		return nil, jerr.Get("error getting pubkeyhash from memo test", err)
	}
	return addressPkHash, nil
}

func newMemo(txn *db.Transaction, out *db.TransactionOut, block *db.Block) error {
	fmt.Printf("Saving new memo (txn: %s)\n", txn.GetChainHash().String())
	inputAddress, err := getInputPkHash(txn)
	if err != nil {
		return jerr.Get("error getting pk hash from input", err)
	}
	var blockId uint
	if block != nil {
		blockId = block.Id
	}
	// Used for ordering
	var parentHash []byte
	if len(txn.TxIn) == 1 {
		parentHash = txn.TxIn[0].PreviousOutPointHash
	}
	var memoTest = db.MemoTest{
		TxHash:   txn.Hash,
		PkHash:   inputAddress.ScriptAddress(),
		PkScript: out.PkScript,
		Address:  inputAddress.EncodeAddress(),
		BlockId:  blockId,
	}
	err = memoTest.Save()
	if err != nil {
		return jerr.Get("error saving memo_test", err)
	}
	switch out.PkScript[3] {
	case memo.CodePost:
		var memoPost = db.MemoPost{
			TxHash:     txn.Hash,
			PkHash:     inputAddress.ScriptAddress(),
			PkScript:   out.PkScript,
			ParentHash: parentHash,
			Address:    inputAddress.EncodeAddress(),
			Message:    html.EscapeString(string(out.PkScript[5:])),
			BlockId:    blockId,
		}
		err := memoPost.Save()
		if err != nil {
			return jerr.Get("error saving memo_post", err)
		}
	case memo.CodeSetName:
		var memoSetName = db.MemoSetName{
			TxHash:     txn.Hash,
			PkHash:     inputAddress.ScriptAddress(),
			PkScript:   out.PkScript,
			ParentHash: parentHash,
			Address:    inputAddress.EncodeAddress(),
			Name:       html.EscapeString(string(out.PkScript[5:])),
			BlockId:    blockId,
		}
		err := memoSetName.Save()
		if err != nil {
			return jerr.Get("error saving memo_set_name", err)
		}
	}
	return nil
}

func updateMemo(txn *db.Transaction, out *db.TransactionOut, block *db.Block) error {
	fmt.Printf("Updating existing memo (txn: %s)\n", txn.GetChainHash().String())
	memoTest, err := db.GetMemoTest(txn.Hash)
	if err != nil {
		return jerr.Get("error getting memo_test", err)
	}
	if block == nil || memoTest.BlockId != 0 {
		// Nothing to update
		return nil
	}
	memoTest.BlockId = block.Id
	err = memoTest.Save()
	if err != nil {
		return jerr.Get("error saving memo_test", err)
	}
	switch out.PkScript[3] {
	case memo.CodePost:
		memoPost, err := db.GetMemoPost(txn.Hash)
		if err != nil {
			return jerr.Get("error getting memo_post", err)
		}
		memoPost.BlockId = block.Id
		err = memoPost.Save()
		if err != nil {
			return jerr.Get("error saving memo_post", err)
		}
	case memo.CodeSetName:
		memoSetName, err := db.GetMemoSetName(txn.Hash)
		if err != nil {
			return jerr.Get("error getting memo_set_name", err)
		}
		memoSetName.BlockId = block.Id
		err = memoSetName.Save()
		if err != nil {
			return jerr.Get("error saving memo_set_name", err)
		}
	}
	return nil
}
