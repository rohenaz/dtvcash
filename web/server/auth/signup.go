package auth

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/res"
	"net/http"
)

const (
	MsgErrorParsingWif         = "error parsing wif"
	MsgErrorGettingSessionUser = "error getting session user"
	MsgErrorCreatingNewPrivKey = "error creating new private key"
	MsgErrorSavingKey          = "error saving key"
	MsgErrorImportingKey       = "error importing key"
	MsgErrorUserAlreadyExists  = "user already exists"
	MsgErrorSigningUp          = "error signing up"
)

var signupRoute = web.Route{
	Pattern: res.UrlSignup,
	Handler: func(r *web.Response) {
		r.Helper["Nav"] = "signup"
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
		oldCookieId := r.Session.CookieId
		r.ResetOrCreateSession()
		db.UpdateCsrfTokenSession(oldCookieId, r.Session.CookieId)
		username := r.Request.GetFormValue("username")
		password := r.Request.GetFormValue("password")
		wif := r.Request.GetFormValue("wif")

		// Before creating account, make sure we have a valid private key
		if wif != "" {
			_, err := wallet.ImportPrivateKey(wif)
			if err != nil {
				r.Error(jerr.Get(MsgErrorParsingWif, err), http.StatusUnprocessableEntity)
				return
			}
		}

		err := auth.Signup(r.Session.CookieId, username, password)
		if auth.UserAlreadyExists(err) {
			r.Error(jerr.Get(MsgErrorUserAlreadyExists, err), http.StatusForbidden)
			return
		} else if err != nil {
			r.Error(jerr.Get(MsgErrorSigningUp, err), http.StatusUnauthorized)
			return
		}
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get(MsgErrorGettingSessionUser, err), http.StatusInternalServerError)
			return
		}
		if wif == "" {
			key, err := db.GenerateKey(username+"-generated", password, user.Id)
			if err != nil {
				r.Error(jerr.Get(MsgErrorCreatingNewPrivKey, err), http.StatusInternalServerError)
				return
			}
			recentBlock, err := db.GetRecentBlock()
			// No need to check back for a new key
			key.MaxCheck = recentBlock.Height
			err = key.Save()
			if err != nil {
				r.Error(jerr.Get(MsgErrorSavingKey, err), http.StatusInternalServerError)
				return
			}
		} else {
			_, err = db.ImportKey(username+"-imported", password, wif, user.Id)
			if err != nil {
				r.Error(jerr.Get(MsgErrorImportingKey, err), http.StatusInternalServerError)
				return
			}
		}
	},
}
