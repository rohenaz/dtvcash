package auth

import (
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var signupRoute = web.Route{
	Pattern: res.UrlSignup,
	Handler: func(r *web.Response) {
		if auth.IsLoggedIn(r.Session.CookieId) {
			r.SetRedirect(res.GetUrlWithBaseUrl(res.UrlIndex, r))
			return
		}
		r.Render()
	},
}

var signupSubmitRoute = web.Route{
	Pattern:     res.UrlSignupSubmit,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		if auth.IsLoggedIn(r.Session.CookieId) {
			r.SetRedirect(res.GetUrlWithBaseUrl(res.UrlIndex, r))
			return
		}
		// Protects against some session hi-jacking attacks
		r.ResetOrCreateSession()
		username := r.Request.GetFormValue("username")
		password := r.Request.GetFormValue("password")
		wif := r.Request.GetFormValue("wif")

		// Before creating account, make sure we have a valid private key
		if wif != "" {
			_, err := wallet.ImportPrivateKey(wif)
			if err != nil {
				r.Error(jerr.Get("error parsing WIF", err), http.StatusUnprocessableEntity)
				return
			}
		}

		err := auth.Signup(r.Session.CookieId, username, password)
		if err != nil {
			r.Error(err, http.StatusUnauthorized)
		}
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		if wif == "" {
			_, err = db.GenerateKey(username+"-default", password, user.Id)
			if err != nil {
				r.Error(jerr.Get("error creating new private key", err), http.StatusInternalServerError)
			}
		} else {
			_, err = db.ImportKey(username, password, wif, user.Id)
			if err != nil {
				r.Error(jerr.Get("error importing key", err), http.StatusInternalServerError)
			}
		}
	},
}
