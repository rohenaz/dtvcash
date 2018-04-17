package server

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var newPostsRoute = web.Route{
	Pattern: res.UrlNewPosts,
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
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
		posts, err := profile.GetRecentPosts(userPkHash, uint(offset))
		if err != nil {
			r.Error(jerr.Get("error getting recent posts", err), http.StatusInternalServerError)
			return
		}
		var prevOffset int
		if offset > 25 {
			prevOffset = offset - 25
		}
		r.Helper["PrevOffset"] = prevOffset
		r.Helper["NextOffset"] = offset + 25
		r.Helper["Posts"] = posts
		r.Render()
	},
}
