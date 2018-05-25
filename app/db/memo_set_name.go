package db

import (
	"bytes"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/script"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"html"
	"sort"
	"time"
)

type MemoSetName struct {
	Id         uint   `gorm:"primary_key"`
	TxHash     []byte `gorm:"unique;size:50"`
	ParentHash []byte
	PkHash     []byte `gorm:"index:pk_hash"`
	PkScript   []byte `gorm:"size:500"`
	Address    string
	Name       string `gorm:"size:500"`
	BlockId    uint
	Block      *Block
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m MemoSetName) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoSetName) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoSetName) GetAddressString() string {
	pkHash, err := btcutil.NewAddressPubKeyHash(m.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		jerr.Get("error getting pubkeyhash from memo post", err).Print()
		return ""
	}
	return pkHash.EncodeAddress()
}

func (m MemoSetName) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoSetName) GetTimeString() string {
	if m.BlockId != 0 {
		return m.Block.Timestamp.Format("2006-01-02 15:04:05")
	}
	return "Unconfirmed"
}

func GetMemoSetNameById(id uint) (*MemoSetName, error) {
	var memoSetName MemoSetName
	err := find(&memoSetName, MemoSetName{
		Id: id,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo set name", err)
	}
	return &memoSetName, nil
}

func GetMemoSetName(txHash []byte) (*MemoSetName, error) {
	var memoSetName MemoSetName
	err := find(&memoSetName, MemoSetName{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo set name", err)
	}
	return &memoSetName, nil
}

type memoSetNameSortByDate []*MemoSetName

func (txns memoSetNameSortByDate) Len() int      { return len(txns) }
func (txns memoSetNameSortByDate) Swap(i, j int) { txns[i], txns[j] = txns[j], txns[i] }
func (txns memoSetNameSortByDate) Less(i, j int) bool {
	if bytes.Equal(txns[i].ParentHash, txns[j].TxHash) {
		return true
	}
	if bytes.Equal(txns[i].TxHash, txns[j].ParentHash) {
		return false
	}
	if txns[i].Block == nil && txns[j].Block == nil {
		return false
	}
	if txns[i].Block == nil {
		return true
	}
	if txns[j].Block == nil {
		return false
	}
	return txns[i].Block.Height > txns[j].Block.Height
}

func GetNameForPkHash(pkHash []byte) (*MemoSetName, error) {
	names, err := GetSetNamesForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting set names for pk hash", err)
	}
	if len(names) == 0 {
		return nil, nil
	}
	return names[0], nil
}

func GetNamesForPkHashes(pkHashes [][]byte) ([]*MemoSetName, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Table("memo_set_names").
		Select("memo_set_names.*, blocks.*").
		Joins("JOIN blocks ON (memo_set_names.block_id = blocks.id)").
		Order("blocks.timestamp DESC").
		Where("pk_hash IN (?)", pkHashes)
	rows, err := query.Rows()
	if err != nil {
		return nil, jerr.Get("error getting set names", err)
	}
	var memoSetNames []*MemoSetName
	for rows.Next() {
		var memoSetName = MemoSetName{
			Block: &Block{},
		}
		err = rows.Scan(
			&memoSetName.Id,
			&memoSetName.TxHash,
			&memoSetName.ParentHash,
			&memoSetName.PkHash,
			&memoSetName.PkScript,
			&memoSetName.Address,
			&memoSetName.Name,
			&memoSetName.BlockId,
			&memoSetName.CreatedAt,
			&memoSetName.UpdatedAt,
			&memoSetName.Block.Id,
			&memoSetName.Block.Height,
			&memoSetName.Block.Timestamp,
			&memoSetName.Block.Hash,
			&memoSetName.Block.PrevBlock,
			&memoSetName.Block.MerkleRoot,
			&memoSetName.Block.Nonce,
			&memoSetName.Block.TxnCount,
			&memoSetName.Block.Version,
			&memoSetName.Block.Bits,
			&memoSetName.Block.CreatedAt,
			&memoSetName.Block.UpdatedAt,
		)
		if err != nil {
			return nil, jerr.Get("error scanning set name", err)
		}
		memoSetNames = append(memoSetNames, &memoSetName)
	}

	var setNames []*MemoSetName
SetNameLoop:
	for _, memoSetName := range memoSetNames {
		for _, setName := range setNames {
			if bytes.Equal(setName.PkHash, memoSetName.PkHash) {
				continue SetNameLoop
			}
		}
		setNames = append(setNames, memoSetName)
	}
	return setNames, nil
}

func GetUniqueMemoAPkHashesMatchName(searchString string, offset int) ([][]byte, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	joinSelect := "JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_set_names" +
		"	GROUP BY pk_hash" +
		") sq ON (sq.id = memo_set_names.id)"
	rows, err := db.
		Table("memo_set_names").
		Select("DISTINCT(pk_hash)").
		Joins(joinSelect).
		Where("name LIKE ?", fmt.Sprintf("%%%s%%", searchString)).
		Limit(25).
		Offset(offset).
		Rows()
	if err != nil {
		return nil, jerr.Get("error getting distinct pk hashes", err)
	}
	defer rows.Close()
	var pkHashes [][]byte
	for rows.Next() {
		var pkHash []byte
		err := rows.Scan(&pkHash)
		if err != nil {
			return nil, jerr.Get("error scanning row with pkHash", err)
		}
		pkHashes = append(pkHashes, pkHash)
	}
	return pkHashes, nil
}

func GetSetNamesForPkHash(pkHash []byte) ([]*MemoSetName, error) {
	var memoSetNames []*MemoSetName
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoSetNames, &MemoSetName{
		PkHash: pkHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo names", err)
	}
	sort.Sort(memoSetNameSortByDate(memoSetNames))
	return memoSetNames, nil
}

func GetCountMemoSetName() (uint, error) {
	cnt, err := count(&MemoSetName{})
	if err != nil {
		return 0, jerr.Get("error getting total count", err)
	}
	return cnt, nil
}
