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

var refreshKeySubmitRoute = web.Route{
	Pattern:     res.UrlKeyRefreshSubmit,
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
		recentBlock, err := db.GetRecentBlock()
		if err != nil {
			r.Error(jerr.Get("error getting recent block", err), http.StatusInternalServerError)
			return
		}
		if key.MaxCheck != recentBlock.Height {
			r.Error(jerr.Newf(
				"key is not up to date (MaxCheck: %d, recentBlock.Height: %d)",
				key.MaxCheck,
				recentBlock.Height,
			), http.StatusUnprocessableEntity)
			return
		}
		// Use same pointer as bitcoin node so updates can take effect
		for _, nodeKey := range node.BitcoinNode.Keys {
			if key.Id == nodeKey.Id {
				key = nodeKey
			}
		}

		transactions, err := db.GetTransactionsForKey(key.Id)
		if err != nil {
			r.Error(jerr.Get("error getting transactions", err), http.StatusInternalServerError)
			return
		}

		var maxBlockHeight uint
		for _, transaction := range transactions {
			if transaction.BlockId == 0 {
				err = transaction.Delete()
				if err != nil {
					r.Error(jerr.Get("error deleting transaction", err), http.StatusInternalServerError)
					return
				}
			} else if transaction.Block.Height > maxBlockHeight {
				maxBlockHeight = transaction.Block.Height
			}
		}

		if maxBlockHeight > 0 {
			key.MaxCheck = maxBlockHeight - 1
			err = key.Save()
			if err != nil {
				r.Error(jerr.Get("error saving key", err), http.StatusInternalServerError)
				return
			}
			node.BitcoinNode.QueueMore()
		}
	},
}
