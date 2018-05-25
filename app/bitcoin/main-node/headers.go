package main_node

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
)

func sendGetHeaders(n *Node, startingBlock *chainhash.Hash) {
	msgGetHeaders := wire.NewMsgGetHeaders()
	msgGetHeaders.BlockLocatorHashes = []*chainhash.Hash{
		startingBlock,
	}
	n.Peer.QueueMessage(msgGetHeaders, nil)
}

func onHeaders(n *Node, msg *wire.MsgHeaders) {
	var lastBlock *db.Block
	for _, header := range msg.Headers {
		block := db.ConvertMessageHeaderToBlock(header)
		dbBlock, err := db.GetBlockByHash(*block.GetChainhash())
		if err != nil && ! db.IsRecordNotFoundError(err) {
			jerr.Get("error finding existing block", err).Print()
			return
		}
		if dbBlock != nil {
			// Block already exists
			continue
		}
		parentBlock, err := db.GetBlockByHash(header.PrevBlock)
		if err != nil {
			jerr.Getf(err, "error finding parent block in db (%s)", header.PrevBlock.String()).Print()
			return
		}
		block.Height = parentBlock.Height + 1
		err = block.Save()
		if err != nil {
			jerr.Get("error saving block", err).Print()
		}
		if block.Height % 10000 == 0 {
			fmt.Printf("Header scan at height: %d\n", block.Height)
		}
		lastBlock = block
	}
	if len(msg.Headers) == 0 {
		if ! n.HeaderSyncComplete {
			fmt.Println("Header sync complete")
			n.HeaderSyncComplete = true
		}
		queueBlocks(n)
		return
	}
	if lastBlock == nil {
		jerr.New("Unexpected nil lastBlock").Print()
		return
	}
	sendGetHeaders(n, lastBlock.GetChainhash())
}
