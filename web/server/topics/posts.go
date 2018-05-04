package topics

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var postsMoreRoute = web.Route{
	Pattern:    res.UrlTopicsMorePosts,
	Handler: func(r *web.Response) {
		firstPostId := r.Request.GetUrlParameterUInt("firstPostId")
		topic := r.Request.GetUrlParameter("topic")
		var userPkHash []byte
		if auth.IsLoggedIn(r.Session.CookieId) {
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
			userPkHash = key.PkHash
		}
		posts, err := profile.GetOlderPostsForTopic(topic, userPkHash, firstPostId)
		if err != nil {
			r.Error(jerr.Get("error getting posts", err), http.StatusInternalServerError)
			return
		}
		if len(posts) == 0 {
			return
		}
		if len(userPkHash) > 0 {
			err = profile.AttachReputationToPosts(posts)
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to posts", err), http.StatusInternalServerError)
				return
			}
		}
		err = profile.AttachLikesToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Posts"] = posts
		r.Helper["FirstPostId"] = posts[0].Memo.Id
		r.Render()
	},
}
