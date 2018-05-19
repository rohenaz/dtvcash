package profile

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var settingsRoute = web.Route{
	Pattern:    res.UrlProfileSettings,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		r.RenderTemplate(res.TmplProfileSettings)
	},
}

var settingsSubmitRoute = web.Route{
	Pattern:     res.UrlProfileSettingsSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		defaultTip := r.Request.GetFormValueUint("defaultTip")
		integrations := r.Request.GetFormValue("integrations")
		theme := r.Request.GetFormValue("theme")
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		if ! db.IsValidDefaultTip(defaultTip) {
			r.Error(jerr.New("invalid default tip"), http.StatusUnprocessableEntity)
			return
		}
		if ! db.IsValidIntegrationsSetting(integrations) {
			r.Error(jerr.New("invalid default tip"), http.StatusUnprocessableEntity)
			return
		}
		if ! db.IsValidThemeSetting(theme) {
			r.Error(jerr.New("invalid default tip"), http.StatusUnprocessableEntity)
			return
		}
		userSettings, err := db.SaveSettingsForUser(user.Id, defaultTip, integrations, theme)
		if err != nil {
			r.Error(jerr.Get("error saving settings for user", err), http.StatusInternalServerError)
			return
		}
		err = cache.SetUserSettings(userSettings)
		if err != nil {
			r.Error(jerr.Get("error updating cache", err), http.StatusInternalServerError)
			return
		}
	},
}
