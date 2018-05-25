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
	"time"
)

var archiveRoute = web.Route{
	Pattern: res.UrlPostsArchive,
	Handler: func(r *web.Response) {
		preHandler(r)
		offset := r.Request.GetUrlParameterInt("offset")
		day := r.Request.GetUrlParameter("day")
		today := time.Now().Format("2006-01-02")
		if day == "" {
			day = today
		}
		timeStart, err := time.Parse("2006-01-02", day)
		day = timeStart.Format("2006-01-02")
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
		if err != nil {
			r.Error(jerr.Get("error parsing time", err), http.StatusUnprocessableEntity)
			return
		}
		timeEnd := timeStart.Add(24 * time.Hour)
		if day == today {
			timeEnd = time.Time{}
			r.Helper["Today"] = true
		} else {
			r.Helper["Today"] = false
		}
		posts, err := profile.GetTopPosts(userPkHash, uint(offset), timeStart, timeEnd, false)
		if err != nil {
			r.Error(jerr.Get("error getting recent posts", err), http.StatusInternalServerError)
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
		err = profile.AttachPollsToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching polls to posts", err), http.StatusInternalServerError)
			return
		}
		if len(userPkHash) > 0 {
			err = profile.AttachReputationToPosts(posts)
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to posts", err), http.StatusInternalServerError)
				return
			}
		}
		err = profile.SetShowMediaForPosts(posts, userId)
		if err != nil {
			r.Error(jerr.Get("error setting show media for posts", err), http.StatusInternalServerError)
			return
		}
		res.SetPageAndOffset(r, offset)
		r.Helper["OffsetLink"] = fmt.Sprintf("%s?day=%s", strings.TrimLeft(res.UrlPostsArchive, "/"), day)
		r.Helper["PrevDay"] = timeStart.Add(-24 * time.Hour).Format("2006-01-02")
		r.Helper["NextDay"] = timeStart.Add(24 * time.Hour).Format("2006-01-02")
		r.Helper["Posts"] = posts
		r.Helper["Day"] = day
		r.Render()
	},
}
