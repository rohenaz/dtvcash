package node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

func onMerkleBlock(n *Node, msg *wire.MsgMerkleBlock) {
	hash := msg.Header.BlockHash().String()
	block, ok := n.QueuedBlocks[hash]
	if !ok {
		fmt.Println(jerr.New("got merkle block that wasn't queued!"))
		return
	}
	delete(n.QueuedBlocks, hash)
	n.LastMerkleBlock = block

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

	if block.Height%1000 == 0 {
		saveKeys(n)
	}

	if len(n.QueuedBlocks) == 0 {
		if n.LastMerkleBlock.Height == 0 {
			fmt.Printf("checked entire chain!")
			return
		}
		fmt.Printf("Querying more... (current height checked: %d, txns: %d, block time: %s, time: %s)\n", n.LastMerkleBlock.Height, n.CheckedTxns, n.LastMerkleBlock.Timestamp.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
		queueMerkleBlocks(n, n.LastMerkleBlock.Height, 0)
	}
}

func queueMerkleBlocks(n *Node, endingBlockHeight uint, startingBlockHeight uint) {
	if startingBlockHeight == 0 {
		startingBlockHeight = endingBlockHeight - 2000
	}
	blocks, err := db.GetBlocksInHeightRange(startingBlockHeight, endingBlockHeight)
	if err != nil {
		fmt.Println(jerr.Get("error getting blocks in height range", err))
		return
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
			return
		}
	}
	n.Peer.QueueMessage(msgGetData, nil)
}
