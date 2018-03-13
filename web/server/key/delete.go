package key

import (
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var deleteKeySubmitRoute = web.Route{
	Pattern:     res.UrlKeyDeleteSubmit,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		id := r.Request.GetFormValueUint("id")

		privateKey, err := db.GetKey(uint(id), user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting private key", err), http.StatusUnprocessableEntity)
			return
		}

		err = privateKey.Delete()
		if err != nil {
			r.Error(jerr.Get("error deleting private key", err), http.StatusInternalServerError)
		}
	},
}
