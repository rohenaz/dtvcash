package main_node

import "github.com/jchavannes/btcd/wire"

func queueMempool(n  *Node) {
	msgMemPool := wire.NewMsgMemPool()
	n.Peer.QueueMessage(msgMemPool, nil)
}
