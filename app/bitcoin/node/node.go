package node

import (
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/btcsuite/btcd/peer"
	"log"
	"net"
)

type Node struct {
	Peer    *peer.Peer
	Address string
}

func (n Node) Start() {
	var p, err = peer.NewOutboundPeer(&peer.Config{
		UserAgentName:    "bch-lite-node",
		UserAgentVersion: "0.1.0",
		ChainParams:      &wallet.MainNetParams,
		DisableRelayTx:   true,
		Listeners: peer.MessageListeners{
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
