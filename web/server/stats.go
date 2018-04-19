package server

import (
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var statsRoute = web.Route{
	Pattern: res.UrlStats,
	Handler: func(r *web.Response) {
		memoFollowCount, err := db.GetCountMemoFollows()
		if err != nil {
			r.Error(jerr.Get("error getting memo follow count", err), http.StatusInternalServerError)
			return
		}
		memoLikeCount, err := db.GetCountMemoLikes()
		if err != nil {
			r.Error(jerr.Get("error getting memo like count", err), http.StatusInternalServerError)
			return
		}
		memoPostCount, err := db.GetCountMemoPosts()
		if err != nil {
			r.Error(jerr.Get("error getting memo post count", err), http.StatusInternalServerError)
			return
		}
		memoSetNameCount, err := db.GetCountMemoSetName()
		if err != nil {
			r.Error(jerr.Get("error getting memo set name count", err), http.StatusInternalServerError)
			return
		}
		r.Helper["MemoFollowCount"] = memoFollowCount
		r.Helper["MemoLikeCount"] = memoLikeCount
		r.Helper["MemoPostCount"] = memoPostCount
		r.Helper["MemoSetNameCount"] = memoSetNameCount
		r.Helper["MemoTotalActionCount"] = memoFollowCount + memoLikeCount + memoPostCount + memoSetNameCount

		r.Render()
	},
}
