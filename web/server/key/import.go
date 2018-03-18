package key

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/bitcoin/node"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var importKeyRoute = web.Route{
	Pattern:    res.UrlKeyImport,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		r.Render()
	},
}

var importKeySubmitRoute = web.Route{
	Pattern:     res.UrlKeyImportSubmit,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		name := r.Request.GetFormValue("name")
		password := r.Request.GetFormValue("password")
		wif := r.Request.GetFormValue("wif")
		_, err = db.ImportKey(name, password, wif, user.Id)
		if err != nil {
			r.Error(jerr.Get("error importing key", err), http.StatusInternalServerError)
		}
		node.BitcoinNode.SetKeys()
	},
}
