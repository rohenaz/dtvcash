package node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/bloom"
	"github.com/jchavannes/jgo/jerr"
)

func SetBloomFilters(n *Node) {
	// Set bloom filter
	bloomFilter := bloom.NewFilter(uint32(len(n.Keys)), 0, 0, wire.BloomUpdateNone)
	for _, key := range n.Keys {
		fmt.Printf("Adding filter for address: %s\n", key.GetAddress().GetEncoded())
		bloomFilter.Add(key.GetAddress().GetScriptAddress())
	}
	n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
	// Start checking all blocks
	recentBlock, err := db.GetRecentBlock()
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}
	n.QueueMerkleBlocks(recentBlock)
}
