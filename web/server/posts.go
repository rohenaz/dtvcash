package server

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
)

var newPostsRoute = web.Route{
	Pattern: res.UrlNewPosts,
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
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
		posts, err := profile.GetRecentPosts(userPkHash, uint(offset))
		if err != nil {
			r.Error(jerr.Get("error getting top posts", err), http.StatusInternalServerError)
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
		var prevOffset int
		if offset > 25 {
			prevOffset = offset - 25
		}
		r.Helper["PrevOffset"] = prevOffset
		r.Helper["NextOffset"] = offset + 25
		r.Helper["Posts"] = posts
		r.Render()
	},
}

var topPostsRoute = web.Route{
	Pattern: res.UrlTopPosts,
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
		timeRange := r.Request.GetUrlParameter("range")
		if timeRange == "" {
			timeRange = profile.TimeRange24Hours
		} else if ! profile.StringIsTimeRange(timeRange) {
			r.Error(jerr.New("range not valid time range"), http.StatusUnprocessableEntity)
			return
		}
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
		posts, err := profile.GetTopPosts(userPkHash, uint(offset), timeRange)
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
		var prevOffset int
		if offset > 25 {
			prevOffset = offset - 25
		}
		r.Helper["OffsetLink"] = fmt.Sprintf("%s?range=%s", strings.TrimLeft(res.UrlTopPosts, "/"), timeRange)
		r.Helper["PrevOffset"] = prevOffset
		r.Helper["NextOffset"] = offset + 25
		r.Helper["Posts"] = posts
		r.Render()
	},
}
