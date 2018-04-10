package transaction

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/memo/app/db"
	"github.com/cpacia/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
)

func FindAndSaveMemos(txn *db.Transaction) {
	var address []byte
	for _, in := range txn.TxIn {
		address = in.GetAddress().GetScriptAddress()
	}
	for _, out := range txn.TxOut {
		if len(out.PkScript) < 5 || ! bytes.Equal(out.PkScript[0:3], []byte{
			txscript.OP_RETURN,
			txscript.OP_DATA_2,
			memo.CodePrefix,
		}) {
			continue
		}
		var test = db.MemoTest{
			TxHash: txn.GetChainHash().CloneBytes(),
			Address: address,
			PkScript: out.PkScript,
		}
		err := test.Save()
		if err != nil {
			jerr.Get("error saving memo test", err).Print()
			return
		}
	}
}
