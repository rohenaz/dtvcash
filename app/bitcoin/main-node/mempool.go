package main_node

import "github.com/cpacia/btcd/wire"

func queueMempool(n  *Node) {
	msgMemPool := wire.NewMsgMemPool()
	n.Peer.QueueMessage(msgMemPool, nil)
}
