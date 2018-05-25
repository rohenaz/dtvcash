package memo

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/mutex"
	"github.com/memocash/memo/app/res"
	"net/http"
	"strings"
)

var waitRoute = web.Route{
	Pattern:    res.UrlMemoWait + "/" + urlTxHash.UrlPart(),
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetUrlNamedQueryVariable(urlTxHash.Id)
		hash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		r.Helper["TxHash"] = hash.String()
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
			r.Error(jerr.Getf(err, "error waiting for transaction (%s)", txHashString), http.StatusInternalServerError)
			return
		}
		txn, err := db.GetTransactionByHashWithOutputs(txHash.CloneBytes())
		if err != nil {
			r.Error(jerr.Get("error getting transaction from db", err), http.StatusInternalServerError)
			return
		}
		out, err := transaction.GetMemoOutputIfExists(txn)
		if err != nil {
			r.Error(jerr.Get("error checking for memo output", err), http.StatusInternalServerError)
			return
		}
		mutex.Unlock(out.KeyPkHash)
		switch out.PkScript[3] {
		case memo.CodeFollow, memo.CodeUnfollow:
			follow, err := db.GetMemoFollow(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting follow from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlProfileView + "/" + follow.GetFollowAddressString(), "/"))
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
		case memo.CodeSetName:
			setName, err := db.GetMemoSetName(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting set name from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlProfileView + "/" + setName.GetAddressString(), "/"))
		case memo.CodeSetProfile:
			setProfile, err := db.GetMemoSetProfile(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting set profile from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlProfileView + "/" + setProfile.GetAddressString(), "/"))
		case memo.CodeReply:
			post, err := db.GetMemoPost(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting post from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlMemoPost + "/" + post.GetTransactionHashString(), "/"))
		case memo.CodeTopicMessage:
			post, err := db.GetMemoPost(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting post from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlTopicView + "/" + post.Topic, "/"))
		case memo.CodePollOption:
			memoPollOption, err := db.GetMemoPollOption(txHash.CloneBytes())
			if err != nil {
				r.Error(jerr.Get("error getting memo_poll_option from db", err), http.StatusInternalServerError)
				return
			}
			r.Write(strings.TrimLeft(res.UrlMemoPost + "/" + memoPollOption.GetPollTransactionHashString(), "/"))
		}
	},
}
