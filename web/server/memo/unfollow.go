package memo

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/bitcoin/transaction"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var unfollowRoute = web.Route{
	Pattern:    res.UrlMemoUnfollow + "/" + urlAddress.UrlPart(),
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		addressString := r.Request.GetUrlNamedQueryVariable(urlAddress.Id)
		address := wallet.GetAddressFromString(addressString)
		pkHash := address.GetScriptAddress()
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
		if bytes.Equal(key.PkHash, pkHash) {
			r.SetRedirect(res.GetUrlWithBaseUrl(res.UrlIndex, r))
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

		pf, err := profile.GetProfileAndSetFollowers(pkHash, key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profile for hash", err), http.StatusInternalServerError)
			return
		}
		if ! pf.CanUnFollow() {
			r.Error(jerr.New("unable to unfollow user"), http.StatusUnprocessableEntity)
			return
		}
		r.Helper["Profile"] = pf
		r.RenderTemplate(res.UrlMemoUnfollow)
	},
}

var unfollowSubmitRoute = web.Route{
	Pattern:     res.UrlMemoUnfollowSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		addressString := r.Request.GetFormValue("address")
		followAddress := wallet.GetAddressFromString(addressString)
		if followAddress.GetEncoded() != addressString {
			r.Error(jerr.New("error parsing address"), http.StatusUnprocessableEntity)
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
		var fee = int64(283 - memo.MaxPostSize + len(address.GetScriptAddress()))
		tx, err := transaction.Create(txOut, privateKey, []transaction.SpendOutput{{
			Type:    transaction.SpendOutputTypeP2PK,
			Address: address,
			Amount:  txOut.Value - fee,
		}, {
			Type: transaction.SpendOutputTypeMemoUnfollow,
			Data: followAddress.GetScriptAddress(),
		}})
		if err != nil {
			r.Error(jerr.Get("error creating tx", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(transaction.GetTxInfo(tx))
		err = transaction.QueueAndWaitForTx(tx)
		if err != nil {
			r.Error(jerr.Get("error waiting for transaction", err), http.StatusInternalServerError)
			return
		}
	},
}
