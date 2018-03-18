package node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func OnVerAck(n *Node, msg *wire.MsgVerAck) {
	block, err := db.GetRecentBlock()
	n.LastBlock = block
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}
	n.QueuedBlocks = make(map[string]*db.Block)
	n.SendGetHeaders(block.GetChainhash())
}
