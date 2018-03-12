package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"strconv"
)

var createPrivateKeySubmitRoute = web.Route{
	Pattern:     UrlCreatePrivateKeySubmit,
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
		_, err = db.CreateNewPrivateKey(name, password, user.Id)
		if err != nil {
			r.Error(jerr.Get("error creating new private key", err), http.StatusInternalServerError)
		}
	},
}

var viewKeyRoute = web.Route{
	Pattern:    UrlKeyView + "/" + urlId.UrlPart(),
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
		privateKey, err := db.GetPrivateKey(uint(id), user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting private key", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Key"] = privateKey
		r.RenderTemplate(UrlKeyView)
	},
}

var loadKeyRoute = web.Route{
	Pattern:    UrlKeyLoad,
	CsrfProtect: true,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}

		id := r.Request.GetFormValueUint("id")
		password := r.Request.GetFormValue("password")

		dbPrivateKey, err := db.GetPrivateKey(uint(id), user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting private key", err), http.StatusInternalServerError)
			return
		}
		privateKey, err := dbPrivateKey.GetPrivateKey(password)
		if err != nil {
			r.Error(jerr.Get("error unlocking private key", err), http.StatusUnauthorized)
			return
		}
		r.Helper["PrivateKey"] = privateKey
		r.RenderTemplate(UrlKeyLoad)
	},
}
