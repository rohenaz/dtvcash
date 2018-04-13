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

var changePasswordRoute = web.Route{
	Pattern:    res.UrlKeyChangePassword,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		r.Render()
	},
}

var changePasswordSubmitRoute = web.Route{
	Pattern:     res.UrlKeyChangePasswordSubmit,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}

		oldPassword := r.Request.GetFormValue("oldPassword")
		newPassword := r.Request.GetFormValue("newPassword")

		key, err := db.GetKeyForUser(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting key", err), http.StatusInternalServerError)
			return
		}
		key = node.BitcoinNode.GetKey(key.Id)
		err = key.UpdatePassword(oldPassword, newPassword)
		if err != nil {
			r.Error(jerr.Get("error updating key password", err), http.StatusUnauthorized)
			return
		}
		err = auth.UpdatePassword(user.Id, oldPassword, newPassword)
		if err != nil {
			r.Error(jerr.Get("error updating user password", err), http.StatusUnauthorized)
			return
		}
		r.Render()
	},
}
