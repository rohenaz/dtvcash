package node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/bloom"
	"log"
	"net"
)

/**
TODO:
- Create and broadcast transaction
 */

const BitcoinPeerAddress = "dev1.jasonc.me:8333"

var BitcoinNode Node

type Node struct {
	Peer            *peer.Peer
	NetAddress      string
	Keys            []*db.Key
	scriptAddresses []*wallet.Address
	BloomFilter     *bloom.Filter
	CheckedTxns     uint
	LastBlock       *db.Block
	SyncComplete    bool
	LastMerkleBlock *db.Block
	QueuedBlocks    map[string]*db.Block
	BlockHashes     map[string]*db.Block
	PrevBlockHashes map[string]*db.Block
}

func (n *Node) Start() {
	var p, err = peer.NewOutboundPeer(&peer.Config{
		UserAgentName:    "bch-lite-node",
		UserAgentVersion: "0.1.0",
		ChainParams:      &wallet.MainNetParams,
		DisableRelayTx:   true,
		Listeners: peer.MessageListeners{
			OnVerAck:      n.OnVerAck,
			OnHeaders:     n.OnHeaders,
			OnInv:         n.OnInv,
			OnMerkleBlock: n.OnMerkleBlock,
			OnTx:          n.OnTx,
		},
	}, n.NetAddress)
	if err != nil {
		log.Fatal(err)
	}
	n.Peer = p
	fmt.Printf("Starting bitcoin node: %s\n", n.NetAddress)
	conn, err := net.Dial("tcp", n.NetAddress)
	if err != nil {
		log.Fatal(err)
	}
	p.AssociateConnection(conn)
}

func (n *Node) SetKeys() {
	setKeys(n)
	if n.SyncComplete {
		setBloomFilters(n)
	}
}

func (n *Node) OnVerAck(p *peer.Peer, msg *wire.MsgVerAck) {
	onVerAck(n, msg)
}

func (n *Node) OnHeaders(p *peer.Peer, msg *wire.MsgHeaders) {
	onHeaders(n, msg)
}

func (n *Node) OnInv(p *peer.Peer, msg *wire.MsgInv) {
	onInv(n, msg)
}

func (n *Node) OnTx(p *peer.Peer, msg *wire.MsgTx) {
	onTx(n, msg)
}

func (n *Node) OnMerkleBlock(p *peer.Peer, msg *wire.MsgMerkleBlock) {
	onMerkleBlock(n, msg)
}
