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
			key, err := db.GetKeyForUser(user.Id)
			if err != nil {
				r.Error(jerr.Get("error getting key for user", err), http.StatusInternalServerError)
				return
			}
			r.Helper["Key"] = key

			transactions, err := db.GetTransactionsForPkHash(key.PkHash)
			if err != nil {
				r.Error(jerr.Get("error getting transactions for key", err), http.StatusInternalServerError)
				return
			}
			var balance int64
			var balanceBCH float64
			for _, transaction := range transactions {
				for address, value := range transaction.GetValues() {
					if address == key.GetAddress().GetEncoded() {
						balance += value.GetValue()
						balanceBCH += value.GetValueBCH()
					}
				}
			}
			r.Helper["Balance"] = balance
			r.Helper["BalanceBCH"] = balanceBCH

			posts, err := db.GetPostsForPkHash(key.GetPublicKey().GetAddress().GetScriptAddress())
			if err != nil {
				r.Error(jerr.Get("error getting posts for hash", err), http.StatusInternalServerError)
				return
			}
			r.Helper["Posts"] = posts

			r.RenderTemplate("dashboard")
		},
	}
)
