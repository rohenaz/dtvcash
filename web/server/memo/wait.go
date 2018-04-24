package memo

import (
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"git.jasonc.me/main/memo/app/bitcoin/transaction"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"net/http"
	"strings"
)

var waitRoute = web.Route{
	Pattern:    res.UrlMemoWait + "/" + urlTxHash.UrlPart(),
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetUrlNamedQueryVariable(urlTxHash.Id)
		_, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		r.Helper["TxHash"] = txHashString
		r.RenderTemplate(res.UrlMemoWait)
	},
}

var waitSubmitRoute = web.Route{
	Pattern:     res.UrlMemoWaitSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetFormValue("txHash")
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		err = transaction.WaitForTx(txHash)
		if err != nil {
			r.Error(jerr.Get("error waiting for transaction", err), http.StatusInternalServerError)
			return
		}
		txn, err := db.GetTransactionByHash(txHash.CloneBytes())
		out, err := transaction.GetMemoOutputIfExists(txn)
		if err != nil {
			r.Error(jerr.Get("error checking for memo output", err), http.StatusInternalServerError)
			return
		}
		switch out.PkScript[3] {
		case memo.CodeFollow:
		case memo.CodeUnfollow:
			follow, err := db.GetMemoFollow(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting follow from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlProfileView + "/" + follow.Address, "/"))
		case memo.CodeLike:
			like, err := db.GetMemoLike(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting like from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlMemoPost + "/" + like.GetLikeTransactionHashString(), "/"))
		case memo.CodePost:
			post, err := db.GetMemoPost(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting post from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlMemoPost + "/" + post.GetTransactionHashString(), "/"))
		}
	},
}
