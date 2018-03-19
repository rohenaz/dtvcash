package node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func onMerkleBlock(n *Node, msg *wire.MsgMerkleBlock) {
	hash := msg.Header.BlockHash().String()
	block, ok := n.QueuedBlocks[hash]
	if !ok {
		fmt.Println(jerr.New("got merkle block that wasn't queued!"))
		return
	}
	delete(n.QueuedBlocks, hash)

	if block.Height != 0 {
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

	for _, hash := range msg.Hashes {
		n.BlockHashes[hash.String()] = block
	}

	if len(n.QueuedBlocks) == 0 {
		saveKeys(n)
		if block.Height == 0 {
			fmt.Printf("checked entire chain!")
			return
		}
		/*fmt.Printf("Querying merkle blocks... last height checked: %d, txns: %d, block time: %s, time: %s\n",
			block.Height,
			n.CheckedTxns,
			block.Timestamp.Format("2006-01-02 15:04:05"),
			time.Now().Format("2006-01-02 15:04:05"),
		)*/
		queueMoreMerkleBlocks(n)
	}
}

func queueMerkleBlocks(n *Node, endingBlockHeight uint, startingBlockHeight uint) uint {
	blocks, err := db.GetBlocksInHeightRange(startingBlockHeight, endingBlockHeight)
	if err != nil {
		fmt.Println(jerr.Get("error getting blocks in height range", err))
		return 0
	}
	msgGetData := wire.NewMsgGetData()
	for i := len(blocks) - 1; i >= 0; i-- {
		block := blocks[i]
		n.QueuedBlocks[block.GetChainhash().String()] = block
		err := msgGetData.AddInvVect(&wire.InvVect{
			Type: wire.InvTypeFilteredBlock,
			Hash: *block.GetChainhash(),
		})
		if err != nil {
			fmt.Printf("error adding invVect: %s\n", err)
			return 0
		}
	}
	n.PrevBlockHashes = n.BlockHashes
	n.BlockHashes = make(map[string]*db.Block)
	n.Peer.QueueMessage(msgGetData, nil)
	return uint(len(blocks))
}

func queueMoreMerkleBlocks(n *Node) {
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

	var numQueued uint
	// Initially start at the top
	if maxHeightChecked == 0 {
		numQueued += queueMerkleBlocks(n, recentBlock.Height, recentBlock.Height-2000)
	}
	// See if any new blocks need to be checked (usually after restarting)
	if numQueued < 2000 && recentBlock.Height > maxHeightChecked {
		endQueue := recentBlock.Height - 2000 + numQueued
		if endQueue <= maxHeightChecked {
			endQueue = maxHeightChecked + 1
		}
		numQueued += queueMerkleBlocks(n, endQueue, recentBlock.Height)
	}
	// Work way back to genesis
	if numQueued < 2000 && minHeightChecked > 1 {
		var startQueue uint
		if minHeightChecked > 2000-numQueued {
			startQueue = minHeightChecked - 2000 + numQueued
		}
		numQueued += queueMerkleBlocks(n, minHeightChecked, startQueue)
	}
	if numQueued > 0 {
		//fmt.Printf("Queued %d merkle blocks...\n", numQueued)
	} else {
		//fmt.Println("Merkle blocks all caught up!")
	}
}

func findHashBlock(n *Node, hash *chainhash.Hash) *db.Block {
	for _, hashMap := range []map[string]*db.Block{n.BlockHashes, n.PrevBlockHashes} {
		for hashString, block := range hashMap {
			if hashString == hash.String() {
				return block
			}
		}
	}
	return nil
}
