package transaction

import (
	"github.com/memocash/memo/app/bitcoin/queuer"
	"github.com/memocash/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

const waitTime = 200 * time.Millisecond

func QueueAndWaitForTx(tx *wire.MsgTx) error {
	QueueTx(tx)
	txHash := tx.TxHash()
	return WaitForTx(&txHash)
}

func QueueTx(tx *wire.MsgTx) {
	doneChan := make(chan struct{}, 1)
	queuer.Node.Peer.QueueMessage(tx, doneChan)
	<-doneChan
}

func WaitForTx(txHash *chainhash.Hash) error {
	// wait up to 30 seconds
	for i := 0; i < 150; i++ {
		_, err := db.GetTransactionByHash(txHash.CloneBytes())
		if err == nil {
			return nil
		}
		if ! db.IsRecordNotFoundError(err) {
			return jerr.Get("error looking for transaction", err)
		}
		time.Sleep(waitTime)
	}
	return jerr.New("unable to find transaction")
}
