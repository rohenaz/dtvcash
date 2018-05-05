package auth

import (
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/res"
	"github.com/jchavannes/jgo/web"
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
		r.SetRedirect(res.GetUrlWithBaseUrl(res.UrlIndex, r))
	},
}
