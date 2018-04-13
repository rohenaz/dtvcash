package profile

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var allRoute = web.Route{
	Pattern:    res.UrlProfiles,
	Handler: func(r *web.Response) {
		var selfPkHash []byte
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
			selfPkHash = key.PkHash
		}
		profiles, err := profile.GetProfiles(selfPkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profiles", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Profiles"] = profiles
		r.RenderTemplate(res.TmplProfiles)
	},
}
