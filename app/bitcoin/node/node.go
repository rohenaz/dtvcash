package node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
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
}

func (n Node) Start() {
	n.BloomFilter = bloom.NewFilter(2, 0, 0.0000001, wire.BloomUpdateAll)
	var p, err = peer.NewOutboundPeer(&peer.Config{
		UserAgentName:    "bch-lite-node",
		UserAgentVersion: "0.1.0",
		ChainParams:      &wallet.MainNetParams,
		DisableRelayTx:   true,
		Listeners: peer.MessageListeners{
			OnVerAck: n.OnVerAck,
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
	// Load data
}
