package transaction

import (
	"github.com/jchavannes/btcd/peer"
	"github.com/jchavannes/btcd/wire"
)

func Broadcast(tx *wire.MsgTx, peers []*peer.Peer) error {
	for _, p := range peers {
		p.QueueMessage(tx, nil)
	}
	return nil
}
