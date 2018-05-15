package topics

import (
	"fmt"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var createRoute = web.Route{
	Pattern: res.UrlTopicsCreate,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		preHandler(r)
		r.Render()
	},
}

var createSubmitRoute = web.Route{
	Pattern:     res.UrlTopicsCreateSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		topicName := r.Request.GetFormValue("topic")
		message := r.Request.GetFormValue("message")
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

		address := key.GetAddress()
		var fee = int64(424 - memo.MaxTagMessageSize + len([]byte(message)) + len([]byte(topicName)))
		var minInput = fee + transaction.DustMinimumOutput

		txOut, err := db.GetSpendableTxOut(key.PkHash, minInput)
		if err != nil {
			r.Error(jerr.Get("error getting spendable tx out", err), http.StatusInternalServerError)
			return
		}

		tx, err := transaction.Create(txOut, privateKey, []transaction.SpendOutput{{
			Type:    transaction.SpendOutputTypeP2PK,
			Address: address,
			Amount:  txOut.Value - fee,
		}, {
			Type:    transaction.SpendOutputTypeMemoTopicMessage,
			RefData: []byte(topicName),
			Data:    []byte(message),
		}})
		if err != nil {
			r.Error(jerr.Get("error creating tx", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(transaction.GetTxInfo(tx))
		transaction.QueueTx(tx)
		r.Write(tx.TxHash().String())
	},
}
