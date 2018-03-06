package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var logoutRoute = web.Route{
	Pattern: UrlLogout,
	Handler: func(r *web.Response) {
		if auth.IsLoggedIn(r.Session.CookieId) {
			err := auth.Logout(r.Session.CookieId)
			if err != nil {
				r.Error(err, http.StatusInternalServerError)
				return
			}
		}
		r.SetRedirect(getUrlWithBaseUrl(UrlIndex, r))
	},
}
