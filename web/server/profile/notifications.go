package profile

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/rohenaz/dtvcash/app/auth"
	"github.com/rohenaz/dtvcash/app/cache"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/rohenaz/dtvcash/app/notify"
	"github.com/rohenaz/dtvcash/app/res"
	"net/http"
)

var notificationsRoute = web.Route{
	Pattern:    res.UrlProfileNotifications,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
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
		notifications, err := notify.GetNotificationsFeed(pkHash, uint(offset))
		if err != nil {
			r.Error(jerr.Get("error getting recent notifications for user", err), http.StatusInternalServerError)
			return
		}
		lastNotificationId, err := db.GetLastNotificationId(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting last notification id from db", err), http.StatusInternalServerError)
			return
		}
		var newLastNotificationId = lastNotificationId
		for _, notification := range notifications {
			if notification.GetId() > newLastNotificationId {
				newLastNotificationId = notification.GetId()
			}
		}
		if newLastNotificationId > lastNotificationId {
			err = db.SetLastNotificationId(user.Id, newLastNotificationId)
			if err != nil {
				r.Error(jerr.Get("error setting last notification id", err), http.StatusInternalServerError)
				return
			}
			_, err = cache.GetAndSetUnreadNotificationCount(user.Id)
			if err != nil {
				r.Error(jerr.Get("error updating unread notification count in cache", err), http.StatusInternalServerError)
				return
			}
		}
		r.Helper["Notifications"] = notifications
		res.SetPageAndOffset(r, offset)
		r.RenderTemplate(res.TmplProfileNotifications)
	},
}
