package node

import (
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func onInv(n *Node, msg *wire.MsgInv) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(jerr.Get("error saving blocks", fmt.Errorf("Recover: %#v\n", err)))
			}
		}()
		for _, inv := range msg.InvList {
			switch inv.Type {
			case wire.InvTypeBlock:
				fmt.Printf("Got InvTypeBlock\n")
				sendGetHeaders(n, &inv.Hash)
			case wire.InvTypeTx:
				getTransaction(n, inv.Hash)
			default:
				continue
			}

		}
	}()
}
