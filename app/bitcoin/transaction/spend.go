package transaction

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/peer"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	btcdPeer "github.com/cpacia/btcd/peer"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

type SpendOutput struct {
	Address   wallet.Address
	Amount    int64
	Type      SpendOutputType
	ReplyHash []byte
	Data      string
}

func Spend(txOut *db.TransactionOut, privateKey *wallet.PrivateKey, spendOutputs []SpendOutput, peerIds []uint) error {
	var peers []*btcdPeer.Peer
	for _, id := range peerIds {
		n, err := peer.Get(id)
		if err != nil {
			return jerr.Getf(err, "error getting peer (id: %d)", id)
		}
		peers = append(peers, n.Peer)
	}

	tx, err := Create(txOut, privateKey, spendOutputs)
	if err != nil {
		return jerr.Get("error getting low fee transaction", err)
	}

	fmt.Println("Sleeping 5 seconds to allow nodes to connect...")
	time.Sleep(5 * time.Second)

	err = Broadcast(tx, peers)
	if err != nil {
		return jerr.Get("error sending low fee transaction", err)
	}
	fmt.Println(GetTxInfo(tx))
	err = SaveTransaction(tx, nil)
	if err != nil {
		return jerr.Get("error saving low fee transaction", err)
	}
	return nil
}
