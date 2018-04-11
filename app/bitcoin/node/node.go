package node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/cpacia/btcd/peer"
	"github.com/cpacia/btcd/wire"
	"log"
	"net"
)

/**
TODO:
- Make usable with public (untrusted) nodes
  - Basically just need to connect to a few and verify you're getting the same block headers
- Create and broadcast transaction
- Handle transaction between two known keys
 */

const BitcoinPeerAddress = "dev1.jasonc.me:8333"

var BitcoinNode Node

type Node struct {
	Peer               *peer.Peer
	NetAddress         string
	Keys               []*db.Key
	scriptAddresses    []*wallet.Address
	CheckedTxns        uint
	LastBlock          *db.Block
	SyncComplete       bool
	NeedsSetKeys       bool
	LastMerkleBlock    *db.Block
	QueuedMerkleBlocks map[string]*db.Block
	BlockHashes        map[string]*db.Block
	PrevBlockHashes    map[string]*db.Block
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
			OnReject:      n.OnReject,
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

func (n *Node) QueueSetKeys() {
	if len(n.QueuedMerkleBlocks) > 0 {
		n.NeedsSetKeys = true
	} else {
		n.SetKeys()
	}
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

func (n *Node) OnReject(p *peer.Peer, msg *wire.MsgReject) {
	onReject(n, msg)
}

func (n *Node) QueueMore() {
	queueMoreMerkleBlocks(n)
}
