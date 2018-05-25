package auth

import (
	"github.com/jchavannes/jgo/web"
	"github.com/rohenaz/dtvcash/app/auth"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/rohenaz/dtvcash/app/res"
	"net/http"
)

var loginRoute = web.Route{
	Pattern: res.UrlLogin,
	Handler: func(r *web.Response) {
		if auth.IsLoggedIn(r.Session.CookieId) {
			r.SetRedirect(res.GetUrlWithBaseUrl(res.UrlIndex, r))
			return
		}
		r.Render()
	},
}

var loginSubmitRoute = web.Route{
	Pattern:     res.UrlLoginSubmit,
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

		err := auth.Login(r.Session.CookieId, username, password)
		if err != nil {
			if auth.IsBadUsernamePasswordError(err) {
				r.Error(err, http.StatusUnauthorized)
			} else {
				r.Error(err, http.StatusInternalServerError)
			}
		}
	},
}
