package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var (
	indexRoute = web.Route{
		Pattern: UrlIndex,
		Handler: func(r *web.Response) {
			if ! auth.IsLoggedIn(r.Session.CookieId) {
				r.Render()
				return
			}
			user, err := auth.GetSessionUser(r.Session.CookieId)
			if err != nil {
				r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
				return
			}
			privateKeys, err := db.GetPublicKeysForUser(user.Id)
			if err != nil {
				r.Error(jerr.Get("error getting private keys for user", err), http.StatusInternalServerError)
				return
			}
			r.Helper["PrivateKeys"] = privateKeys
			r.RenderTemplate("dashboard")
		},
	}
)
