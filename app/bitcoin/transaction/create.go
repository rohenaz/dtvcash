package transaction

import (
	"bytes"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/txscript"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/rohenaz/dtvcash/app/bitcoin/memo"
	"github.com/rohenaz/dtvcash/app/bitcoin/wallet"
	"github.com/rohenaz/dtvcash/app/db"
)

const DustMinimumOutput int64 = 546

type SpendOutputType uint

const (
	SpendOutputTypeP2PK             SpendOutputType = iota
	SpendOutputTypeReturn
	SpendOutputTypeMemoMessage
	SpendOutputTypeMemoSetName
	SpendOutputTypeMemoFollow
	SpendOutputTypeMemoUnfollow
	SpendOutputTypeMemoLike
	SpendOutputTypeMemoReply
	SpendOutputTypeMemoSetProfile
	SpendOutputTypeMemoTopicMessage
	SpendOutputTypeMemoPollQuestionSingle
	SpendOutputTypeMemoPollQuestionMulti
	SpendOutputTypeMemoPollOption
	SpendOutputTypeMemoPollVote
)

func Create(spendOuts []*db.TransactionOut, privateKey *wallet.PrivateKey, spendOutputs []SpendOutput) (*wire.MsgTx, error) {
	var txOuts []*wire.TxOut
	for _, spendOutput := range spendOutputs {
		switch spendOutput.Type {
		case SpendOutputTypeP2PK:
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_DUP).
				AddOp(txscript.OP_HASH160).
				AddData(spendOutput.Address.GetScriptAddress()).
				AddOp(txscript.OP_EQUALVERIFY).
				AddOp(txscript.OP_CHECKSIG).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating pay to addr output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeReturn:
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating op return output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(0, pkScript))
		case SpendOutputTypeMemoMessage:
			if len(spendOutput.Data) > memo.MaxPostSize {
				return nil, jerr.New("message size too large")
			}
			if len(spendOutput.Data) == 0 {
				return nil, jerr.New("empty message")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodePost}).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo message output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoSetName:
			if len(spendOutput.Data) > memo.MaxPostSize {
				return nil, jerr.New("name too large")
			}
			if len(spendOutput.Data) == 0 {
				return nil, jerr.New("empty name")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodeSetName}).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo set name output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoFollow:
			if len(spendOutput.Data) > memo.MaxPostSize {
				return nil, jerr.New("data too large")
			}
			if len(spendOutput.Data) == 0 {
				return nil, jerr.New("empty data")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodeFollow}).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo follow output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoUnfollow:
			if len(spendOutput.Data) > memo.MaxPostSize {
				return nil, jerr.New("data too large")
			}
			if len(spendOutput.Data) == 0 {
				return nil, jerr.New("empty data")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodeUnfollow}).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo unfollow output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoLike:
			if len(spendOutput.Data) > memo.MaxPostSize {
				return nil, jerr.New("data too large")
			}
			if len(spendOutput.Data) == 0 {
				return nil, jerr.New("empty data")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodeLike}).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo like output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoReply:
			if len(spendOutput.Data) > memo.MaxReplySize {
				return nil, jerr.New("data too large")
			}
			if len(spendOutput.Data) == 0 {
				return nil, jerr.New("empty data")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodeReply}).
				AddData(spendOutput.RefData).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo reply output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoSetProfile:
			if len(spendOutput.Data) > memo.MaxPostSize {
				return nil, jerr.New("profile too large")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodeSetProfile}).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo set profile output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoTopicMessage:
			if len(spendOutput.Data)+len(spendOutput.RefData) > memo.MaxTagMessageSize {
				return nil, jerr.New("data too large")
			}
			if len(spendOutput.Data) == 0 || len(spendOutput.RefData) == 0 {
				return nil, jerr.New("empty data")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodeTopicMessage}).
				AddData(spendOutput.RefData).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo tag message output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoPollQuestionSingle, SpendOutputTypeMemoPollQuestionMulti:
			var question = spendOutput.Data
			var optionCount = spendOutput.RefData
			if len(question) > memo.MaxPollQuestionSize {
				return nil, jerr.New("question size too large")
			}
			if len(question) == 0 {
				return nil, jerr.New("empty question")
			}
			if len(optionCount) == 0 {
				return nil, jerr.New("empty option count")
			}
			var pollType byte
			switch spendOutput.Type {
			case SpendOutputTypeMemoPollQuestionSingle:
				pollType = memo.CodePollTypeSingle
			case SpendOutputTypeMemoPollQuestionMulti:
				pollType = memo.CodePollTypeMulti
			default:
				return nil, jerr.New("invalid poll type")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodePollCreate}).
				AddData([]byte{pollType}).
				AddData(optionCount).
				AddData(question).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo question output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoPollOption:
			var option = spendOutput.Data
			var parentTxHash = spendOutput.RefData
			if len(option) > memo.MaxPollOptionSize {
				return nil, jerr.New("option size too large")
			}
			if len(option) == 0 {
				return nil, jerr.New("empty option")
			}
			if len(parentTxHash) != 32 {
				return nil, jerr.Newf("parent tx hash length incorrect (expected 32, got: %d)", len(parentTxHash))
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodePollOption}).
				AddData(parentTxHash).
				AddData(option).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo option output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		case SpendOutputTypeMemoPollVote:
			if len(spendOutput.Data) != 32 {
				return nil, jerr.New("invalid txn hash")
			}
			if len(spendOutput.RefData) > memo.MaxVoteCommentSize {
				return nil, jerr.New("comment data too large")
			}
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{memo.CodePrefix, memo.CodePollVote}).
				AddData(spendOutput.Data).
				AddData(spendOutput.RefData).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo poll vote output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		}
	}

	var txIns []*wire.TxIn
	var totalValue int64
	for _, spendOut := range spendOuts {
		hash, err := chainhash.NewHash(spendOut.TransactionHash)
		if err != nil {
			return nil, jerr.Get("error getting transaction hash", err)
		}
		newTxIn := wire.NewTxIn(&wire.OutPoint{
			Hash:  *hash,
			Index: uint32(spendOut.Index),
		}, nil)
		txIns = append(txIns, newTxIn)
		totalValue += spendOut.Value
	}

	var tx = &wire.MsgTx{
		Version:  wire.TxVersion,
		TxIn:     txIns,
		TxOut:    txOuts,
		LockTime: 0,
	}

	signature, err := txscript.SignatureScript(
		tx,
		0,
		spendOuts[0].PkScript,
		txscript.SigHashAll+wallet.SigHashForkID,
		privateKey.GetBtcEcPrivateKey(),
		true,
		totalValue,
	)

	if err != nil {
		return nil, jerr.Get("error signing transaction", err)
	}
	txIns[0].SignatureScript = signature

	fmt.Printf("Signature: %x\n", signature)
	writer := new(bytes.Buffer)
	err = tx.BtcEncode(writer, 1)
	if err != nil {
		return nil, jerr.Get("error encoding transaction", err)
	}
	fmt.Printf("Txn: %s\nHex: %x\n", tx.TxHash().String(), writer.Bytes())
	return tx, nil
}
