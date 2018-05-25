package memo

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/rohenaz/dtvcash/app/auth"
	"github.com/rohenaz/dtvcash/app/bitcoin/memo"
	"github.com/rohenaz/dtvcash/app/bitcoin/transaction"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/rohenaz/dtvcash/app/mutex"
	"github.com/rohenaz/dtvcash/app/res"
	"net/http"
)

var newRoute = web.Route{
	Pattern:    res.UrlMemoNew,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		r.Helper["Nav"] = "new"
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
		hasSpendableTxOut, err := db.HasSpendable(key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error getting spendable tx out", err), http.StatusInternalServerError)
			return
		}
		if ! hasSpendableTxOut {
			r.SetRedirect(res.UrlNeedFunds)
			return
		}
		r.Render()
	},
}

var newSubmitRoute = web.Route{
	Pattern:     res.UrlMemoNewSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
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
		var fee = int64(memo.MaxTxFee - memo.MaxPostSize + len([]byte(message)))
		var minInput = fee + transaction.DustMinimumOutput

		mutex.Lock(key.PkHash)
		txOut, err := db.GetSpendableTxOut(key.PkHash, minInput)
		if err != nil {
			mutex.Unlock(key.PkHash)
			r.Error(jerr.Get("error getting spendable tx out", err), http.StatusInternalServerError)
			return
		}

		tx, err := transaction.Create([]*db.TransactionOut{txOut}, privateKey, []transaction.SpendOutput{{
			Type:    transaction.SpendOutputTypeP2PK,
			Address: address,
			Amount:  txOut.Value - fee,
		}, {
			Type: transaction.SpendOutputTypeMemoMessage,
			Data: []byte(message),
		}})
		if err != nil {
			mutex.Unlock(key.PkHash)
			r.Error(jerr.Get("error creating tx", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(transaction.GetTxInfo(tx))
		transaction.QueueTx(tx)
		r.Write(tx.TxHash().String())
	},
}
