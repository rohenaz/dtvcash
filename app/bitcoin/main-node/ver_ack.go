package main_node

import (
	"fmt"
	"github.com/memocash/memo/app/db"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func onVerAck(n *Node, msg *wire.MsgVerAck) {
	block, err := db.GetRecentBlock()
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}
	sendGetHeaders(n, block.GetChainhash())
}
