package node

import (
	"fmt"
	"github.com/cpacia/bchutil/bloom"
	"github.com/cpacia/btcd/wire"
)

func setBloomFilters(n *Node) {
	fmt.Printf("Setting bloom filter (keys: %d)...\n", len(n.Keys))
	bloomFilter := bloom.NewFilter(uint32(len(n.Keys)), 0, 0, wire.BloomUpdateNone)
	for _, key := range n.Keys {
		fmt.Printf("Adding filter for address: %s\n", key.GetAddress().GetEncoded())
		bloomFilter.Add(key.GetAddress().GetScriptAddress())
	}
	n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
	queueMoreMerkleBlocks(n)
}
