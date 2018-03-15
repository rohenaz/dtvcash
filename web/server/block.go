package server

import (
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"strconv"
)

var blockRoute = web.Route{
	Pattern: res.UrlBlockView + "/" + paramHeight.UrlPart(),
	Handler: func(r *web.Response) {
		height := r.Request.GetUrlNamedQueryVariableUInt(paramHeight.Id)
		block, err := db.GetBlockByHeight(height)
		if err != nil {
			if ! db.IsRecordNotFoundError(err) {
				r.Error(jerr.Get("error getting block by height", err), http.StatusUnprocessableEntity)
				return
			}
			block, err = db.GetRecentBlock()
			if err != nil {
				r.Error(jerr.Get("error getting most recent block", err), http.StatusUnprocessableEntity)
				return
			}
			r.SetRedirect(res.UrlBlockView + "/" + strconv.Itoa(int(block.Height)))
			return
		}
		r.Helper["Block"] = block
		r.RenderTemplate(res.UrlBlockView)
	},
}
