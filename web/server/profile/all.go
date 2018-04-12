package profile

import (
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var allRoute = web.Route{
	Pattern:    res.UrlProfiles,
	Handler: func(r *web.Response) {
		profiles, err := profile.GetProfiles()
		if err != nil {
			r.Error(jerr.Get("error getting profiles", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Profiles"] = profiles
		r.RenderTemplate(res.TmplProfiles)
	},
}
