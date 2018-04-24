package auth

import (
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/bitcoin/node"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

const (
	MsgErrorParsingWif = "error parsing wif"
	MsgErrorGettingSessionUser = "error getting session user"
	MsgErrorCreatingNewPrivKey = "error creating new private key"
	MsgErrorSavingKey = "error saving key"
	MsgErrorImportingKey = "error importing key"
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
			r.Error(err, http.StatusForbidden)
		} else if err != nil {
			r.Error(err, http.StatusUnauthorized)
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
			}
			recentBlock, err := db.GetRecentBlock()
			// No need to check back for a new key
			key.MaxCheck = recentBlock.Height
			err = key.Save()
			if err != nil {
				r.Error(jerr.Get(MsgErrorSavingKey, err), http.StatusInternalServerError)
			}
		} else {
			_, err = db.ImportKey(username+"-imported", password, wif, user.Id)
			if err != nil {
				r.Error(jerr.Get(MsgErrorImportingKey, err), http.StatusInternalServerError)
			}
		}
		node.BitcoinNode.QueueSetKeys()
	},
}
