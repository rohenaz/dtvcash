package db

import (
	"git.jasonc.me/main/bitcoin/wallet"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"time"
)

const BlockTable = "Block"

type Block struct {
	Id         uint   `gorm:"primary_key"`
	Height     uint   `gorm:"unique;"`
	Timestamp  time.Time
	Hash       []byte `gorm:"unique;"`
	PrevBlock  []byte
	MerkleRoot []byte
	Nonce      uint32
	TxnCount   uint32
	Version    int32
	Bits       uint32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (b *Block) Save() error {
	result := save(b)
	if result.Error != nil {
		return jerr.Get("error saving block", result.Error)
	}
	return nil
}

func (b Block) GetChainhash() *chainhash.Hash {
	hash, _ := chainhash.NewHash(b.Hash)
	return hash
}

func (b Block) GetPrevBlockChainhash() *chainhash.Hash {
	hash, _ := chainhash.NewHash(b.PrevBlock)
	return hash
}

func (b Block) GetMerkleRoot() *chainhash.Hash {
	hash, _ := chainhash.NewHash(b.MerkleRoot)
	return hash
}

func SaveBlocks(blocks []*Block) error {
	for _, block := range blocks {
		err := block.Save()
		if err != nil {
			return jerr.Get("error saving blocks", err)
		}
	}
	return nil
}

func GetGenesis() (*Block, error) {
	var block = Block{
		Height:     0,
		Timestamp:  time.Unix(1231006505, 0),
		Hash:       wallet.GenesisBlock.Hash.CloneBytes(),
		MerkleRoot: wallet.GenesisBlock.MerkleRoot.CloneBytes(),
		Nonce:      2083236893,
		TxnCount:   1,
		Version:    1,
		Bits:       0x1d00ffff,
	}
	err := find(&block, &block)
	if err == nil {
		return &block, nil
	}
	if ! IsRecordNotFoundError(err) {
		return nil, jerr.Get("error finding genesis block", err)
	}
	err = create(&block)
	if err != nil {
		return nil, jerr.Get("error creating block", err)
	}
	return &block, nil
}

func ConvertMessageToBlock(msg *wire.MsgMerkleBlock) (*Block) {
	block := ConvertMessageHeaderToBlock(&msg.Header)
	block.TxnCount = msg.Transactions
	return block
}

func ConvertMessageHeaderToBlock(header *wire.BlockHeader) (*Block) {
	blockHash := header.BlockHash()
	return &Block{
		Timestamp:  header.Timestamp,
		Hash:       blockHash.CloneBytes(),
		PrevBlock:  header.PrevBlock.CloneBytes(),
		MerkleRoot: header.MerkleRoot.CloneBytes(),
		Nonce:      header.Nonce,
		Version:    header.Version,
		Bits:       header.Bits,
	}
}

func GetBlockByHeight(height uint) (*Block, error) {
	var block = Block{
		Height: height,
	}
	err := find(&block, &block)
	if err != nil {
		return nil, jerr.Get("error getting block", err)
	}
	return &block, nil
}

func GetBlockByHash(hash chainhash.Hash) (*Block, error) {
	var block = Block{
		Hash: hash.CloneBytes(),
	}
	err := find(&block, &block)
	if err != nil {
		return nil, jerr.Get("error getting block", err)
	}
	return &block, nil
}

func AddBlock(block *Block) error {
	parent, err := GetBlockByHash(*block.GetPrevBlockChainhash())
	if err != nil {
		return jerr.Get("error getting parent", err)
	}
	block.Height = parent.Height + 1
	_, err = GetBlockByHash(*block.GetChainhash())
	if err == nil {
		return jerr.New("block already exists")
	}
	result := save(&block)
	if result.Error != nil {
		return jerr.Get("error saving block", result.Error)
	}
	/*transactions := transaction.GetTransactionsFromMerkleBlock(msg)
	if len(transactions) > 0 {
		for _, txn := range transactions {
			txn.GetTxId()
		}
		fmt.Printf("Txns: %d (height: %d)\n", len(transactions), newBlock.Height)
	}*/
	return nil
}

func GetRecentBlock() (*Block, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var block Block
	result := db.Order("height desc").First(&block)
	if result.Error != nil {
		if IsRecordNotFoundError(result.Error) {
			return GetGenesis()
		}
		return nil, jerr.Get("error querying first block", err)
	}
	return &block, nil
}

func GetBlocksInHeightRange(startHeight uint, endHeight uint) ([]*Block, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var blocks []*Block
	query := db
	if startHeight > endHeight {
		query = query.Order("height desc")
		startHeight, endHeight = endHeight, startHeight
	}
	query = query.Where("height >= ? AND height <= ?", startHeight, endHeight)
	result := query.Find(&blocks)
	if result.Error != nil {
		return nil, jerr.Get("error querying blocks", result.Error)
	}
	return blocks, nil
}
