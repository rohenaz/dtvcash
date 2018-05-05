package wallet

import "github.com/btcsuite/btcd/chaincfg/chainhash"

var GenesisBlock Block

func init() {
	hash, _ := chainhash.NewHashFromStr("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	merkleRoot, _ := chainhash.NewHashFromStr("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b")
	GenesisBlock = Block{
		Hash:       hash,
		MerkleRoot: merkleRoot,
	}
}

type Block struct {
	Hash       *chainhash.Hash
	MerkleRoot *chainhash.Hash
}
