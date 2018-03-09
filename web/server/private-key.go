package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var createPrivateKeySubmitRoute = web.Route{
	Pattern:     UrlCreatePrivateKeySubmit,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		var user, err = auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		var name = r.Request.GetFormValue("name")
		_, err = db.CreateNewPrivateKey(name, user.Id)
		if err != nil {
			r.Error(jerr.Get("error creating new private key", err), http.StatusInternalServerError)
		}
	},
}
