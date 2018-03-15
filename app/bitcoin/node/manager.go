package node

import (
	"git.jasonc.me/main/memo/app/db"
)

const BitcoinPeerAddress = "dev1.jasonc.me:8333"

var BitcoinNode BlockNode

var addressNodes map[string]*AddressNode

func GetAddressNode(key db.Key) (*AddressNode) {
	id := key.GetAddress().GetEncoded()
	existingNode, ok := addressNodes[id]
	if ok {
		return existingNode
	}
	if addressNodes == nil {
		addressNodes = make(map[string]*AddressNode)
	}
	var node = AddressNode{
		Key: key,
	}
	node.Start()
	addressNodes[id] = &node
	return &node
}
