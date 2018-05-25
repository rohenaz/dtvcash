package profile

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/rohenaz/dtvcash/app/auth"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/rohenaz/dtvcash/app/html-parser"
	"github.com/rohenaz/dtvcash/app/profile"
	"github.com/rohenaz/dtvcash/app/res"
	"net/http"
	"strings"
)

var allRoute = web.Route{
	Pattern:    res.UrlProfiles,
	Handler: func(r *web.Response) {
		r.Helper["Nav"] = "profiles"
		offset := r.Request.GetUrlParameterInt("offset")
		searchString := html_parser.EscapeWithEmojis(r.Request.GetUrlParameter("s"))
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
		profiles, err := profile.GetProfiles(selfPkHash, searchString, offset)
		if err != nil {
			r.Error(jerr.Get("error getting profiles", err), http.StatusInternalServerError)
			return
		}
		res.SetPageAndOffset(r, offset)
		r.Helper["SearchString"] = searchString
		if searchString != "" {
			r.Helper["OffsetLink"] = fmt.Sprintf("%s?s=%s", strings.TrimLeft(res.UrlProfiles, "/"), searchString)
		} else {
			r.Helper["OffsetLink"] = fmt.Sprintf("%s?", res.UrlProfiles)
		}
		r.Helper["Profiles"] = profiles
		r.RenderTemplate(res.TmplProfiles)
	},
}
