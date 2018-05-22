package poll

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/mutex"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var createRoute = web.Route{
	Pattern: res.UrlPollCreate,
	Handler: func(r *web.Response) {
		r.Render()
	},
}

var createSubmitRoute = web.Route{
	Pattern: res.UrlPollCreateSubmit,
	Handler: func(r *web.Response) {
		pollType := r.Request.GetFormValue("pollType")
		question := r.Request.GetFormValue("question")
		responses := r.Request.GetFormValueSlice("responses")
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

		var questionFee = int64(memo.MaxTxFee - memo.MaxPollQuestionSize + len([]byte(question)))
		var responseFees []int64
		var totalResponseFee int64
		for _, response := range responses {
			responseFee := int64(memo.MaxTxFee - memo.MaxPollResponseSize + len([]byte(response)))
			responseFees = append(responseFees, responseFee)
			totalResponseFee += responseFee
		}
		var minInput = questionFee + totalResponseFee + transaction.DustMinimumOutput

		mutex.Lock(key.PkHash)
		txOuts, err := db.GetSpendableTxOuts(key.PkHash, minInput)
		if err != nil {
			mutex.Unlock(key.PkHash)
			r.Error(jerr.Get("error getting spendable tx out", err), http.StatusInternalServerError)
			return
		}
		address := key.GetAddress()
		var totalValue int64
		for _, txOut := range txOuts {
			totalValue += txOut.Value
		}

		tx, err := transaction.Create(txOuts, privateKey, []transaction.SpendOutput{{
			Type:    transaction.SpendOutputTypeP2PK,
			Address: address,
			Amount:  totalValue - questionFee - totalResponseFee,
		}, {
			Type:    transaction.SpendOutputTypeMemoPollQuestion,
			Data:    []byte(question),
			RefData: []byte(pollType),
		}})
		if err != nil {
			mutex.Unlock(key.PkHash)
			r.Error(jerr.Get("error creating tx", err), http.StatusInternalServerError)
			return
		}

		mutex.Unlock(key.PkHash) // remove after testing
		fmt.Println(transaction.GetTxInfo(tx))
		/*transaction.QueueTx(tx)
		r.Write(tx.TxHash().String())*/
	},
}
