package main_node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/bitcoin/transaction"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
)

func onTx(n *Node, msg *wire.MsgTx) {
	block := findHashBlock([]map[string]*db.Block{n.BlockHashes, n.PrevBlockHashes}, msg.TxHash())
	err := transaction.SaveTransaction(msg, block)
	n.CheckedTxns++
	if err != nil {
		fmt.Println(jerr.Get("error saving transaction", err))
	}
	//fmt.Println(transaction.GetTxInfo(msg))
}

func getTransaction(n *Node, txId chainhash.Hash) {
	msgGetData := wire.NewMsgGetData()
	err := msgGetData.AddInvVect(&wire.InvVect{
		Type: wire.InvTypeTx,
		Hash: txId,
	})
	if err != nil {
		jerr.Get("error adding invVect: %s\n", err).Print()
		return
	}
	n.Peer.QueueMessage(msgGetData, nil)
}
