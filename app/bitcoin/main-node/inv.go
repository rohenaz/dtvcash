package main_node

import (
	"fmt"
	"github.com/memocash/memo/app/db"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func onInv(n *Node, msg *wire.MsgInv) {
	for _, inv := range msg.InvList {
		switch inv.Type {
		case wire.InvTypeBlock:
			//fmt.Printf("Got InvTypeBlock: %s\n", inv.Hash.String())
			recentBlock, err := db.GetRecentBlock()
			if err != nil {
				fmt.Println(jerr.Get("error getting recent block", err))
				return
			}
			sendGetHeaders(n, recentBlock.GetChainhash())
		case wire.InvTypeTx:
			//fmt.Printf("Got InvTypeTx: %s\n", inv.Hash.String())
			getTransaction(n, inv.Hash)
		default:
			fmt.Println("Unknown inventory")
			continue
		}

	}
}
