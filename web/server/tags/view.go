package tags

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var viewRoute = web.Route{
	Pattern: res.UrlTagView + "/" + urlTagName.UrlPart(),
	Handler: func(r *web.Response) {
		preHandler(r)
		tagRaw := r.Request.GetUrlNamedQueryVariable(urlTagName.Id)
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
		tagPosts, err := profile.GetPostsForTag(tagRaw, userPkHash, 0)
		if err != nil {
			r.Error(jerr.Get("error getting tag posts from db", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Tag"] = tagRaw
		r.Helper["Posts"] = tagPosts
		r.RenderTemplate(res.TmplTagView)
	},
}
