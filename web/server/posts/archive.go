package posts

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
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
		for i := 0; i < len(posts); i++ {
			post := posts[i]
			if strings.ToLower(post.Name) == "memo" && ! bytes.Equal(post.Memo.PkHash, []byte{0x9a, 0x60, 0xa8, 0x54, 0x27, 0xc, 0x2f, 0xc2, 0xdd, 0x4d, 0xd4, 0xd3, 0xba, 0x0, 0xf2, 0x6, 0x8f, 0xd, 0x75, 0xd6}) {
				posts = append(posts[:i], posts[i+1:]...)
				i--
			}
		}
		err = profile.AttachLikesToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		if len(userPkHash) > 0 {
			err = profile.AttachReputationToPosts(posts)
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to posts", err), http.StatusInternalServerError)
				return
			}
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
