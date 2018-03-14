package node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/res"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/bloom"
	"github.com/jchavannes/jgo/jerr"
	"log"
	"net"
	"time"
)

type AddressNode struct {
	Peer         *peer.Peer
	Key          db.Key
	BloomFilter  *bloom.Filter
	QueuedBlocks uint
}

func (n *AddressNode) Start() {
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
			OnTx:          n.OnTx,
		},
	}, res.BitcoinPeerAddress)
	if err != nil {
		log.Fatal(err)
	}
	n.Peer = p
	conn, err := net.Dial("tcp", res.BitcoinPeerAddress)
	if err != nil {
		log.Fatal(err)
	}
	p.AssociateConnection(conn)
}

func (n *AddressNode) OnVerAck(p *peer.Peer, msg *wire.MsgVerAck) {
	// Set bloom filter
	n.BloomFilter.Add(n.Key.GetAddress().GetScriptAddress())
	n.Peer.QueueMessage(n.BloomFilter.MsgFilterLoad(), nil)
	// Start checking all blocks
	genesisBlock, err := db.GetGenesis()
	if err != nil {
		fmt.Println(jerr.Get("error getting genesis block", err))
		return
	}
	n.SendGetHeaders(genesisBlock.GetChainhash())
}

func (n *AddressNode) OnHeaders(p *peer.Peer, msg *wire.MsgHeaders) {
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

func (n *AddressNode) SendGetHeaders(startingBlock *chainhash.Hash) {
	msgGetHeaders := wire.NewMsgGetHeaders()
	msgGetHeaders.BlockLocatorHashes = []*chainhash.Hash{
		startingBlock,
	}
	n.Peer.QueueMessage(msgGetHeaders, nil)
}

func (n *AddressNode) OnTx(p *peer.Peer, msg *wire.MsgTx) {
	scriptAddress := n.Key.GetAddress().GetScriptAddress()
	fmt.Printf("Transaction - version: %d, locktime: %d, inputs: %d, outputs: %d\n", msg.Version, msg.LockTime, len(msg.TxIn), len(msg.TxOut))
	for _, in := range msg.TxIn {
		unlockScript, err := txscript.DisasmString(in.SignatureScript)
		if in.SignatureScript == scriptAddress {

		}
		if err != nil {
			fmt.Printf("Error disassembling unlockScript: %s\n", err.Error())
		}
		fmt.Printf("  TxIn - Sequence: %d\n"+
			"    prevOut: %s\n"+
			"    unlockScript: %s\n",
			in.Sequence, in.PreviousOutPoint.String(), unlockScript)
	}
	for _, out := range msg.TxOut {
		lockScript, err := txscript.DisasmString(out.PkScript)
		if err != nil {
			fmt.Printf("Error disassembling lockScript: %s\n", err.Error())
		}
		scriptClass, addresses, sigCount, err := txscript.ExtractPkScriptAddrs(out.PkScript, &wallet.MainNetParams)
		fmt.Printf("  TxOut - value: %d\n"+
			"    lockScript: %s\n"+
			"    scriptClass: %s\n"+
			"    requiredSigs: %d\n",
			out.Value, lockScript, scriptClass, sigCount)
		for _, address := range addresses {
			fmt.Printf("    address: %s\n", address.String())
		}
	}
}

func (n *AddressNode) OnMerkleBlock(p *peer.Peer, msg *wire.MsgMerkleBlock) {
	block := db.ConvertMessageToBlock(msg)
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
