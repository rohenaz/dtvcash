package node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/memo"
	"github.com/cpacia/bchutil/bloom"
	"github.com/cpacia/btcd/wire"
)

func setBloomFilters(n *Node) {
	codes := memo.GetAllCodes()
	fmt.Printf("Setting bloom filter (keys: %d, codes: %d)...\n", len(n.Keys), len(codes))
	bloomFilter := bloom.NewFilter(uint32(len(n.Keys)+len(codes)), 0, 0, wire.BloomUpdateNone)
	for _, key := range n.Keys {
		fmt.Printf("Adding filter for address: %s\n", key.GetAddress().GetEncoded())
		bloomFilter.Add(key.GetAddress().GetScriptAddress())
	}
	for _, code := range codes {
		bloomFilter.Add(code)
	}
	n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
	queueMoreMerkleBlocks(n)
}
