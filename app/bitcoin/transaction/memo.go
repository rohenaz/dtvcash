package transaction

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/html-parser"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
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
	//fmt.Printf("Found new memo (txn: %s)\n", txn.GetChainHash().String())
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
			Message:    html_parser.EscapeWithEmojis(string(out.PkScript[5:])),
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
			Name:       html_parser.EscapeWithEmojis(string(out.PkScript[5:])),
			BlockId:    blockId,
		}
		err := memoSetName.Save()
		if err != nil {
			return jerr.Get("error saving memo_set_name", err)
		}
	case memo.CodeFollow:
		address := wallet.GetAddressFromPkHash(out.PkScript[5:])
		if ! bytes.Equal(address.GetScriptAddress(), out.PkScript[5:]) {
			return jerr.New("unable to parse follow address")
		}
		var memoFollow = db.MemoFollow{
			TxHash:       txn.Hash,
			PkHash:       inputAddress.ScriptAddress(),
			PkScript:     out.PkScript,
			ParentHash:   parentHash,
			Address:      inputAddress.EncodeAddress(),
			FollowPkHash: address.GetScriptAddress(),
			BlockId:      blockId,
		}
		err := memoFollow.Save()
		if err != nil {
			return jerr.Get("error saving memo_follow", err)
		}
	case memo.CodeUnfollow:
		address := wallet.GetAddressFromPkHash(out.PkScript[5:])
		if ! bytes.Equal(address.GetScriptAddress(), out.PkScript[5:]) {
			return jerr.New("unable to parse follow address")
		}
		var memoFollow = db.MemoFollow{
			TxHash:       txn.Hash,
			PkHash:       inputAddress.ScriptAddress(),
			PkScript:     out.PkScript,
			ParentHash:   parentHash,
			Address:      inputAddress.EncodeAddress(),
			FollowPkHash: address.GetScriptAddress(),
			BlockId:      blockId,
			Unfollow:     true,
		}
		err := memoFollow.Save()
		if err != nil {
			return jerr.Get("error saving memo_follow", err)
		}
	case memo.CodeLike:
		txHash, err := chainhash.NewHash(out.PkScript[5:37])
		if err != nil {
			return jerr.Get("error parsing transaction hash", err)
		}
		var tipPkHash []byte
		var tipAmount int64
		for _, txOut := range txn.TxOut {
			if len(txOut.KeyPkHash) == 0 || bytes.Equal(txOut.KeyPkHash, inputAddress.ScriptAddress()) {
				continue
			}
			if len(tipPkHash) != 0 {
				return jerr.New("error found multiple tip outputs, unable to process")
			}
			tipAmount += txOut.Value
			tipPkHash = txOut.KeyPkHash
		}
		var memoLike = db.MemoLike{
			TxHash:     txn.Hash,
			PkHash:     inputAddress.ScriptAddress(),
			PkScript:   out.PkScript,
			ParentHash: parentHash,
			Address:    inputAddress.EncodeAddress(),
			LikeTxHash: txHash.CloneBytes(),
			BlockId:    blockId,
			TipPkHash:  tipPkHash,
			TipAmount:  tipAmount,
		}
		err = memoLike.Save()
		if err != nil {
			return jerr.Get("error saving memo_like", err)
		}
	case memo.CodeReply:
		if len(out.PkScript) < 38 {
			return jerr.Newf("invalid reply, length too short (%d)", len(out.PkScript))
		}
		txHash, err := chainhash.NewHash(out.PkScript[5:37])
		if err != nil {
			return jerr.Get("error parsing transaction hash", err)
		}
		var memoPost = db.MemoPost{
			TxHash:       txn.Hash,
			PkHash:       inputAddress.ScriptAddress(),
			PkScript:     out.PkScript,
			ParentHash:   parentHash,
			Address:      inputAddress.EncodeAddress(),
			ParentTxHash: txHash.CloneBytes(),
			Message:      html_parser.EscapeWithEmojis(string(out.PkScript[38:])),
			BlockId:      blockId,
		}
		err = memoPost.Save()
		if err != nil {
			return jerr.Get("error saving memo_reply", err)
		}
	}
	return nil
}

func updateMemo(txn *db.Transaction, out *db.TransactionOut, block *db.Block) error {
	//fmt.Printf("Updating existing memo (txn: %s)\n", txn.GetChainHash().String())
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
	case memo.CodeFollow:
		memoFollow, err := db.GetMemoFollow(txn.Hash)
		if err != nil {
			return jerr.Get("error getting memo_follow", err)
		}
		memoFollow.BlockId = block.Id
		err = memoFollow.Save()
		if err != nil {
			return jerr.Get("error saving memo_follow", err)
		}
	case memo.CodeUnfollow:
		memoFollow, err := db.GetMemoFollow(txn.Hash)
		if err != nil {
			return jerr.Get("error getting memo_follow", err)
		}
		memoFollow.BlockId = block.Id
		err = memoFollow.Save()
		if err != nil {
			return jerr.Get("error saving memo_follow", err)
		}
	case memo.CodeLike:
		memoLike, err := db.GetMemoLike(txn.Hash)
		if err != nil {
			return jerr.Get("error getting memo_like", err)
		}
		memoLike.BlockId = block.Id
		err = memoLike.Save()
		if err != nil {
			return jerr.Get("error saving memo_like", err)
		}
	case memo.CodeReply:
		memoPost, err := db.GetMemoPost(txn.Hash)
		if err != nil {
			return jerr.Get("error getting memo_reply", err)
		}
		memoPost.BlockId = block.Id
		err = memoPost.Save()
		if err != nil {
			return jerr.Get("error saving memo_reply", err)
		}
	}
	return nil
}
