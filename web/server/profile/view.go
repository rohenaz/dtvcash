package profile

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var viewRoute = web.Route{
	Pattern:    res.UrlProfileView + "/" + urlAddress.UrlPart(),
	Handler: func(r *web.Response) {
		addressString := r.Request.GetUrlNamedQueryVariable(urlAddress.Id)
		address := wallet.GetAddressFromString(addressString)
		pkHash := address.GetScriptAddress()
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

		offset := r.Request.GetUrlParameterInt("offset")
		posts, err := profile.GetPostsForHash(pkHash, userPkHash, uint(offset))
		if err != nil {
			r.Error(jerr.Get("error getting posts for hash", err), http.StatusInternalServerError)
			return
		}
		err = profile.AttachParentToPosts(posts)
		if err != nil {
			r.Error(jerr.Get("error attaching parent to post", err), http.StatusInternalServerError)
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
		err = profile.SetShowMediaForPosts(posts, userId)
		if err != nil {
			r.Error(jerr.Get("error setting show media for posts", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Posts"] = posts

		pf, err := profile.GetProfile(pkHash, userPkHash)
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
		if len(userPkHash) > 0 {
			err = pf.SetReputation()
			if err != nil {
				r.Error(jerr.Get("error getting reputation", err), http.StatusInternalServerError)
				return
			}
			err = pf.SetCanFollow()
			if err != nil {
				r.Error(jerr.Get("error setting can follow for profile", err), http.StatusInternalServerError)
				return
			}
		}
		err = pf.SetQr()
		if err != nil {
			r.Error(jerr.Get("error creating qr", err), http.StatusInternalServerError)
			return
		}

		r.Helper["Profile"] = pf

		memoLikes, err := profile.GetLikesForPkHash(pkHash)
		r.Helper["Likes"] = memoLikes
		r.Helper["Title"] = fmt.Sprintf("Memo - %s's Profile", pf.Name)
		res.SetPageAndOffset(r, offset)
		r.RenderTemplate(res.UrlProfileView)
	},
}
