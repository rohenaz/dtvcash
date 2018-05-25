package poll

import (
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var votesAjaxRoute = web.Route{
	Pattern: res.UrlPollVotesAjax,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetUrlParameter("txHash")
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		votes, err := profile.GetVotesForTxHash(txHash.CloneBytes())
		if err != nil {
			r.Error(jerr.Get("error getting votes for tx hash", err), http.StatusInternalServerError)
			return
		}
		r.Helper["Votes"] = votes
		r.Render()
	},
}
