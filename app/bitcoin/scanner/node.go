package scanner

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/bchutil/bloom"
	"github.com/cpacia/btcd/peer"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/main-node"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/config"
	"github.com/memocash/memo/app/db"
	"log"
	"net"
)

var Node SNode

type SNode struct {
	Peer            *peer.Peer
	BlocksQueued    int
	BlockHashes     map[string]*db.Block
	PrevBlockHashes map[string]*db.Block
	MemoTxnsFound   int
	AllTxnsFound    int
}

func (n *SNode) Start() {
	bitcoinNodeConfig := config.GetBitcoinNode()
	var p, err = peer.NewOutboundPeer(&peer.Config{
		UserAgentName:    "bch-lite-node",
		UserAgentVersion: "0.1.0",
		ChainParams:      &wallet.MainNetParams,
		DisableRelayTx:   true,
		Listeners: peer.MessageListeners{
			OnVerAck:      n.OnVerAck,
			OnReject:      n.OnReject,
			OnPing:        n.OnPing,
			OnMerkleBlock: n.OnMerkleBlock,
			OnTx:          n.OnTx,
		},
	}, bitcoinNodeConfig.GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	n.Peer = p
	fmt.Printf("Starting bitcoin queuer node: %s\n", bitcoinNodeConfig.GetConnectionString())
	conn, err := net.Dial("tcp", bitcoinNodeConfig.GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	p.AssociateConnection(conn)
}

func (n *SNode) OnVerAck(p *peer.Peer, msg *wire.MsgVerAck) {
	// set bloom filters
	setBloomFilters(n)
	// start scanning
	queueMerkleBlocks(n, main_node.MinCheckHeight)
}

func (n *SNode) OnReject(p *peer.Peer, msg *wire.MsgReject) {
	fmt.Printf("Hash: %s\nCmd: %s\nCode: %s\nReason: %s\n", msg.Hash.String(), msg.Cmd, msg.Code.String(), msg.Reason)
}

func (n *SNode) OnMerkleBlock(p *peer.Peer, msg *wire.MsgMerkleBlock) {
	onMerkleBlock(n, msg)
}

func (n *SNode) OnTx(p *peer.Peer, msg *wire.MsgTx) {
	onTx(n, msg)
}

func (n *SNode) OnPing(p *peer.Peer, msg *wire.MsgPing) {
	fmt.Printf("Received ping: %d\n", msg.Nonce)
	n.Peer.QueueMessage(wire.NewMsgPong(msg.Nonce), nil)
}

func setBloomFilters(n *SNode) {
	allKeys, err := db.GetAllKeys()
	if err != nil {
		jerr.Get("error getting keys from db", err).Print()
		return
	}
	fmt.Printf("Setting bloom filter (keys: %d)...\n", len(allKeys))
	bloomFilter := bloom.NewFilter(uint32(len(allKeys)*2), 0, 0, wire.BloomUpdateNone)
	for _, key := range allKeys {
		fmt.Printf("Adding filter for address: %s\n", key.GetAddress().GetEncoded())
		bloomFilter.Add(key.GetAddress().GetScriptAddress())
		bloomFilter.Add(key.GetPublicKey().GetSerialized())
	}
	n.Peer.QueueMessage(bloomFilter.MsgFilterLoad(), nil)
}

func queueMerkleBlocks(n *SNode, startingBlockHeight uint) error {
	if n.BlocksQueued > 0 {
		return jerr.New("blocks already queued")
	}
	//fmt.Printf("Queueing more merkle blocks (start: %d, end %d)\n", startingBlockHeight, endingBlockHeight)
	blocks, err := db.GetBlocksInHeightRange(startingBlockHeight, startingBlockHeight+999)
	if err != nil {
		return jerr.Get("error getting blocks in height range", err)
	}
	if len(blocks) == 0 {
		fmt.Printf("Finished scanning merkle blocks.\n")
		n.Peer.Disconnect()
		return nil
	}
	msgGetData := wire.NewMsgGetData()
	for _, block := range blocks {
		err := msgGetData.AddInvVect(&wire.InvVect{
			Type: wire.InvTypeFilteredBlock,
			Hash: *block.GetChainhash(),
		})
		if err != nil {
			return jerr.Get("error adding inventory vector: %s\n", err)
		}
	}
	n.PrevBlockHashes = n.BlockHashes
	n.BlockHashes = make(map[string]*db.Block)
	n.BlocksQueued += len(blocks)
	n.Peer.QueueMessage(msgGetData, nil)
	fmt.Printf("Queued %d merkle blocks\n", len(blocks))
	return nil
}

func onMerkleBlock(n *SNode, msg *wire.MsgMerkleBlock) {
	block, err := db.GetBlockByHash(msg.Header.BlockHash())
	if err != nil {
		jerr.Get("error getting block from db", err).Print()
		return
	}

	transactionHashes := transaction.GetTransactionsFromMerkleBlock(msg)
	for _, transactionHash := range transactionHashes {
		n.BlockHashes[transactionHash.GetTxId().String()] = block
	}

	n.BlocksQueued--
	if n.BlocksQueued == 0 {
		fmt.Printf("At height: %d, txns found: %d, memo txns found: %d\n", block.Height, n.AllTxnsFound, n.MemoTxnsFound)
		n.AllTxnsFound = 0
		n.MemoTxnsFound = 0
		queueMerkleBlocks(n, block.Height+1)
	}
}

func onTx(n *SNode, msg *wire.MsgTx) {
	block := findHashBlock([]map[string]*db.Block{n.BlockHashes, n.PrevBlockHashes}, msg.TxHash())
	savedTxn, savedMemo, err := transaction.ConditionallySaveTransaction(msg, block)
	if err != nil {
		jerr.Get("error conditionally saving transaction", err).Print()
	}
	if savedTxn {
		n.AllTxnsFound++
		if savedMemo {
			n.MemoTxnsFound++
		}
	}
}

func findHashBlock(blockHashes []map[string]*db.Block, hash chainhash.Hash) *db.Block {
	for _, hashMap := range blockHashes {
		for hashString, block := range hashMap {
			if hashString == hash.String() {
				return block
			}
		}
	}
	return nil
}
