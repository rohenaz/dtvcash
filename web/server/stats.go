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
		r.Helper["MemoFollowCount"] = int64(memoFollowCount)
		r.Helper["MemoLikeCount"] = int64(memoLikeCount)
		r.Helper["MemoPostCount"] = int64(memoPostCount)
		r.Helper["MemoSetNameCount"] = int64(memoSetNameCount)
		r.Helper["MemoTotalActionCount"] = int64(memoFollowCount + memoLikeCount + memoPostCount + memoSetNameCount)

		r.Render()
	},
}
