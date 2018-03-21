package server

import (
	"fmt"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/web"
)

var spendRoute = web.Route{
	Pattern:    res.UrlSpend + "/" + paramId.UrlPart(),
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		utxoId := r.Request.GetUrlNamedQueryVariableUInt(paramId.Id)
		fmt.Printf("Utxo: %d\n", utxoId)
		r.RenderTemplate(res.UrlSpend)
	},
}
