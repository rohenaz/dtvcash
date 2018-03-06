package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"github.com/jchavannes/jgo/web"
	"net/http"
)


var loginRoute = web.Route{
	Pattern: UrlLogin,
	Handler: func(r *web.Response) {
		r.Render()
	},
}

var loginSubmitRoute = web.Route{
	Pattern:     UrlLoginSubmit,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		if auth.IsLoggedIn(r.Session.CookieId) {
			r.SetRedirect(getUrlWithBaseUrl(UrlIndex, r))
			return
		}
		// Protects against some session hi-jacking attacks
		r.ResetOrCreateSession()
		username := r.Request.GetFormValue("username")
		password := r.Request.GetFormValue("password")

		err := auth.Login(r.Session.CookieId, username, password)
		if err != nil {
			r.Error(err, http.StatusUnauthorized)
		}
	},
}
