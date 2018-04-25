package memo

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/bitcoin/transaction"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/cpacia/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var replyRoute = web.Route{
	Pattern:    res.UrlMemoReply + "/" + urlTxHash.UrlPart(),
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetUrlNamedQueryVariable(urlTxHash.Id)
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		var pkHash []byte
		if auth.IsLoggedIn(r.Session.CookieId) {
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
			pkHash = key.PkHash
		}
		post, err := profile.GetPostByTxHash(txHash.CloneBytes(), pkHash)
		if err != nil {
			r.Error(jerr.Get("error getting post", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachLikesToPosts([]*profile.Post{post})
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Post"] = post

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
		r.RenderTemplate(res.UrlMemoReply)
	},
}

var replySubmitRoute = web.Route{
	Pattern:     res.UrlMemoReplySubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetFormValue("txHash")
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
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
		transactions, err := db.GetTransactionsForPkHash(key.PkHash)
		var txOut *db.TransactionOut
		for _, txn := range transactions {
			for _, out := range txn.TxOut {
				if out.TxnIn == nil && out.Value > 1000 && bytes.Equal(out.KeyPkHash, key.PkHash) {
					txOut = out
				}
			}
		}
		if txOut == nil {
			r.Error(jerr.New("unable to find an output to spend"), http.StatusUnprocessableEntity)
			return
		}

		privateKey, err := key.GetPrivateKey(password)
		if err != nil {
			r.Error(jerr.Get("error getting private key", err), http.StatusUnauthorized)
			return
		}

		address := key.GetAddress()
		var fee = int64(283 - memo.MaxPostSize + len([]byte(message)))
		tx, err := transaction.Create(txOut, privateKey, []transaction.SpendOutput{{
			Type:    transaction.SpendOutputTypeP2PK,
			Address: address,
			Amount:  txOut.Value - fee,
		}, {
			Type:      transaction.SpendOutputTypeMemoReply,
			ReplyHash: txHash.CloneBytes(),
			Data:      []byte(message),
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
