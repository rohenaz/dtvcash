package node

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

func SendGetHeaders(n *Node, startingBlock *chainhash.Hash) {
	msgGetHeaders := wire.NewMsgGetHeaders()
	msgGetHeaders.BlockLocatorHashes = []*chainhash.Hash{
		startingBlock,
	}
	n.Peer.QueueMessage(msgGetHeaders, nil)
}

func OnHeaders(n *Node, msg *wire.MsgHeaders) {
	var blocksToSave []*db.Block
	for _, header := range msg.Headers {
		block := db.ConvertMessageHeaderToBlock(header)
		lastChainHash := n.LastBlock.GetChainhash().CloneBytes()
		if bytes.Equal(block.Hash, lastChainHash) {
			// Skipping block since we already have it
			continue
		}
		if ! bytes.Equal(block.PrevBlock, lastChainHash) {
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
	statusMsg := fmt.Sprintf("(current height: %d, time: %s)\n", n.LastBlock.Height, time.Now().Format("2006-01-02 15:04:05"))
	if len(blocksToSave) == 0 {
		if ! n.SyncComplete {
			n.SyncComplete = true
			fmt.Printf("done... " + statusMsg + "all caught up!\n")
			n.SetBloomFilters()
		}
		fmt.Printf("n.LastBlock.Height: %d\n", n.LastBlock.Height)
		return
	}
	err := db.SaveBlocks(blocksToSave)
	if err != nil {
		fmt.Println(jerr.Get("error saving blocks", err))
		return
	}

	fmt.Printf("Querying more... " + statusMsg)
	n.SendGetHeaders(n.LastBlock.GetChainhash())
}
