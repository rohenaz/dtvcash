package key

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var viewKeyRoute = web.Route{
	Pattern:    res.UrlKeyExport,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
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
		r.Helper["Key"] = key
		r.RenderTemplate(res.UrlKeyExport)
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
