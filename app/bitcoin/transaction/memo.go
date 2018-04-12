package transaction

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
	"html"
)

func FindAndSaveMemos(txn *db.Transaction, block *db.Block) error {
	var pkHash []byte
	for _, in := range txn.TxIn {
		tmpPkHash := in.GetAddress().GetScriptAddress()
		if len(tmpPkHash) > 0 {
			pkHash = tmpPkHash
		}
	}
	if len(pkHash) == 0 {
		// Unknown script type
		return nil
	}
	addressPkHash, err := btcutil.NewAddressPubKeyHash(pkHash, &wallet.MainNetParamsOld)
	if err != nil {
		return jerr.Get("error getting pubkeyhash from memo test", err)
	}
	address := addressPkHash.EncodeAddress()
	txnHash := txn.GetChainHash().CloneBytes()
	for _, out := range txn.TxOut {
		if len(out.PkScript) < 5 || ! bytes.Equal(out.PkScript[0:3], []byte{
			txscript.OP_RETURN,
			txscript.OP_DATA_2,
			memo.CodePrefix,
		}) {
			continue
		}
		// Save MemoTest
		var test = db.MemoTest{
			TxHash:   txnHash,
			PkHash:   pkHash,
			PkScript: out.PkScript,
			Address:  address,
		}
		if block != nil {
			test.BlockId = block.Id
		}
		err := test.Save()
		if err != nil {
			if ! db.IsAlreadyExistsError(err) || block == nil {
				return jerr.Get("error saving memo test", err)
			}
			memoTest, err := db.GetMemoTest(txnHash)
			if err != nil {
				return jerr.Get("error getting existing memo test", err)
			}
			memoTest.BlockId = block.Id
			err = memoTest.Save()
			if err != nil {
				return jerr.Get("error saving existing memo test", err)
			}
		}
		switch out.PkScript[3] {
		case memo.CodePost:
			var post = db.MemoPost{
				TxHash:   txnHash,
				PkHash:   pkHash,
				PkScript: out.PkScript,
				Address:  address,
				Message:  html.EscapeString(string(out.PkScript[5:])),
			}
			if len(txn.TxIn) == 1 {
				post.ParentHash = txn.TxIn[0].PreviousOutPointHash
			}
			if block != nil {
				post.BlockId = block.Id
			}
			err := post.Save()
			if err != nil {
				if ! db.IsAlreadyExistsError(err) || block == nil {
					return jerr.Get("error saving memo post", err)
				}
				post, err := db.GetMemoTest(txnHash)
				if err != nil {
					return jerr.Get("error getting existing memo post", err)
				}
				post.BlockId = block.Id
				err = post.Save()
				if err != nil {
					return jerr.Get("error saving existing memo post", err)
				}
			}
		case memo.CodeSetName:
			var setName = db.MemoSetName{
				TxHash:   txnHash,
				PkHash:   pkHash,
				PkScript: out.PkScript,
				Address:  address,
				Name:     html.EscapeString(string(out.PkScript[5:])),
			}
			if len(txn.TxIn) == 1 {
				setName.ParentHash = txn.TxIn[0].PreviousOutPointHash
			}
			if block != nil {
				setName.BlockId = block.Id
			}
			err := setName.Save()
			if err != nil {
				if ! db.IsAlreadyExistsError(err) || block == nil {
					return jerr.Get("error saving memo set name", err)
				}
				setName, err := db.GetMemoTest(txnHash)
				if err != nil {
					return jerr.Get("error getting existing memo set name", err)
				}
				setName.BlockId = block.Id
				err = setName.Save()
				if err != nil {
					return jerr.Get("error saving existing memo set name", err)
				}
			}
		}
	}
	return nil
}
