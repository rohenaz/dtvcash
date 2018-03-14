package node

import (
	"git.jasonc.me/main/memo/app/db"
)

var addressNodes map[string]*AddressNode

func GetAddressNode(key db.Key) (*AddressNode) {
	id := key.GetAddress().GetEncoded()
	existingNode, ok := addressNodes[id]
	if ok {
		return existingNode
	}
	var node = AddressNode{
		Key: key,
	}
	node.Start()
	addressNodes[id] = &node
	return &node
}
