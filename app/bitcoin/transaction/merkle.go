package transaction

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cpacia/btcd/wire"
)

func GetTransactionsFromMerkleBlock(merkleBlock *wire.MsgMerkleBlock) []MerkleTransaction {
	txns, err := getTransactionsFromMerkleBlock(merkleBlock)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return []MerkleTransaction{}
	}
	var merkleTransactions []MerkleTransaction
	for _, txn := range txns {
		merkleTransactions = append(merkleTransactions, MerkleTransaction{
			transaction: *txn,
		})
	}
	return merkleTransactions
}

type MerkleTransaction struct {
	transaction chainhash.Hash
}

func (m *MerkleTransaction) GetTxId() chainhash.Hash {
	return m.transaction
}
