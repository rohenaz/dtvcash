package main_node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/bitcoin/transaction"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

const MinCheckHeight = 520000

func onBlock(n *Node, msg *wire.MsgBlock) {
	block, err := db.GetBlockByHash(msg.Header.BlockHash())
	if err != nil {
		jerr.Get("error getting block from db", err).Print()
		return
	}

	if block.Height == n.NodeStatus.LastBlock+1 {
		n.NodeStatus.LastBlock = block.Height
	} else if block.Height == n.NodeStatus.LastBlock-1 {
		n.NodeStatus.LastBlock = block.Height
	} else {
		fmt.Printf("Got block out of order (block.Height: %d, n.NodeStatus.LastBlock: %d)\n",
			block.Height, n.NodeStatus.LastBlock)
	}

	transactionHashes := transaction.GetTransactionsFromMerkleBlock(msg)
	for _, transactionHash := range transactionHashes {
		n.BlockHashes[transactionHash.GetTxId().String()] = block
	}

	if len(n.QueuedMerkleBlocks) == 0 {
		saveKeys(n)
		recentBlock, err := db.GetRecentBlock()
		if err != nil {
			jerr.Get("error getting recent block", err).Print()
			return
		}
		if block.Height == 0 || block.Height == recentBlock.Height {
			return
		}
		if n.NeedsSetKeys {
			n.NeedsSetKeys = false
			n.SetKeys()
		} else {
			queueMoreBlocks(n)
		}
	}
}

func queueBlocks(n *Node, startingBlockHeight uint, endingBlockHeight uint) uint {
	//fmt.Printf("Queueing more merkle blocks (start: %d, end %d)\n", startingBlockHeight, endingBlockHeight)
	blocks, err := db.GetBlocksInHeightRange(startingBlockHeight, endingBlockHeight)
	if err != nil {
		jerr.Get("error getting blocks in height range", err).Print()
		return 0
	}
	msgGetData := wire.NewMsgGetData()
	for _, block := range blocks {
		n.QueuedMerkleBlocks[block.GetChainhash().String()] = block
		err := msgGetData.AddInvVect(&wire.InvVect{
			Type: wire.InvTypeFilteredBlock,
			Hash: *block.GetChainhash(),
		})
		if err != nil {
			jerr.Get("error adding invVect: %s\n", err).Print()
			return 0
		}
	}
	n.PrevBlockHashes = n.BlockHashes
	n.BlockHashes = make(map[string]*db.Block)
	n.Peer.QueueMessage(msgGetData, nil)
	return uint(len(blocks))
}

func queueMoreBlocks(n *Node) {
	if ! n.SyncComplete {
		return
	}
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
		jerr.Get("error getting recent block", err).Print()
		return
	}

	var numQueued uint
	// Initially start at the top
	if maxHeightChecked == 0 {
		numQueued += queueBlocks(n, recentBlock.Height, recentBlock.Height-2000)
	}
	// See if any new blocks need to be checked (usually after restarting)
	if numQueued < 2000 && recentBlock.Height > maxHeightChecked {
		var endQueue = maxHeightChecked + 2000 - numQueued
		if endQueue > recentBlock.Height {
			endQueue = recentBlock.Height
		}
		numQueued += queueBlocks(n, maxHeightChecked+1, endQueue)
	}
	// Work way back to genesis
	if numQueued < 2000 && minHeightChecked > 1 {
		var endQueue = minHeightChecked - 2000 + numQueued
		if endQueue < 0 || endQueue > minHeightChecked {
			endQueue = 0
		}
		numQueued += queueBlocks(n, minHeightChecked, endQueue)
	}
}

func findHashBlock(blockHashes []map[string]*db.Block, hash chainhash.Hash) *db.Block {
	for _, hashMap := range blockHashes {
		for hashString, block := range hashMap {
			if hashString == hash.String() {
				return block
			}
		}
	}
	return nil
}
