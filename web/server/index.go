package server

import (
	"bytes"
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"strings"
)

var indexRoute = web.Route{
	Pattern: res.UrlIndex,
	Handler: func(r *web.Response) {
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

		err = pf.SetFollowing()
		if err != nil {
			r.Error(jerr.Get("error setting following for profile", err), http.StatusInternalServerError)
			return
		}

		var pkHashes [][]byte
		for _, following := range pf.Following {
			pkHashes = append(pkHashes, following.PkHash)
		}
		posts, err := profile.GetPostsForHashes(pkHashes, key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error getting posts for hashes", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachLikesToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching likes to posts", err), http.StatusInternalServerError)
			return
		}
		for i := 0; i < len(posts); i++ {
			post := posts[i]
			if strings.ToLower(post.Name) == "memo" && ! bytes.Equal(post.Memo.PkHash, []byte{0xfe, 0x68, 0x6b, 0x9b, 0x2a, 0xb5, 0x89, 0xa3, 0xcb, 0x33, 0x68, 0xd0, 0x22, 0x11, 0xca, 0x1a, 0x9b, 0x88, 0xaa, 0x42}) {
				posts = append(posts[:i], posts[i+1:]...)
				i--
			}
		}
		r.Helper["Posts"] = posts

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
