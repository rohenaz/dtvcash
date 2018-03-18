package node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/bloom"
	"log"
	"net"
)

/**
Strategy:
- First node catches up and downloads all block headers
- After caught up, add bloom filter for all addresses
- Start with most recent block and work back to genesis, getting merkle blocks and searching transactions
- Track starting block and progress, if restarted later, only update starting block once all new blocks have been checked
  - e.g.
    Start height: 20,000
    End height: 10,000
  - Restart at 25,000, only update Start once 25,000-20,000 have been checked, then skip to 10,000 and continue
- During normal checking, update End every once in awhile (e.g. every 2,000 blocks)
- If a new address is added, start over
- Each address independently tracks progress

TODO:
  - Capture new blocks
  - Save transactions to database and display in UI
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
	fmt.Printf("Starting bitcoin block node: %s\n", n.NetAddress)
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
