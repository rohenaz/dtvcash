package node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

func OnMerkleBlock(n *Node, msg *wire.MsgMerkleBlock) {
	hash := msg.Header.BlockHash().String()
	block, ok := n.QueuedBlocks[hash]
	if !ok {
		fmt.Println(jerr.New("got merkle block that wasn't queued!"))
		return
	}
	delete(n.QueuedBlocks, hash)
	n.LastMerkleBlock = block

	if len(n.QueuedBlocks) == 0 {
		if n.LastMerkleBlock.Height == 0 {
			fmt.Printf("checked entire chain!")
			return
		}
		fmt.Printf("Querying more... (current height checked: %d, txns: %d, block time: %s, time: %s)\n", n.LastMerkleBlock.Height, n.CheckedTxns, n.LastMerkleBlock.Timestamp.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
		n.QueueMerkleBlocks(n.LastMerkleBlock)
	}
}

func QueueMerkleBlocks(n *Node, startingBlock *db.Block) {
	blocks, err := db.GetBlocksInHeightRange(startingBlock.Height-1999, startingBlock.Height)
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
