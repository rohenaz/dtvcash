package node

import (
	"fmt"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func OnInv(n *Node, msg *wire.MsgInv) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(jerr.Get("error saving blocks", fmt.Errorf("Recover: %#v\n", err)))
			}
		}()
		for _, inv := range msg.InvList {
			switch inv.Type {
			case wire.InvTypeBlock:
				n.SendGetHeaders(&inv.Hash)
			case wire.InvTypeTx:
				n.GetTransaction(inv.Hash)
			default:
				continue
			}

		}
	}()
}
