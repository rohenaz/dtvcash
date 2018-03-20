package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var transactionRoute = web.Route{
	Pattern:    res.UrlTransactionView + "/" + paramId.UrlPart(),
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		id := r.Request.GetUrlNamedQueryVariableUInt(paramId.Id)
		transaction, err := db.GetTransactionById(id)
		if err != nil {
			r.Error(jerr.Get("error getting transaction by id", err), http.StatusInternalServerError)
			return
		}
		if transaction.Key.UserId != user.Id {
			r.Error(jerr.New("unauthorized access to transaction"), http.StatusUnauthorized)
			return
		}
		r.Helper["Transaction"] = transaction
		r.RenderTemplate(res.UrlTransactionView)
	},
}
