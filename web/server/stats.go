package server

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/res"
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
		memoPollQuestionCount, err := db.GetCountMemoPollQuestion()
		if err != nil {
			r.Error(jerr.Get("error getting memo poll question count", err), http.StatusInternalServerError)
			return
		}
		memoPollVoteCount, err := db.GetCountMemoPollVote()
		if err != nil {
			r.Error(jerr.Get("error getting memo poll vote count", err), http.StatusInternalServerError)
			return
		}
		r.Helper["MemoFollowCount"] = int64(memoFollowCount)
		r.Helper["MemoLikeCount"] = int64(memoLikeCount)
		r.Helper["MemoPostCount"] = int64(memoPostCount)
		r.Helper["MemoSetNameCount"] = int64(memoSetNameCount)
		r.Helper["MemoPollQuestionCount"] = int64(memoPollQuestionCount)
		r.Helper["MemoPollVoteCount"] = int64(memoPollVoteCount)
		r.Helper["MemoTotalActionCount"] = int64(memoFollowCount +
			memoLikeCount +
			memoPostCount +
			memoSetNameCount +
			memoPollQuestionCount +
			memoPollVoteCount)

		r.Render()
	},
}
