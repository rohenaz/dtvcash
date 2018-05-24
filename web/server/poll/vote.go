package poll

import (
	"fmt"
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/html-parser"
	"github.com/memocash/memo/app/mutex"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var voteSubmitRoute = web.Route{
	Pattern: res.UrlPollVoteSubmit,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetFormValue("txHash")
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		option := html_parser.EscapeWithEmojis(r.Request.GetFormValue("option"))
		message := r.Request.GetFormValue("message")

		memoPollOption, err := db.GetMemoPollOptionByOption(txHash.CloneBytes(), option)
		if err != nil {
			r.Error(jerr.Get("error getting memo_post", err), http.StatusInternalServerError)
			return
		}

		password := r.Request.GetFormValue("password")
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		key, err := db.GetKeyForUser(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting key for user", err), http.StatusInternalServerError)
			return
		}

		privateKey, err := key.GetPrivateKey(password)
		if err != nil {
			r.Error(jerr.Get("error getting private key", err), http.StatusUnauthorized)
			return
		}

		userAddress := key.GetAddress()
		postAddress := memoPollOption.GetAddress()

		var tx *wire.MsgTx

		var fee = int64(memo.MaxTxFee - memo.MaxVoteCommentSize + len(message))
		tip := int64(r.Request.GetFormValueInt("tip"))
		var minInput = fee + transaction.DustMinimumOutput + tip

		transactions := []transaction.SpendOutput{{
			Type:    transaction.SpendOutputTypeMemoPollVote,
			Data:    memoPollOption.TxHash,
			RefData: []byte(message),
		}}

		mutex.Lock(key.PkHash)
		txOut, err := db.GetSpendableTxOut(key.PkHash, minInput)
		if err != nil {
			mutex.Unlock(key.PkHash)
			r.Error(jerr.Get("error getting spendable tx out", err), http.StatusPaymentRequired)
			return
		}
		remaining := txOut.Value

		if tip != 0 {
			if tip < transaction.DustMinimumOutput {
				mutex.Unlock(key.PkHash)
				r.Error(jerr.Get("error tip not above dust limit", err), http.StatusUnprocessableEntity)
				return
			}
			if tip > 1e8 {
				mutex.Unlock(key.PkHash)
				r.Error(jerr.Get("error trying to tip too much", err), http.StatusUnprocessableEntity)
				return
			}
			transactions = append(transactions, transaction.SpendOutput{
				Type:    transaction.SpendOutputTypeP2PK,
				Address: postAddress,
				Amount:  tip,
			})
			remaining -= tip
			if remaining < transaction.DustMinimumOutput {
				mutex.Unlock(key.PkHash)
				r.Error(jerr.New("not enough funds"), http.StatusUnprocessableEntity)
				return
			}
			fee += memo.AdditionalOutputFee
		}
		transactions = append(transactions, transaction.SpendOutput{
			Type:    transaction.SpendOutputTypeP2PK,
			Address: userAddress,
			Amount:  remaining - fee,
		})
		tx, err = transaction.Create([]*db.TransactionOut{txOut}, privateKey, transactions)
		if err != nil {
			mutex.Unlock(key.PkHash)
			r.Error(jerr.Get("error creating tx", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(transaction.GetTxInfo(tx))
		//transaction.QueueTx(tx)
		r.Write(tx.TxHash().String())
	},
}
