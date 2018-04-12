package profile

import (
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/profile"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var viewRoute = web.Route{
	Pattern:    res.UrlProfileView + "/" + urlAddress.UrlPart(),
	Handler: func(r *web.Response) {
		addressString := r.Request.GetUrlNamedQueryVariable(urlAddress.Id)
		address := wallet.GetAddressFromString(addressString)
		pkHash := address.GetScriptAddress()

		posts, err := db.GetPostsForPkHash(pkHash)
		if err != nil {
			r.Error(jerr.Get("error getting posts for hash", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Posts"] = posts

		pf, err := profile.GetProfile(pkHash)
		if err != nil {
			r.Error(jerr.Get("error getting profile for hash", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Profile"] = pf

		r.RenderTemplate(res.UrlProfileView)
	},
}
