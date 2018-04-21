package transaction

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/btcd/wire"
)

func GetTransactionsFromBlock(block *wire.MsgMerkleBlock) []BlockTransaction {
	block.Hashes
	txns, err := getTransactionsFromMerkleBlock(merkleBlock)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return []BlockTransaction{}
	}
	var merkleTransactions []BlockTransaction
	for _, txn := range txns {
		merkleTransactions = append(merkleTransactions, BlockTransaction{
			transaction: *txn,
		})
	}
	return merkleTransactions
}

type BlockTransaction struct {
	transaction chainhash.Hash
}

func (m *BlockTransaction) GetTxId() chainhash.Hash {
	return m.transaction
}
