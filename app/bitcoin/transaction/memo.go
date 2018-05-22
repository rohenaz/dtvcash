package transaction

import (
	"bytes"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/btcd/txscript"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/html-parser"
	"github.com/memocash/memo/app/notify"
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
	case memo.CodeTopicMessage:
		err = saveMemoTopicMessage(txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_post tag message", err)
		}
	case memo.CodePollSingle:
		err = saveMemoPollQuestion(memo.CodePollSingle, txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_poll_question (single)", err)
		}
	case memo.CodePollMulti:
		err = saveMemoPollQuestion(memo.CodePollMulti, txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo_poll_question (multi)", err)
		}
	case memo.CodePollOption:
		err = saveMemoPollOption(txn, out, blockId, inputAddress, parentHash)
		if err != nil {
			return jerr.Get("error saving memo poll option", err)
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
	pushData, err := txscript.PushedData(out.PkScript)
	if err != nil {
		return jerr.Get("error parsing push data from message", err)
	}
	if len(pushData) != 2 {
		return jerr.Newf("invalid message, incorrect push data (%d)", len(pushData))
	}
	var message = string(pushData[1])
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
	pushData, err := txscript.PushedData(out.PkScript)
	if err != nil {
		return jerr.Get("error parsing push data from set name", err)
	}
	if len(pushData) != 2 {
		return jerr.Newf("invalid set name, incorrect push data (%d)", len(pushData))
	}
	var name = string(pushData[1])
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
	go func() {
		err = notify.AddLikeNotification(memoLike, true)
		if err != nil {
			jerr.Get("error adding like notification", err).Print()
		}
	}()
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
	pushData, err := txscript.PushedData(out.PkScript)
	if err != nil {
		return jerr.Get("error parsing push data from memo reply message", err)
	}
	if len(pushData) != 3 {
		return jerr.Newf("invalid reply message, incorrect push data (%d)", len(pushData))
	}
	var replyTxHash = pushData[1]
	var messageRaw = pushData[2]
	txHash, err := chainhash.NewHash(replyTxHash)
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
		Message:      html_parser.EscapeWithEmojis(string(messageRaw)),
		BlockId:      blockId,
	}
	err = memoPost.Save()
	if err != nil {
		return jerr.Get("error saving memo_reply", err)
	}
	go func() {
		err = notify.AddReplyNotification(memoPost, true)
		if err != nil {
			jerr.Get("error adding reply notification", err).Print()
		}
	}()
	return nil
}

func saveMemoTopicMessage(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
	memoPost, err := db.GetMemoPost(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo topic message", err)
	}
	if memoPost != nil {
		if memoPost.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoPost.BlockId = blockId
		err = memoPost.Save()
		if err != nil {
			return jerr.Get("error saving memo topic message", err)
		}
		return nil
	}

	pushData, err := txscript.PushedData(out.PkScript)
	if err != nil {
		return jerr.Get("error parsing push data from memo topic message", err)
	}
	if len(pushData) != 3 {
		return jerr.Newf("invalid topic message, incorrect push data (%d)", len(pushData))
	}
	var topicNameRaw = pushData[1]
	var messageRaw = pushData[2]
	if len(topicNameRaw) == 0 || len(messageRaw) == 0 {
		return jerr.Newf("empty topic or message (%d, %d)", len(topicNameRaw), len(messageRaw))
	}
	topicName := html_parser.EscapeWithEmojis(string(topicNameRaw))
	message := html_parser.EscapeWithEmojis(string(messageRaw))
	memoPost = &db.MemoPost{
		TxHash:     txn.Hash,
		PkHash:     inputAddress.ScriptAddress(),
		PkScript:   out.PkScript,
		ParentHash: parentHash,
		Address:    inputAddress.EncodeAddress(),
		Topic:      topicName,
		Message:    message,
		BlockId:    blockId,
	}
	err = memoPost.Save()
	if err != nil {
		return jerr.Get("error saving memo topic message", err)
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
	pushData, err := txscript.PushedData(out.PkScript)
	if err != nil {
		return jerr.Get("error parsing push data from profile text", err)
	}
	if len(pushData) != 2 {
		return jerr.Newf("invalid profile text, incorrect push data (%d)", len(pushData))
	}
	var profile = string(pushData[1])
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

func saveMemoPollQuestion(pollType int, txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
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
			return jerr.Get("error saving memo_poll_question", err)
		}
		return nil
	}
	pushData, err := txscript.PushedData(out.PkScript)
	if err != nil {
		return jerr.Get("error parsing push data from poll question", err)
	}
	if len(pushData) != 3 {
		return jerr.Newf("invalid poll question, incorrect push data (%d)", len(pushData))
	}
	if len(pushData[1]) == 0 {
		return jerr.New("invalid push data for poll question, num options empty")
	}
	if len(pushData[2]) == 0 {
		return jerr.New("invalid push data for poll question, question empty")
	}
	var numOptions = uint(pushData[1][0])
	var question = string(pushData[2])
	memoPost = &db.MemoPost{
		TxHash:     txn.Hash,
		PkHash:     inputAddress.ScriptAddress(),
		PkScript:   out.PkScript,
		Message:    html_parser.EscapeWithEmojis(question),
		ParentHash: parentHash,
		Address:    inputAddress.EncodeAddress(),
		BlockId:    blockId,
		IsPoll:     true,
	}
	err = memoPost.Save()
	if err != nil {
		return jerr.Get("error saving memo_post for poll question", err)
	}
	memoPollQuestion := &db.MemoPollQuestion{
		TxHash:     txn.Hash,
		NumOptions: numOptions,
		PollType:   pollType,
	}
	err = memoPollQuestion.Save()
	if err != nil {
		return jerr.Get("error saving memo_set_profile", err)
	}
	return nil
}

func saveMemoPollOption(txn *db.Transaction, out *db.TransactionOut, blockId uint, inputAddress *btcutil.AddressPubKeyHash, parentHash []byte) error {
	memoPollOption, err := db.GetMemoPollOption(txn.Hash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return jerr.Get("error getting memo_poll_option", err)
	}
	if memoPollOption != nil {
		if memoPollOption.BlockId != 0 || blockId == 0 {
			return nil
		}
		memoPollOption.BlockId = blockId
		err = memoPollOption.Save()
		if err != nil {
			return jerr.Get("error saving memo_poll_option", err)
		}
		return nil
	}
	pushData, err := txscript.PushedData(out.PkScript)
	if err != nil {
		return jerr.Get("error parsing push data from poll option", err)
	}
	if len(pushData) != 2 {
		return jerr.Newf("invalid poll option, incorrect push data (%d)", len(pushData))
	}
	var option = string(pushData[1])
	if len(option) == 0 {
		return jerr.New("invalid push data for poll option, option empty")
	}
	memoPollOption = &db.MemoPollOption{
		TxHash:     txn.Hash,
		PkHash:     inputAddress.ScriptAddress(),
		PkScript:   out.PkScript,
		Option:     html_parser.EscapeWithEmojis(option),
		ParentHash: parentHash,
		BlockId:    blockId,
	}
	err = memoPollOption.Save()
	if err != nil {
		return jerr.Get("error saving memo_post for poll option", err)
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
