package server

import (
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var testsRoute = web.Route{
	Pattern:    res.UrlTests,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		memoTests, err := db.GetMemoTests()
		if err != nil {
			r.Error(jerr.Get("error getting memo tests", err), http.StatusInternalServerError)
			return
		}
		r.Helper["MemoTests"] = memoTests
		r.Render()
	},
}
