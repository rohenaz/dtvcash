package auth

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/web"
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
		r.ResetOrCreateSession()
		username := r.Request.GetFormValue("username")
		password := r.Request.GetFormValue("password")

		err := auth.Login(r.Session.CookieId, username, password)
		if err != nil {
			r.Error(err, http.StatusUnauthorized)
			r.Write(err.GetDisplayMessage())
		}
	},
}
