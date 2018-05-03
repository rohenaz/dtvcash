package transaction

import (
	"bytes"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/cache"
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
	err = saveMemoTest(txn, out, blockId, inputAddress)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error saving memo_test", err)
	}
	switch out.PkScript[3] {
	case memo.CodePost:
		err = saveMemoPost(txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_post", err)
		}
	case memo.CodeSetName:
		err = saveMemoSetName(txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_set_name", err)
		}
	case memo.CodeFollow:
		err = saveMemoFollow(txn, out, blockId, inputAddress, parentHash, false)
		if err != nil {
			return jerr.Get("error saving memo_follow", err)
		}
	case memo.CodeUnfollow:
		err = saveMemoFollow(txn, out, blockId, inputAddress, parentHash, true)
		if err != nil {
			return jerr.Get("error saving memo_follow", err)
		}
	case memo.CodeLike:
		err = saveMemoLike(txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_like", err)
		}
	case memo.CodeReply:
		err = saveMemoReply(txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_post reply", err)
		}
	case memo.CodeSetProfile:
		err = saveMemoSetProfile(txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_set_profile", err)
		}
	case memo.CodeTagMessage:
		err = saveMemoTagMessage(txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_post tag message", err)
		}
	}
	return nil
}

func saveMemoTest(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash) error {
	memoTest, err := db.GetMemoTest(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_test", err)
	}
	if memoTest != nil {
		if memoTest.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoTest.BlockId = blockId
		err = memoTest.Save()
		if err != nil {
			return jerr.Get("error saving memo_test", err)
		}
		return nil
	}
	memoTest = &db.MemoTest{
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
	return nil
}

func saveMemoPost(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
	memoPost, err := db.GetMemoPost(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_post", err)
	}
	if memoPost != nil {
		if memoPost.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoPost.BlockId = blockId
		err = memoPost.Save()
		if err != nil {
			return jerr.Get("error saving memo_post", err)
		}
		return nil
	}
	var message string
	if len(out.PkScript) > 81 {
		message = string(out.PkScript[6:])
	} else {
		message = string(out.PkScript[5:])
	}
	memoPost = &db.MemoPost{
		TxHash:     txn.Hash,
		PkHash:     inputAddress.ScriptAddress(),
		PkScript:   out.PkScript,
		ParentHash: parentHash,
		Address:    inputAddress.EncodeAddress(),
		Message:    html_parser.EscapeWithEmojis(message),
		BlockId:    blockId,
	}
	err = memoPost.Save()
	if err != nil {
		return jerr.Get("error saving memo_post", err)
	}
	return nil
}

func saveMemoSetName(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
	memoSetName, err := db.GetMemoSetName(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_set_name", err)
	}
	if memoSetName != nil {
		if memoSetName.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoSetName.BlockId = blockId
		err = memoSetName.Save()
		if err != nil {
			return jerr.Get("error saving memo_set_name", err)
		}
		return nil
	}
	var name string
	if len(out.PkScript) > 81 {
		name = string(out.PkScript[6:])
	} else {
		name = string(out.PkScript[5:])
	}
	memoSetName = &db.MemoSetName{
		TxHash:     txn.Hash,
		PkHash:     inputAddress.ScriptAddress(),
		PkScript:   out.PkScript,
		ParentHash: parentHash,
		Address:    inputAddress.EncodeAddress(),
		Name:       html_parser.EscapeWithEmojis(name),
		BlockId:    blockId,
	}
	err = memoSetName.Save()
	if err != nil {
		return jerr.Get("error saving memo_set_name", err)
	}
	return nil
}

func saveMemoFollow(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte, unfollow bool) error {
	memoFollow, err := db.GetMemoFollow(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_follow", err)
	}
	if memoFollow != nil {
		if memoFollow.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoFollow.BlockId = blockId
		err = memoFollow.Save()
		if err != nil {
			return jerr.Get("error saving memo_follow", err)
		}
		return nil
	}
	address := wallet.GetAddressFromPkHash(out.PkScript[5:])
	if ! bytes.Equal(address.GetScriptAddress(), out.PkScript[5:]) {
		return jerr.New("unable to parse follow address")
	}
	memoFollow = &db.MemoFollow{
		TxHash:       txn.Hash,
		PkHash:       inputAddress.ScriptAddress(),
		PkScript:     out.PkScript,
		ParentHash:   parentHash,
		Address:      inputAddress.EncodeAddress(),
		FollowPkHash: address.GetScriptAddress(),
		BlockId:      blockId,
		Unfollow:     unfollow,
	}
	err = memoFollow.Save()
	if err != nil {
		return jerr.Get("error saving memo_follow", err)
	}
	err = cache.ClearReputation(memoFollow.PkHash, memoFollow.FollowPkHash)
	if err != nil && ! cache.IsMissError(err) {
		return jerr.Get("error clearing cache", err)
	}
	return nil
}

func saveMemoLike(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
	memoLike, err := db.GetMemoLike(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_like", err)
	}
	if memoLike != nil {
		if memoLike.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoLike.BlockId = blockId
		err = memoLike.Save()
		if err != nil {
			return jerr.Get("error saving memo_like", err)
		}
		return nil
	}

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
	memoLike = &db.MemoLike{
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
	return nil
}

func saveMemoReply(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
	memoPost, err := db.GetMemoPost(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_reply", err)
	}
	if memoPost != nil {
		if memoPost.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoPost.BlockId = blockId
		err = memoPost.Save()
		if err != nil {
			return jerr.Get("error saving memo_reply", err)
		}
		return nil
	}

	if len(out.PkScript) < 38 {
		return jerr.Newf("invalid reply, length too short (%d)", len(out.PkScript))
	}
	txHash, err := chainhash.NewHash(out.PkScript[5:37])
	if err != nil {
		return jerr.Get("error parsing transaction hash", err)
	}
	memoPost = &db.MemoPost{
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
	return nil
}

func saveMemoTagMessage(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
	memoPost, err := db.GetMemoPost(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo tag message", err)
	}
	if memoPost != nil {
		if memoPost.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoPost.BlockId = blockId
		err = memoPost.Save()
		if err != nil {
			return jerr.Get("error saving memo tag message", err)
		}
		return nil
	}

	pushData, err := txscript.PushedData(out.PkScript)
	if err != nil {
		return jerr.Get("error parsing push data from memo tag message", err)
	}
	if len(pushData) != 3 {
		return jerr.Newf("invalid tag message, incorrect push data (%d)", len(pushData))
	}
	var tagNameRaw = pushData[1]
	var messageRaw = pushData[2]
	if len(tagNameRaw) == 0 || len(messageRaw) == 0 {
		return jerr.Newf("empty tag or message (%d, %d)", len(tagNameRaw), len(messageRaw))
	}
	tagName := html_parser.EscapeWithEmojis(string(tagNameRaw))
	message := html_parser.EscapeWithEmojis(string(messageRaw))
	memoPost = &db.MemoPost{
		TxHash:     txn.Hash,
		PkHash:     inputAddress.ScriptAddress(),
		PkScript:   out.PkScript,
		ParentHash: parentHash,
		Address:    inputAddress.EncodeAddress(),
		Tag:        tagName,
		Message:    message,
		BlockId:    blockId,
	}
	err = memoPost.Save()
	if err != nil {
		return jerr.Get("error saving memo tag message", err)
	}
	return nil
}

func saveMemoSetProfile(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
	memoSetProfile, err := db.GetMemoSetProfile(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_set_profile", err)
	}
	if memoSetProfile != nil {
		if memoSetProfile.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoSetProfile.BlockId = blockId
		err = memoSetProfile.Save()
		if err != nil {
			return jerr.Get("error saving memo_set_profile", err)
		}
		return nil
	}

	var profile string
	if len(out.PkScript) > 81 {
		profile = string(out.PkScript[6:])
	} else {
		profile = string(out.PkScript[5:])
	}
	memoSetProfile = &db.MemoSetProfile{
		TxHash:     txn.Hash,
		PkHash:     inputAddress.ScriptAddress(),
		PkScript:   out.PkScript,
		ParentHash: parentHash,
		Address:    inputAddress.EncodeAddress(),
		Profile:    html_parser.EscapeWithEmojis(profile),
		BlockId:    blockId,
	}
	err = memoSetProfile.Save()
	if err != nil {
		return jerr.Get("error saving memo_set_profile", err)
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
