package node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func onInv(n *Node, msg *wire.MsgInv) {
	for _, inv := range msg.InvList {
		switch inv.Type {
		case wire.InvTypeBlock:
			fmt.Printf("Got InvTypeBlock: %s\n", inv.Hash.String())
			recentBlock, err := db.GetRecentBlock()
			if err != nil {
				fmt.Println(jerr.Get("error getting recent block", err))
				return
			}
			sendGetHeaders(n, recentBlock.GetChainhash())
		case wire.InvTypeTx:
			getTransaction(n, inv.Hash)
		default:
			continue
		}

	}
}
