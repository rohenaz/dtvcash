package node

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func sendGetHeaders(n *Node, startingBlock *chainhash.Hash) {
	msgGetHeaders := wire.NewMsgGetHeaders()
	msgGetHeaders.BlockLocatorHashes = []*chainhash.Hash{
		startingBlock,
	}
	n.Peer.QueueMessage(msgGetHeaders, nil)
}

func onHeaders(n *Node, msg *wire.MsgHeaders) {
	var blocksToSave []*db.Block
	for _, header := range msg.Headers {
		block := db.ConvertMessageHeaderToBlock(header)
		if bytes.Equal(block.Hash, n.LastBlock.Hash) {
			// Skipping block since we already have it
			continue
		}
		if ! bytes.Equal(block.PrevBlock, n.LastBlock.Hash) {
			fmt.Println(jerr.New("block prev hash does not match!"))
			fromDb, err := db.GetBlockByHash(*block.GetChainhash())
			if err != nil {
				fmt.Println(jerr.Get("error finding parent block in db", err))
				return
			}
			n.LastBlock = fromDb
		}
		block.Height = n.LastBlock.Height + 1
		blocksToSave = append(blocksToSave, block)
		n.LastBlock = block
	}
	if len(blocksToSave) == 0 {
		fmt.Printf("Current block height: %d, Block time: %s\n",
			n.LastBlock.Height, n.LastBlock.Timestamp.Format("2006-01-02 15:04:05"))
		if ! n.SyncComplete {
			n.SyncComplete = true
			setBloomFilters(n)
		} else {
			queueMoreMerkleBlocks(n)
		}
		return
	}
	err := db.SaveBlocks(blocksToSave)
	if err != nil {
		fmt.Println(jerr.Get("error saving blocks", err))
		return
	}
	sendGetHeaders(n, n.LastBlock.GetChainhash())
}
