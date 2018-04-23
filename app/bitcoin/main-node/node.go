package main_node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/cpacia/btcd/peer"
	"github.com/cpacia/btcd/wire"
	"log"
	"net"
)

const BitcoinPeerAddress = "dev1.jasonc.me:8333"

var BitcoinNode Node

type Node struct {
	Peer               *peer.Peer
	NodeStatus         *db.NodeStatus
	BlocksQueued       int
	HeaderSyncComplete bool
	BlocksSyncComplete bool
}

func Start() {
	BitcoinNode.Start()
}

func WaitForDisconnect() {
	BitcoinNode.Peer.WaitForDisconnect()
}

func (n *Node) Start() {
	nodeStatus, err := db.GetNodeStatus()
	if err != nil {
		log.Fatal(err)
	}
	n.NodeStatus = nodeStatus
	p, err := peer.NewOutboundPeer(&peer.Config{
		UserAgentName:    "bch-lite-node",
		UserAgentVersion: "0.1.0",
		ChainParams:      &wallet.MainNetParams,
		Listeners: peer.MessageListeners{
			OnVerAck:  n.OnVerAck,
			OnHeaders: n.OnHeaders,
			OnInv:     n.OnInv,
			OnBlock:   n.OnBlock,
			OnTx:      n.OnTx,
			OnReject:  n.OnReject,
			OnPing:    n.OnPing,
		},
	}, BitcoinPeerAddress)
	if err != nil {
		log.Fatal(err)
	}
	n.Peer = p
	fmt.Printf("Starting bitcoin node: %s\n", BitcoinPeerAddress)
	conn, err := net.Dial("tcp", BitcoinPeerAddress)
	if err != nil {
		log.Fatal(err)
	}
	p.AssociateConnection(conn)
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

func (n *Node) OnBlock(p *peer.Peer, msg *wire.MsgBlock, buf []byte) {
	onBlock(n, msg)
}

func (n *Node) OnReject(p *peer.Peer, msg *wire.MsgReject) {
	onReject(n, msg)
}

func (n *Node) OnPing(p *peer.Peer, msg *wire.MsgPing) {
	fmt.Printf("Received ping: %d\n", msg.Nonce)
	n.Peer.QueueMessage(wire.NewMsgPong(msg.Nonce), nil)
}
