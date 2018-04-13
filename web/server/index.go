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

var indexRoute = web.Route{
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

		pf, err := profile.GetProfileAndSetBalances(key.PkHash, key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profile for hash", err), http.StatusInternalServerError)
			return
		}
		err = pf.SetFollowing()
		if err != nil {
			r.Error(jerr.Get("error setting following for profile", err), http.StatusInternalServerError)
			return
		}
		err = pf.SetFollowers()
		if err != nil {
			r.Error(jerr.Get("error setting followers for profile", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Profile"] = pf

		var pkHashes [][]byte
		for _, following := range pf.Following {
			pkHashes = append(pkHashes, following.PkHash)
		}
		posts, err := profile.GetPostsForHashes(pkHashes, key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error getting posts for hashes", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachLikesToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Posts"] = posts

		r.RenderTemplate("dashboard")
	},
}

var protocolRoute = web.Route{
	Pattern: res.UrlProtocol,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Protocol"
		r.Render()
	},
}

var disclaimerRoute = web.Route{
	Pattern: res.UrlDisclaimer,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Disclaimer"
		r.Render()
	},
}
