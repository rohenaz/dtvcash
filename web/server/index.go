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

			transactions, err := db.GetTransactionsForKey(key.Id)
			if err != nil {
				r.Error(jerr.Get("error getting transactions for key", err), http.StatusInternalServerError)
				return
			}
			var balance int64
			var balanceBCH float64
			for _, transaction := range transactions {
				balance += transaction.GetValue()
				balanceBCH += transaction.GetValueBCH()
			}
			for i, j := 0, len(transactions)-1; i < j; i, j = i+1, j-1 {
				transactions[i], transactions[j] = transactions[j], transactions[i]
			}
			r.Helper["Balance"] = balance
			r.Helper["BalanceBCH"] = balanceBCH
			r.RenderTemplate("dashboard")
		},
	}
)
