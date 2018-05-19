package profile

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var allRoute = web.Route{
	Pattern:    res.UrlProfiles,
	Handler: func(r *web.Response) {
		r.Helper["Nav"] = "profiles"
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
