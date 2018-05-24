package poll

import (
	"fmt"
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var voteSubmitRoute = web.Route{
	Pattern: res.UrlPollVoteSubmit,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetFormValue("txHash")
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		option := r.Request.GetFormValue("option")
		tip := int64(r.Request.GetFormValueInt("tip"))
		message := r.Request.GetFormValue("message")
		fmt.Printf("Option: %s, tip: %d, message: %s, txHash: %s\n", option, tip, message, txHash.String())
		r.Write(txHash.String())
	},
}
