package server

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/bitcoin/node"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var spendRoute = web.Route{
	Pattern:    res.UrlSpend + "/" + paramId.UrlPart(),
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		utxoId := r.Request.GetUrlNamedQueryVariableUInt(paramId.Id)
		r.Helper["UtxoId"] = utxoId
		r.RenderTemplate(res.UrlSpend)
	},
}

var spendSignRoute = web.Route{
	Pattern:     res.UrlSpendSign,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		utxoId := r.Request.GetFormValueUint("id")
		txOut, err := db.GetTransactionOutputById(utxoId)
		if err != nil {
			r.Error(jerr.Get("error getting transaction output by id", err), http.StatusUnprocessableEntity)
			return
		}
		key := txOut.Transaction.Key
		address := key.GetAddress()

		password := r.Request.GetFormValue("password")

		privateKey, err := key.GetPrivateKey(password)
		if err != nil {
			r.Error(jerr.Get("error getting private key", err), http.StatusUnauthorized)
			return
		}

		pkScript, err := txscript.NewScriptBuilder().
			AddOp(txscript.OP_DUP).
			AddOp(txscript.OP_HASH160).
			AddData(address.GetScriptAddress()).
			AddOp(txscript.OP_EQUALVERIFY).
			AddOp(txscript.OP_CHECKSIG).
			Script()
		if err != nil {
			r.Error(jerr.Get("error creating pay to addr script (manual)", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("pkScript: %x\n", pkScript)

		newTxIn := wire.NewTxIn(&wire.OutPoint{
			Hash:  *txOut.Transaction.GetChainHash(),
			Index: uint32(txOut.Index),
		}, nil, nil)
		newTxOut := wire.NewTxOut(420001, pkScript)

		var unsignedTransaction = &wire.MsgTx{
			Version: wire.TxVersion,
			TxIn: []*wire.TxIn{
				newTxIn,
			},
			TxOut: []*wire.TxOut{
				newTxOut,
			},
			LockTime: 0,
		}

		fmt.Printf("unsignedTransaction: %#v\n", unsignedTransaction)

		/*var keyDb = wallet.KeyDB{
			Keys: map[string]*btcec.PrivateKey{
				address.GetEncoded(): privateKey.GetBtcEcPrivateKey(),
			},
		}*/
		writer := new(bytes.Buffer)
		unsignedTransaction.BtcEncode(writer, 0, wire.BaseEncoding)
		fmt.Printf("Unsigned: %s\nHex: %x\n", unsignedTransaction.TxHash().String(), writer.Bytes())

		signature, err := txscript.SignatureScript(
			unsignedTransaction,
			0,
			pkScript,
			txscript.SigHashAll+wallet.SigHashForkID,
			privateKey.GetBtcEcPrivateKey(),
			true,
		)

		/*signature, err := txscript.SignTxOutput(
			&wallet.MainNetParams,
			unsignedTransaction,
			0,
			pkScript,
			txscript.SigHashAll + wallet.SigHashForkID,
			keyDb,
			wallet.ScriptDb{},
			txOut.PkScript,
		)*/
		if err != nil {
			r.Error(jerr.Get("error signing transaction", err), http.StatusInternalServerError)
			return
		}
		newTxIn.SignatureScript = signature

		fmt.Printf("Signature: %x\n", signature)
		writer = new(bytes.Buffer)
		err = unsignedTransaction.BtcEncode(writer, 1, wire.WitnessEncoding)
		if err != nil {
			r.Error(jerr.Get("error encoding transaction", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Txn: %s\nHex: %x\n", unsignedTransaction.TxHash().String(), writer.Bytes())
		node.BitcoinNode.Peer.QueueMessage(unsignedTransaction, nil)
	},
}
