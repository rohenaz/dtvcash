package main_node

import (
	"fmt"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/rohenaz/dtvcash/app/db"
)

func onVerAck(n *Node, msg *wire.MsgVerAck) {
	block, err := db.GetRecentBlock()
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}
	sendGetHeaders(n, block.GetChainhash())
}
