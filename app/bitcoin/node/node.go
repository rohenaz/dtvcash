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
)

type Node struct {
	Peer        *peer.Peer
	Address     string
	BloomFilter *bloom.Filter
	LastBlock   *wallet.Block
}

func (n *Node) Start() {
	n.BloomFilter = bloom.NewFilter(2, 0, 0.0000001, wire.BloomUpdateAll)
	var p, err = peer.NewOutboundPeer(&peer.Config{
		UserAgentName:    "bch-lite-node",
		UserAgentVersion: "0.1.0",
		ChainParams:      &wallet.MainNetParams,
		DisableRelayTx:   true,
		Listeners: peer.MessageListeners{
			OnVerAck:  n.OnVerAck,
			OnHeaders: n.OnHeaders,
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

func (n *Node) OnVerAck(p *peer.Peer, msg *wire.MsgVerAck) {
	keys, err := db.GetAllKeys()
	if err != nil {
		fmt.Println(jerr.Get("error getting keys from database", err))
		return
	}
	for _, key := range keys {
		address := key.GetPublicKey().GetAddress()
		n.BloomFilter.Add(address.GetScriptAddress())
	}
	n.Peer.QueueMessage(n.BloomFilter.MsgFilterLoad(), nil)
}

func (n *Node) OnHeaders(p *peer.Peer, msg *wire.MsgHeaders) {
	fmt.Printf("len(msg.Headers): %d\n", len(msg.Headers))
	for _, header := range msg.Headers {
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
		return
	}
}

func (n *Node) SendGetHeaders() {
	msgGetHeaders := wire.NewMsgGetHeaders()
	msgGetHeaders.BlockLocatorHashes = []*chainhash.Hash{
		wallet.GenesisBlock.Hash,
	}
	n.LastBlock = &wallet.GenesisBlock
	n.Peer.QueueMessage(msgGetHeaders, nil)
}
