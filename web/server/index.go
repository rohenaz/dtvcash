package server

import (
	"bytes"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"strings"
)

var indexRoute = web.Route{
	Pattern: res.UrlIndex,
	Handler: func(r *web.Response) {
		r.Helper["Nav"] = "home"
		if ! auth.IsLoggedIn(r.Session.CookieId) {
			r.Render()
			return
		}
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
		r.Helper["Key"] = key

		pf, err := profile.GetProfileAndSetBalances(key.PkHash, key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profile for hash", err), http.StatusInternalServerError)
			return
		}
		err = pf.SetFollowingCount()
		if err != nil {
			r.Error(jerr.Get("error setting following count for profile", err), http.StatusInternalServerError)
			return
		}
		err = pf.SetFollowerCount()
		if err != nil {
			r.Error(jerr.Get("error setting follower count for profile", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Profile"] = pf

		err = setFeed(r, key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error setting feed", err), http.StatusInternalServerError)
			return
		}

		r.RenderTemplate("dashboard")
	},
}

var protocolRoute = web.Route{
	Pattern: res.UrlProtocol,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Protocol"
		r.Render()
	},
}

var disclaimerRoute = web.Route{
	Pattern: res.UrlDisclaimer,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - Disclaimer"
		r.Render()
	},
}

var introducingMemoRoute = web.Route{
	Pattern: res.UrlIntroducing,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Introducing Memo"
		r.Render()
	},
}

var aboutRoute = web.Route{
	Pattern: res.UrlAbout,
	Handler: func(r *web.Response) {
		r.Helper["Title"] = "Memo - About"
		r.Render()
	},
}

var feedRoute = web.Route{
	Pattern: res.UrlFeed,
	Handler: func(r *web.Response) {
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
		r.Helper["Key"] = key
		setFeed(r, key.PkHash)
		r.Render()
	},
}

func setFeed(r *web.Response, selfPkHash []byte) error {
	offset := r.Request.GetUrlParameterInt("offset")
	posts, err := profile.GetPostsFeed(selfPkHash, uint(offset))
	if err != nil {
		return jerr.Get("error getting posts for hashes", err)
	}
	err = profile.AttachLikesToPosts(posts)
	if err != nil {
		return jerr.Get("error attaching likes to posts", err)
	}
	r.Helper["PostCount"] = len(posts)
	for i := 0; i < len(posts); i++ {
		post := posts[i]
		if strings.ToLower(post.Name) == "memo" && ! bytes.Equal(post.Memo.PkHash, []byte{0x9a, 0x60, 0xa8, 0x54, 0x27, 0xc, 0x2f, 0xc2, 0xdd, 0x4d, 0xd4, 0xd3, 0xba, 0x0, 0xf2, 0x6, 0x8f, 0xd, 0x75, 0xd6}) {
			posts = append(posts[:i], posts[i+1:]...)
			i--
		}
	}
	r.Helper["Posts"] = posts
	r.Helper["Offset"] = offset

	var prevOffset int
	if offset > 25 {
		prevOffset = offset - 25
	}
	page := offset / 25 + 1
	r.Helper["Page"] = page
	r.Helper["PrevOffset"] = prevOffset
	r.Helper["NextOffset"] = offset + 25
	return nil
}

var needFundsRoute = web.Route{
	Pattern:    res.UrlNeedFunds,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
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
		r.Helper["Key"] = key
		r.Render()
	},
}

var newPostsRoute = web.Route{
	Pattern: res.UrlNewPosts,
	Handler: func(r *web.Response) {
		r.SetRedirect(res.UrlPostsNew)
	},
}
