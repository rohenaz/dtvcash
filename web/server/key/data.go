package key

import (
	"fmt"
	"git.jasonc.me/main/memo/app/auth"
	"git.jasonc.me/main/memo/app/bitcoin/node"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var dataLoadSubmitRoute = web.Route{
	Pattern:     res.UrlKeyDataLoadSubmit,
	CsrfProtect: true,
	NeedsLogin:  true,
	Handler: func(r *web.Response) {
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		id := r.Request.GetFormValueUint("id")

		key, err := db.GetKey(uint(id), user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting key", err), http.StatusUnprocessableEntity)
			return
		}

		addressNode := node.GetAddressNode(*key)
		fmt.Printf("addressNode.CheckedTxns: %d\n", addressNode.CheckedTxns)
	},
}
