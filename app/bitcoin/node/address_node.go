package node

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
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
	Address      db.Address
	CheckedTxns  uint
	QueuedBlocks uint
}

func (n *AddressNode) Start() {
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
	}, BitcoinPeerAddress)
	if err != nil {
		log.Fatal(err)
	}
	n.Peer = p
	fmt.Printf("Starting bitcoin address node: %s\n", BitcoinPeerAddress)
	conn, err := net.Dial("tcp", BitcoinPeerAddress)
	if err != nil {
		log.Fatal(err)
	}
	p.AssociateConnection(conn)
}

func (n *AddressNode) OnVerAck(p *peer.Peer, msg *wire.MsgVerAck) {
	// Set bloom filter
	bloomFilter := bloom.NewFilter(2, 0, 0.0000001, wire.BloomUpdateAll)
	fmt.Printf("Adding filter for address: %s\n", n.Key.GetAddress().GetEncoded())
	bloomFilter.Add(n.Key.GetAddress().GetScriptAddress())
	n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
	// Start checking all blocks
	firstBlock, err := db.GetBlockByHeight(520000)
	if err != nil {
		fmt.Println(jerr.Get("error getting first block", err))
		return
	}
	n.SendGetHeaders(firstBlock.GetChainhash())
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
	n.CheckedTxns++
	scriptAddress := n.Key.GetAddress().GetScriptAddress()
	//fmt.Printf("Transaction - version: %d, locktime: %d, inputs: %d, outputs: %d\n", msg.Version, msg.LockTime, len(msg.TxIn), len(msg.TxOut))
	var found bool
	var txnInfo string
	for _, in := range msg.TxIn {
		if bytes.Equal(in.SignatureScript, scriptAddress) {
			found = true
		}
		unlockScript, err := txscript.DisasmString(in.SignatureScript)
		if err != nil {
			txnInfo = txnInfo + fmt.Sprintf("Error disassembling unlockScript: %s\n", err.Error())
		}
		txnInfo = txnInfo + fmt.Sprintf("  TxIn - Sequence: %d\n"+
			"    prevOut: %s\n"+
			"    unlockScript: %s\n",
			in.Sequence, in.PreviousOutPoint.String(), unlockScript)
	}
	for _, out := range msg.TxOut {
		lockScript, err := txscript.DisasmString(out.PkScript)
		if err != nil {
			txnInfo = txnInfo + fmt.Sprintf("Error disassembling lockScript: %s\n", err.Error())
			continue
		}
		scriptClass, addresses, sigCount, err := txscript.ExtractPkScriptAddrs(out.PkScript, &wallet.MainNetParams)
		for _, address := range addresses {
			if bytes.Equal(address.ScriptAddress(), scriptAddress) {
				found = true
			}
			txnInfo = txnInfo + fmt.Sprintf("  TxOut - value: %d\n"+
				"    lockScript: %s\n"+
				"    scriptClass: %s\n"+
				"    requiredSigs: %d\n",
				out.Value, lockScript, scriptClass, sigCount)
			txnInfo = txnInfo + fmt.Sprintf("    address: %s\n", address.String())
		}
	}
	if found {
		txnInfo = "Saving transaction...\n" + txnInfo
		fmt.Printf(txnInfo)
		var transaction = db.Transaction{
			Address: scriptAddress,
		}
		err := transaction.Save()
		if err != nil {
			fmt.Println(jerr.Get("error saving transaction", err))
			return
		}
	}
}

func (n *AddressNode) OnMerkleBlock(p *peer.Peer, msg *wire.MsgMerkleBlock) {
	block := db.ConvertMessageToBlock(msg)
	n.Address.HeightChecked++
	/*if n.Address.HeightChecked % 1000 == 0 {

	}*/
	n.QueuedBlocks--
	if n.QueuedBlocks == 0 {
		fmt.Printf("Querying more... (current height checked: %d, txns: %d, time: %s)\n", n.Address.HeightChecked, n.CheckedTxns, time.Now().Format("2006-01-02 15:04:05"))
		/*if recentBlock.Height >= 25000 {
			fmt.Println("Hit max height. Stopping...")
			return
		}*/
		n.SendGetHeaders(block.GetChainhash())
	}
}
