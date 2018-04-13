package transaction

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/btcd/txscript"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

const DustMinimumOutput int64 = 546

type SpendOutputType uint

const (
	SpendOutputTypeP2PK        SpendOutputType = iota
	SpendOutputTypeReturn
	SpendOutputTypeMemoMessage
	SpendOutputTypeMemoSetName
	SpendOutputTypeMemoFollow
	SpendOutputTypeMemoReply
)

func Create(txOut *db.TransactionOut, privateKey *wallet.PrivateKey, spendOutputs []SpendOutput) (*wire.MsgTx, error) {
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
		case SpendOutputTypeMemoReply:
			pkScript, err := txscript.NewScriptBuilder().
				AddOp(txscript.OP_RETURN).
				AddData([]byte{0x6d, 0x00}).
				AddData(spendOutput.ReplyHash).
				AddData(spendOutput.Data).
				Script()
			if err != nil {
				return nil, jerr.Get("error creating memo message output", err)
			}
			fmt.Printf("pkScript: %x\n", pkScript)
			txOuts = append(txOuts, wire.NewTxOut(spendOutput.Amount, pkScript))
		}
	}

	hash, err := chainhash.NewHash(txOut.TransactionHash)
	if err != nil {
		return nil, jerr.Get("error getting transaction hash", err)
	}
	newTxIn := wire.NewTxIn(&wire.OutPoint{
		Hash:  *hash,
		Index: uint32(txOut.Index),
	}, nil)

	var tx = &wire.MsgTx{
		Version: wire.TxVersion,
		TxIn: []*wire.TxIn{
			newTxIn,
		},
		TxOut:    txOuts,
		LockTime: 0,
	}

	signature, err := txscript.SignatureScript(
		tx,
		0,
		txOut.PkScript,
		txscript.SigHashAll+wallet.SigHashForkID,
		privateKey.GetBtcEcPrivateKey(),
		true,
		txOut.Value,
	)

	if err != nil {
		return nil, jerr.Get("error signing transaction", err)
	}
	newTxIn.SignatureScript = signature

	fmt.Printf("Signature: %x\n", signature)
	writer := new(bytes.Buffer)
	err = tx.BtcEncode(writer, 1)
	if err != nil {
		return nil, jerr.Get("error encoding transaction", err)
	}
	fmt.Printf("Txn: %s\nHex: %x\n", tx.TxHash().String(), writer.Bytes())
	return tx, nil
}
