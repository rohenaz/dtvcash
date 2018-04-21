package main_node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func onVerAck(n *Node, msg *wire.MsgVerAck) {
	block, err := db.GetRecentBlock()
	n.LastBlock = block
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}
	n.QueuedMerkleBlocks = make(map[string]*db.Block)
	n.BlockHashes = make(map[string]*db.Block)
	sendGetHeaders(n, block.GetChainhash())
}
