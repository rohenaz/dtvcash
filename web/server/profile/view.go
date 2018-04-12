package profile

import (
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
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
		posts, err := db.GetPostsForPkHash(address.GetScriptAddress())
		if err != nil {
			r.Error(jerr.Get("error getting posts for hash", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Address"] = addressString
		r.Helper["Posts"] = posts
		r.RenderTemplate(res.UrlProfileView)
	},
}
