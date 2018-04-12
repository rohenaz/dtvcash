package profile

import (
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var allRoute = web.Route{
	Pattern:    res.UrlProfiles,
	Handler: func(r *web.Response) {
		pkHashes, err := db.GetUniqueMemoAPkHashes()
		if err != nil {
			r.Error(jerr.Get("error getting unique pk hashes", err), http.StatusInternalServerError)
			return
		}
		var addresses []string
		for _, pkHash := range pkHashes {
			address, err := btcutil.NewAddressPubKeyHash(pkHash, &wallet.MainNetParamsOld)
			if err != nil {
				r.Error(jerr.Get("error getting pub key hash", err), http.StatusInternalServerError)
				return
			}
			addresses = append(addresses, address.String())
		}
		r.Helper["Addresses"] = addresses
		r.RenderTemplate(res.TmplProfiles)
	},
}
