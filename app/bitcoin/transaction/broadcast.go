package transaction

import (
	"github.com/cpacia/btcd/peer"
	"github.com/cpacia/btcd/wire"
)

func Broadcast(tx *wire.MsgTx, peers []*peer.Peer) error {
	for _, p := range peers {
		p.QueueMessage(tx, nil)
	}
	return nil
}
