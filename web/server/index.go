package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"github.com/jchavannes/jgo/web"
)

var (
	indexRoute = web.Route{
		Pattern: UrlIndex,
		Handler: func(r *web.Response) {
			if ! auth.IsLoggedIn(r.Session.CookieId) {
				r.Render()
				return
			}
			r.RenderTemplate("dashboard")
		},
	}
)
