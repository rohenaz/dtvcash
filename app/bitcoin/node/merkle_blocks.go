package node

import (
	"git.jasonc.me/main/memo/app/bitcoin/transaction"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

const MinCheckHeight = 520000

func onMerkleBlock(n *Node, msg *wire.MsgMerkleBlock) {
	hash := msg.Header.BlockHash().String()
	var block *db.Block
	var err error
	block, ok := n.QueuedMerkleBlocks[hash]
	if !ok {
		//jerr.Newf("got merkle block that wasn't queued! (hash: %s)", hash).Print()
		block, err = db.GetBlockByHash(msg.Header.PrevBlock)
		if err != nil {
			jerr.Get("error matching prevHash to a db block", err).Print()
			return
		}
	} else {
		delete(n.QueuedMerkleBlocks, hash)
	}

	if block != nil && block.Height != 0 {
		for _, key := range n.Keys {
			if key.MaxCheck == 0 {
				key.MaxCheck = block.Height
				key.MinCheck = block.Height
			} else if block.Height == key.MaxCheck+1 {
				key.MaxCheck = block.Height
			} else if block.Height == key.MinCheck-1 {
				key.MinCheck = block.Height
			}
		}
	}

	transactionHashes := transaction.GetTransactionsFromMerkleBlock(msg)
	for _, transactionHash := range transactionHashes {
		n.BlockHashes[transactionHash.GetTxId().String()] = block
	}

	if len(n.QueuedMerkleBlocks) == 0 {
		saveKeys(n)
		if block.Height < MinCheckHeight {
			return
		}
		queueMoreMerkleBlocks(n)
	}
}

func queueMerkleBlocks(n *Node, startingBlockHeight uint, endingBlockHeight uint) uint {
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

func queueMoreMerkleBlocks(n *Node) {
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
		numQueued += queueMerkleBlocks(n, recentBlock.Height, recentBlock.Height-2000)
	}
	// See if any new blocks need to be checked (usually after restarting)
	if numQueued < 2000 && recentBlock.Height > maxHeightChecked {
		var endQueue = maxHeightChecked + 2000 - numQueued
		if endQueue > recentBlock.Height {
			endQueue = recentBlock.Height
		}
		numQueued += queueMerkleBlocks(n, maxHeightChecked + 1, endQueue)
	}
	// Work way back to genesis
	if numQueued < 2000 && minHeightChecked > 1 {
		var endQueue = minHeightChecked - 2000 + numQueued
		if endQueue < 0 || endQueue > minHeightChecked {
			endQueue = 0
		}
		numQueued += queueMerkleBlocks(n, minHeightChecked, endQueue)
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
