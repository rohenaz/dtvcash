package transaction

import (
	"fmt"
	"git.jasonc.me/main/memo/app/bitcoin/queuer"
	"git.jasonc.me/main/memo/app/db"
	"github.com/cpacia/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

const waitTime = 500 * time.Millisecond


func QueueAndWaitForTx(tx *wire.MsgTx) error {
	doneChan := make(chan struct{}, 1)
	queuer.Node.Peer.QueueMessage(tx, doneChan)
	<-doneChan
	txHash := tx.TxHash()
	for i := 0; i < 30; i++ {
		_, err := db.GetTransactionByHash(txHash.CloneBytes())
		if err == nil {
			return nil
		}
		if ! db.IsRecordNotFoundError(err) {
			return jerr.Get("error looking for transaction", err)
		}
		time.Sleep(waitTime)
		if i % 5 == 0 {
			fmt.Println("Trying to queue again...")
			queuer.Node.Peer.QueueMessage(tx, doneChan)
			<-doneChan
		}
	}
	return jerr.New("unable to find transaction")
}
