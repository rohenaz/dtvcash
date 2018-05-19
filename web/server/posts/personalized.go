package posts

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
	"strings"
)

var personalizedRoute = web.Route{
	Pattern:    res.UrlPostsPersonalized,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		preHandler(r)
		offset := r.Request.GetUrlParameterInt("offset")
		timeRange := r.Request.GetUrlParameter("range")
		if timeRange == "" {
			timeRange = profile.TimeRange1Hour
		} else if ! profile.StringIsTimeRange(timeRange) {
			r.Error(jerr.New("range not valid time range"), http.StatusUnprocessableEntity)
			return
		}
		var userPkHash []byte
		var userId uint
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
			userId = user.Id
		}
		posts, err := profile.GetTopPostsNamedRange(userPkHash, uint(offset), timeRange, true)
		if err != nil {
			r.Error(jerr.Get("error getting top posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachParentToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching parent to posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachLikesToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachReputationToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching reputation to posts", err), http.StatusInternalServerError)
			return
		}
		err = profile.SetShowMediaForPosts(posts, userId)
		if err != nil {
			r.Error(jerr.Get("error setting show media for posts", err), http.StatusInternalServerError)
			return
		}
		res.SetPageAndOffset(r, offset)
		r.Helper["OffsetLink"] = fmt.Sprintf("%s?range=%s", strings.TrimLeft(res.UrlPostsPersonalized, "/"), timeRange)
		r.Helper["Posts"] = posts
		r.Helper["Range"] = timeRange
		r.Render()
	},
}
