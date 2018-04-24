package main_node

import (
	"fmt"
	"github.com/cpacia/btcd/wire"
)

func onReject(n *Node, msg *wire.MsgReject) {
	fmt.Printf("Hash: %s\nCmd: %s\nCode: %s\nReason: %s\n", msg.Hash.String(), msg.Cmd, msg.Code.String(), msg.Reason)
}
