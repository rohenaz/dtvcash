package main_node

import (
	"fmt"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/bchutil"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

const MinCheckHeight = 525000

func onBlock(n *Node, msg *wire.MsgBlock) {
	block := bchutil.NewBlock(msg)
	dbBlock, err := db.GetBlockByHash(*block.Hash())
	if err != nil {
		jerr.Getf(err, "error getting dbBlock (%s)", block.Hash().String()).Print()
		return
	}
	var memosSaved int
	var txnsSaved int
	for _, txn := range block.Transactions() {
		savedTxn, savedMemo, err := transaction.ConditionallySaveTransaction(txn.MsgTx(), dbBlock)
		if err != nil {
			jerr.Get("error conditionally saving transaction", err).Print()
			continue
		}
		if savedTxn {
			txnsSaved++
		}
		if savedMemo {
			memosSaved++
		}
	}
	fmt.Printf("Block - height: %5d, found: %5d, saved: %5d, memos: %5d (%s)\n",
		dbBlock.Height,
		len(block.Transactions()),
		txnsSaved,
		memosSaved,
		dbBlock.Timestamp.String(),
	)
	if dbBlock.Height == n.NodeStatus.HeightChecked + 1 {
		n.NodeStatus.HeightChecked = dbBlock.Height
		err = n.NodeStatus.Save()
	}
	if err != nil {
		jerr.Get("error saving node status", err).Print()
		return
	}
	n.BlocksQueued--
	if n.BlocksQueued == 0 {
		queueBlocks(n)
	}
}

func queueBlocks(n *Node) {
	if n.BlocksQueued != 0 {
		return
	}
	if n.NodeStatus.HeightChecked < MinCheckHeight {
		n.NodeStatus.HeightChecked = MinCheckHeight
	}
	blocks, err := db.GetBlocksInHeightRange(n.NodeStatus.HeightChecked+1, n.NodeStatus.HeightChecked+2000)
	if err != nil {
		jerr.Get("error getting blocks in height range", err).Print()
		return
	}
	if len(blocks) == 0 {
		if ! n.BlocksSyncComplete {
			n.BlocksSyncComplete = true
			fmt.Println("Block sync complete")
			queueMempool(n)
		}
		return
	}
	msgGetData := wire.NewMsgGetData()
	for _, block := range blocks {
		err := msgGetData.AddInvVect(&wire.InvVect{
			Type: wire.InvTypeBlock,
			Hash: *block.GetChainhash(),
		})
		if err != nil {
			jerr.Get("error adding inventory vector: %s\n", err).Print()
			return
		}
	}
	n.Peer.QueueMessage(msgGetData, nil)
	n.BlocksQueued += len(msgGetData.InvList)
	fmt.Printf("Blocks queued: %d\n", n.BlocksQueued)
}

func getBlock(n *Node, hash chainhash.Hash) {
	getBlocks := wire.NewMsgGetBlocks(&hash)
	n.Peer.QueueMessage(getBlocks, nil)
}
