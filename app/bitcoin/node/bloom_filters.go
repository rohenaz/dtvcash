package node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/bloom"
	"github.com/jchavannes/jgo/jerr"
)

func setBloomFilters(n *Node) {
	// Set bloom filter
	bloomFilter := bloom.NewFilter(uint32(len(n.Keys)), 0, 0, wire.BloomUpdateNone)
	for _, key := range n.Keys {
		fmt.Printf("Adding filter for address: %s\n", key.GetAddress().GetEncoded())
		bloomFilter.Add(key.GetAddress().GetScriptAddress())
	}
	n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
	var minHeightChecked uint
	for _, key := range n.Keys {
		if key.MinCheck == 0 {
			break
		}
		if key.MinCheck > minHeightChecked {
			minHeightChecked = key.MinCheck
		}
	}
	var maxHeightChecked uint
	for _, key := range n.Keys {
		if key.MaxCheck == 0 {
			break
		}
		if maxHeightChecked == 0 || key.MaxCheck < maxHeightChecked {
			maxHeightChecked = key.MaxCheck
		}
	}
	recentBlock, err := db.GetRecentBlock()
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}

	if maxHeightChecked > 0 && recentBlock.Height > maxHeightChecked {
		if recentBlock.Height > maxHeightChecked + 2000 {
			queueMerkleBlocks(n, recentBlock.Height, 0)
			return
		} else {
			queueMerkleBlocks(n, maxHeightChecked, recentBlock.Height)
		}
	}
	var minStartBlock *db.Block
	if minHeightChecked == 0 {
		minStartBlock = recentBlock
	} else {
		minStartBlock, err = db.GetBlockByHeight(minHeightChecked)
		if err != nil {
			fmt.Println(jerr.Get("error getting starting block", err))
			return
		}
	}
	queueMerkleBlocks(n, minStartBlock.Height, 0)
}
