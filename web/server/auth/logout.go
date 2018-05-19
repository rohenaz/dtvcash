package auth

import (
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var logoutRoute = web.Route{
	Pattern: res.UrlLogout,
	Handler: func(r *web.Response) {
		if auth.IsLoggedIn(r.Session.CookieId) {
			err := auth.Logout(r.Session.CookieId)
			if err != nil {
				r.Error(err, http.StatusInternalServerError)
				return
			}
		}
		a := r.Request.GetUrlParameter("a")
		if a == "re-login" {
			r.SetRedirect(res.GetUrlWithBaseUrl(res.UrlLogin, r))
		} else {
			r.SetRedirect(res.GetUrlWithBaseUrl(res.UrlIndex, r))
		}
	},
}
