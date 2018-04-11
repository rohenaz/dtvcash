package transaction

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
)

func FindAndSaveMemos(txn *db.Transaction) {
	var pkHash []byte
	for _, in := range txn.TxIn {
		pkHash = in.GetAddress().GetScriptAddress()
	}
	addressPkHash, err := btcutil.NewAddressPubKeyHash(pkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo test", err).Print()
		return
	}
	address := addressPkHash.EncodeAddress()
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
			TxHash:   txn.GetChainHash().CloneBytes(),
			PkHash:   pkHash,
			PkScript: out.PkScript,
			Address:  address,
		}
		err := test.Save()
		if err != nil {
			jerr.Get("error saving memo test", err).Print()
			return
		}
		switch out.PkScript[3] {
		case memo.CodePost:
			var post = db.MemoPost{
				TxHash:   txn.GetChainHash().CloneBytes(),
				PkHash:   pkHash,
				PkScript: out.PkScript,
				Address:  address,
				Message:  string(out.PkScript[5:]),
			}
			err := post.Save()
			if err != nil {
				jerr.Get("error saving memo post", err).Print()
				return
			}
		}
	}
}
