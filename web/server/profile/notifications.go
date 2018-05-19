package profile

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/notify"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var notificationsRoute = web.Route{
	Pattern:    res.UrlProfileNotifications,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		pkHash, err := cache.GetUserPkHash(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting address", err), http.StatusInternalServerError)
			return
		}
		notifications, err := notify.GetNotificationsFeed(pkHash, 0)
		if err != nil {
			r.Error(jerr.Get("error getting recent notifications for user", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Notifications"] = notifications
		r.RenderTemplate(res.TmplProfileNotifications)
	},
}
