package queuer

import (
	"fmt"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/cpacia/btcd/peer"
	"github.com/cpacia/btcd/wire"
	"log"
	"net"
)

const BitcoinPeerAddress = "dev1.jasonc.me:8333"

var Node QNode

type QNode struct {
	Peer *peer.Peer
}

func (n *QNode) Start() {
	var p, err = peer.NewOutboundPeer(&peer.Config{
		UserAgentName:    "bch-lite-node",
		UserAgentVersion: "0.1.0",
		ChainParams:      &wallet.MainNetParams,
		DisableRelayTx:   true,
		Listeners: peer.MessageListeners{
			OnVerAck: n.OnVerAck,
			OnReject: n.OnReject,
			OnPing:   n.OnPing,
		},
	}, BitcoinPeerAddress)
	if err != nil {
		log.Fatal(err)
	}
	n.Peer = p
	fmt.Printf("Starting bitcoin queuer node: %s\n", BitcoinPeerAddress)
	conn, err := net.Dial("tcp", BitcoinPeerAddress)
	if err != nil {
		log.Fatal(err)
	}
	p.AssociateConnection(conn)
}

func (n *QNode) KeepAlive() {
	for {
		n.Peer.WaitForDisconnect()
		fmt.Println("Queuer disconnected. Restarting.")
		n.Start()
	}
}

func (n *QNode) OnVerAck(p *peer.Peer, msg *wire.MsgVerAck) {
	fmt.Printf("VerAck: %#v\n", msg)
}

func (n *QNode) OnReject(p *peer.Peer, msg *wire.MsgReject) {
	fmt.Printf("Hash: %s\nCmd: %s\nCode: %s\nReason: %s\n", msg.Hash.String(), msg.Cmd, msg.Code.String(), msg.Reason)
}

func (n *QNode) OnPing(p *peer.Peer, msg *wire.MsgPing) {
	fmt.Printf("Received ping: %d\n", msg.Nonce)
	n.Peer.QueueMessage(wire.NewMsgPong(msg.Nonce), nil)
}

func StartAndKeepAlive() {
	Node.Start()
	fmt.Println("Keeping node alive...")
	Node.KeepAlive()
}
