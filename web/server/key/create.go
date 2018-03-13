package key

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var createKeyRoute = web.Route{
	Pattern:    res.UrlKeyCreate,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		r.Render()
	},
}

var createPrivateKeySubmitRoute = web.Route{
	Pattern:     res.UrlKeyCreateSubmit,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		name := r.Request.GetFormValue("name")
		password := r.Request.GetFormValue("password")
		_, err = db.GenerateKey(name, password, user.Id)
		if err != nil {
			r.Error(jerr.Get("error creating new private key", err), http.StatusInternalServerError)
		}
	},
}
