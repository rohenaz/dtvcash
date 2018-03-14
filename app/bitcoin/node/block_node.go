package node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/bloom"
	"github.com/jchavannes/jgo/jerr"
	"log"
	"net"
	"time"
)

type BlockNode struct {
	Peer         *peer.Peer
	Address      string
	BloomFilter  *bloom.Filter
	QueuedBlocks uint
}

func (n *BlockNode) Start() {
	n.BloomFilter = bloom.NewFilter(2, 0, 0.0000001, wire.BloomUpdateAll)
	var p, err = peer.NewOutboundPeer(&peer.Config{
		UserAgentName:    "bch-lite-node",
		UserAgentVersion: "0.1.0",
		ChainParams:      &wallet.MainNetParams,
		DisableRelayTx:   true,
		Listeners: peer.MessageListeners{
			OnVerAck:      n.OnVerAck,
			OnHeaders:     n.OnHeaders,
			OnMerkleBlock: n.OnMerkleBlock,
		},
	}, n.Address)
	if err != nil {
		log.Fatal(err)
	}
	n.Peer = p
	conn, err := net.Dial("tcp", n.Address)
	if err != nil {
		log.Fatal(err)
	}
	p.AssociateConnection(conn)
}

func (n *BlockNode) OnVerAck(p *peer.Peer, msg *wire.MsgVerAck) {}

func (n *BlockNode) OnHeaders(p *peer.Peer, msg *wire.MsgHeaders) {
	msgGetData := wire.NewMsgGetData()
	for _, header := range msg.Headers {
		n.QueuedBlocks++
		err := msgGetData.AddInvVect(&wire.InvVect{
			Type: wire.InvTypeFilteredBlock,
			Hash: header.BlockHash(),
		})
		if err != nil {
			fmt.Printf("error adding invVect: %s\n", err)
			return
		}
	}
	n.Peer.QueueMessage(msgGetData, nil)
}

func (n *BlockNode) SendGetHeaders(startingBlock *chainhash.Hash) {
	msgGetHeaders := wire.NewMsgGetHeaders()
	msgGetHeaders.BlockLocatorHashes = []*chainhash.Hash{
		startingBlock,
	}
	n.Peer.QueueMessage(msgGetHeaders, nil)
}

func printHeader(header *wire.BlockHeader) {
	fmt.Printf("Header:\n"+
		"Bits: %d\n"+
		"Merkle root: %s\n"+
		"Nonce: %d\n"+
		"Previous block: %s\n"+
		"Timestamp: %s\n"+
		"Version: %d\n"+
		"Blockhash: %s\n",
		header.Bits, header.MerkleRoot.String(), header.Nonce, header.PrevBlock.String(),
		header.Timestamp.String(), header.Version, header.BlockHash().String())
}

func (n *BlockNode) OnMerkleBlock(p *peer.Peer, msg *wire.MsgMerkleBlock) {
	err := db.AddBlock(db.ConvertMessageToBlock(msg))
	if err != nil {
		fmt.Println(jerr.Get("error adding block", err))
	}
	n.QueuedBlocks--
	if n.QueuedBlocks == 0 {
		recentBlock, err := db.GetRecentBlock()
		if err != nil {
			fmt.Println(jerr.Get("error getting recent block", err))
			return
		}
		fmt.Printf("Querying more... (current height: %d, time: %s)\n", recentBlock.Height, time.Now().Format("2006-01-02 15:04:05"))
		/*if recentBlock.Height >= 25000 {
			fmt.Println("Hit max height. Stopping...")
			return
		}*/
		n.SendGetHeaders(recentBlock.GetChainhash())
	}
}
