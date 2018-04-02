package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var (
	indexRoute = web.Route{
		Pattern: res.UrlIndex,
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
			privateKeys, err := db.GetKeysForUser(user.Id)
			if err != nil {
				r.Error(jerr.Get("error getting private keys for user", err), http.StatusInternalServerError)
				return
			}
			r.Helper["PrivateKeys"] = privateKeys

			recentBlock, err := db.GetRecentBlock()
			if err != nil {
				r.Error(jerr.Get("error getting recent block", err), http.StatusInternalServerError)
				return
			}

			blocks, err := db.GetBlocksInHeightRange(recentBlock.Height, recentBlock.Height - 10)
			if err != nil {
				r.Error(jerr.Get("error getting blocks in range", err), http.StatusInternalServerError)
				return
			}
			r.Helper["Blocks"] = blocks
			r.RenderTemplate("dashboard")
		},
	}
)
