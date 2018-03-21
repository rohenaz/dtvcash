package key

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"strconv"
)

var viewKeyRoute = web.Route{
	Pattern:    res.UrlKeyView + "/" + urlId.UrlPart(),
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		idString := r.Request.GetUrlNamedQueryVariable(urlId.Id)
		id, err := strconv.Atoi(idString)
		if err != nil {
			r.Error(jerr.Get("error parsing id", err), http.StatusInternalServerError)
			return
		}
		privateKey, err := db.GetKey(uint(id), user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting key", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Key"] = privateKey
		transactions, err := db.GetTransactionsForKey(privateKey.Id)
		if err != nil {
			r.Error(jerr.Get("error getting transactions for key", err), http.StatusInternalServerError)
			return
		}
		var balance int64
		var balanceBCH float64
		for _, transaction := range transactions {
			balance += transaction.GetValue()
			balanceBCH += transaction.GetValueBCH()
		}
		r.Helper["Transactions"] = transactions
		r.Helper["Balance"] = balance
		r.Helper["BalanceBCH"] = balanceBCH

		r.RenderTemplate(res.UrlKeyView)
	},
}

var loadKeyRoute = web.Route{
	Pattern:     res.UrlKeyLoad,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}

		id := r.Request.GetFormValueUint("id")
		password := r.Request.GetFormValue("password")

		dbPrivateKey, err := db.GetKey(uint(id), user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting key", err), http.StatusInternalServerError)
			return
		}
		privateKey, err := dbPrivateKey.GetPrivateKey(password)
		if err != nil {
			r.Error(jerr.Get("error unlocking private key", err), http.StatusUnauthorized)
			return
		}
		r.Helper["PrivateKey"] = privateKey
		r.RenderTemplate(res.UrlKeyLoad)
	},
}
