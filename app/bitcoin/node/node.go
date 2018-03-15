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

type Node struct {
	Peer         *peer.Peer
	NetAddress   string
	Keys         []*db.Key
	Address      db.Address
	BloomFilter  *bloom.Filter
	CheckedTxns  uint
	LastBlock    *db.Block
	QueuedBlocks uint
}

func (n *Node) Start() {
	n.BloomFilter = bloom.NewFilter(2, 0, 0.0000001, wire.BloomUpdateAll)
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

func (n *Node) OnVerAck(p *peer.Peer, msg *wire.MsgVerAck) {
	block, err := db.GetRecentBlock()
	n.LastBlock = block
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}
	n.SendGetHeaders(block.GetChainhash())
}

func (n *Node) OnHeaders(p *peer.Peer, msg *wire.MsgHeaders) {
	var blocksToSave []*db.Block
	for _, header := range msg.Headers {
		block := db.ConvertMessageHeaderToBlock(header)
		lastChainHash := n.LastBlock.GetChainhash().CloneBytes()
		if bytes.Equal(block.Hash, lastChainHash) {
			// Skipping first block since we already have it
			continue
		}
		if ! bytes.Equal(block.PrevBlock, lastChainHash) {
			fmt.Println(jerr.New("block prev hash does not match!"))
			return
		}
		block.Height = n.LastBlock.Height + 1
		blocksToSave = append(blocksToSave, block)
		n.LastBlock = block
	}
	statusMsg := fmt.Sprintf("(current height: %d, time: %s)\n", n.LastBlock.Height, time.Now().Format("2006-01-02 15:04:05"))
	if len(blocksToSave) == 0 {
		fmt.Printf("done... " + statusMsg + "all caught up!\n")
		n.SetBloomFilters()
		return
	}
	err := db.SaveBlocks(blocksToSave)
	if err != nil {
		fmt.Println(jerr.Get("error saving blocks", err))
		return
	}

	fmt.Printf("Querying more... " + statusMsg)
	n.SendGetHeaders(n.LastBlock.GetChainhash())
}

func (n *Node) SetBloomFilters() {
	// Set bloom filter
	bloomFilter := bloom.NewFilter(2, 0, 0.0000001, wire.BloomUpdateAll)
	for _, key := range n.Keys {
		fmt.Printf("Adding filter for address: %s\n", key.GetAddress().GetEncoded())
		bloomFilter.Add(key.GetAddress().GetScriptAddress())
	}
	n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
	// Start checking all blocks
	recentBlock, err := db.GetRecentBlock()
	if err != nil {
		fmt.Println(jerr.Get("error getting recent block", err))
		return
	}
	n.QueueMerkleBlocks(recentBlock)
}

func (n *Node) QueueMerkleBlocks(startingBlock *db.Block) {
	blocks, err := db.GetBlocksInHeightRange(startingBlock.Height-1999, startingBlock.Height)
	if err != nil {
		fmt.Println(jerr.Get("error getting blocks in height range", err))
		return
	}
	fmt.Printf("Got blocks: %d\n", len(blocks))
	msgGetData := wire.NewMsgGetData()
	for _, block := range blocks {
		n.QueuedBlocks++
		err := msgGetData.AddInvVect(&wire.InvVect{
			Type: wire.InvTypeFilteredBlock,
			Hash: *block.GetChainhash(),
		})
		if err != nil {
			fmt.Printf("error adding invVect: %s\n", err)
			return
		}
	}
	n.Peer.QueueMessage(msgGetData, nil)
}

func (n *Node) SendGetHeaders(startingBlock *chainhash.Hash) {
	msgGetHeaders := wire.NewMsgGetHeaders()
	msgGetHeaders.BlockLocatorHashes = []*chainhash.Hash{
		startingBlock,
	}
	n.Peer.QueueMessage(msgGetHeaders, nil)
}

func (n *Node) OnInv(p *peer.Peer, m *wire.MsgInv) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(jerr.Get("error saving blocks", fmt.Errorf("Recover: %#v\n", err)))
			}
		}()
		for _, inv := range m.InvList {
			switch inv.Type {
			case wire.InvTypeBlock:
				n.SendGetHeaders(&inv.Hash)
			case wire.InvTypeTx:
				n.GetTransaction(inv.Hash)
			default:
				continue
			}

		}
	}()
}

func (n *Node) GetTransaction(txId chainhash.Hash) {
	msgGetData := wire.NewMsgGetData()
	err := msgGetData.AddInvVect(&wire.InvVect{
		Type: wire.InvTypeTx,
		Hash: txId,
	})
	if err != nil {
		fmt.Printf("error adding invVect: %s\n", err)
		return
	}
	n.Peer.QueueMessage(msgGetData, nil)
}

func (n *Node) OnTx(p *peer.Peer, msg *wire.MsgTx) {
	n.CheckedTxns++
	var scriptAddresses []wallet.Address
	for _, key := range n.Keys {
		scriptAddresses = append(scriptAddresses, key.GetAddress())
	}
	//fmt.Printf("Transaction - version: %d, locktime: %d, inputs: %d, outputs: %d\n", msg.Version, msg.LockTime, len(msg.TxIn), len(msg.TxOut))
	var found bool
	var txnInfo string
	for _, in := range msg.TxIn {
		for _, key := range n.Keys {
			if bytes.Equal(in.SignatureScript, key.GetAddress().GetScriptAddress()) {
				found = true
			}
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
			for _, key := range n.Keys {
				if bytes.Equal(address.ScriptAddress(), key.GetAddress().GetScriptAddress()) {
					found = true
				}
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
		txnInfo = "Found transaction...\n" + txnInfo
		fmt.Printf(txnInfo)
		/*var transaction = db.Transaction{
			Address: scriptAddress,
		}
		err := transaction.Save()
		if err != nil {
			fmt.Println(jerr.Get("error saving transaction", err))
			return
		}*/
	}
}

func (n *Node) OnMerkleBlock(p *peer.Peer, msg *wire.MsgMerkleBlock) {
	block := db.ConvertMessageToBlock(msg)
	n.Address.HeightChecked = block.Height
	/*if n.Address.HeightChecked % 1000 == 0 {

	}*/
	n.QueuedBlocks--
	if n.QueuedBlocks == 0 {
		fmt.Printf("Querying more... (current height checked: %d, txns: %d, time: %s)\n", n.Address.HeightChecked, n.CheckedTxns, time.Now().Format("2006-01-02 15:04:05"))
		/*if recentBlock.Height >= 25000 {
			fmt.Println("Hit max height. Stopping...")
			return
		}*/
		n.QueueMerkleBlocks(block)
	}
}
