package server

import (
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var blockRoute = web.Route{
	Pattern: res.UrlBlockView + "/" + paramHeight.UrlPart(),
	Handler: func(r *web.Response) {
		height := r.Request.GetUrlNamedQueryVariableUInt(paramHeight.Id)
		block, err := db.GetBlockByHeight(height)
		if err != nil {
			r.Error(jerr.Get("error getting block height", err), http.StatusUnprocessableEntity)
			return
		}
		r.Helper["Block"] = block
		r.RenderTemplate(res.UrlBlockView)
	},
}
