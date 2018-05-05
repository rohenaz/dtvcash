package profile

import (
	"fmt"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"strings"
)

var followersRoute = web.Route{
	Pattern:    res.UrlProfileFollowers + "/" + urlAddress.UrlPart(),
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
		addressString := r.Request.GetUrlNamedQueryVariable(urlAddress.Id)
		address := wallet.GetAddressFromString(addressString)
		pkHash := address.GetScriptAddress()
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

		pf, err := profile.GetProfile(pkHash, userPkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profile for hash", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Profile"] = pf
		followers, err := profile.GetFollowers(userPkHash, pkHash, offset)
		if err != nil {
			r.Error(jerr.Get("error setting followers for hash", err), http.StatusInternalServerError)
			return
		}
		if len(userPkHash) > 0 {
			err = profile.AttachReputationToFollowers(followers)
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to followers", err), http.StatusInternalServerError)
				return
			}
		}
		r.Helper["Followers"] = followers
		r.Helper["OffsetLink"] = fmt.Sprintf("%s/%s", strings.TrimLeft(res.UrlProfileFollowers, "/"), address.GetEncoded())
		res.SetPageAndOffset(r, offset)
		r.RenderTemplate(res.UrlProfileFollowers)
	},
}


var followingRoute = web.Route{
	Pattern:    res.UrlProfileFollowing + "/" + urlAddress.UrlPart(),
	Handler: func(r *web.Response) {
		offset := r.Request.GetUrlParameterInt("offset")
		addressString := r.Request.GetUrlNamedQueryVariable(urlAddress.Id)
		address := wallet.GetAddressFromString(addressString)
		pkHash := address.GetScriptAddress()
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

		pf, err := profile.GetProfile(pkHash, userPkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profile for hash", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Profile"] = pf
		following, err := profile.GetFollowing(userPkHash, pkHash, offset)
		if err != nil {
			r.Error(jerr.Get("error setting following for hash", err), http.StatusInternalServerError)
			return
		}
		if len(userPkHash) > 0 {
			err = profile.AttachReputationToFollowers(following)
			if err != nil {
				r.Error(jerr.Get("error attaching reputation to following", err), http.StatusInternalServerError)
				return
			}
		}
		r.Helper["Following"] = following
		r.Helper["OffsetLink"] = fmt.Sprintf("%s/%s", strings.TrimLeft(res.UrlProfileFollowing, "/"), address.GetEncoded())
		res.SetPageAndOffset(r, offset)
		r.RenderTemplate(res.UrlProfileFollowing)
	},
}
