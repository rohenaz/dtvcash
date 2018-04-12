package profile

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
)

type Profile struct {
	Name       string
	PkHash     []byte
	NameTx     []byte
	Self       bool
	Balance    int64
	BalanceBCH float64
}

func (p Profile) HasBalance() bool {
	return p.Balance != 0
}

func (p Profile) NameSet() bool {
	return len(p.NameTx) > 0
}

func (p Profile) GetNameTx() string {
	hash, err := chainhash.NewHash(p.NameTx)
	if err != nil {
		return ""
	}
	return hash.String()
}

func (p Profile) GetAddressString() string {
	addr, err := btcutil.NewAddressPubKeyHash(p.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		return ""
	}
	return addr.String()
}

func (p *Profile) SetBalances() error {
	transactions, err := db.GetTransactionsForPkHash(p.PkHash)
	if err != nil {
		return jerr.Get("error getting transactions for key", err)
	}
	var balance int64
	var balanceBCH float64
	for _, transaction := range transactions {
		for address, value := range transaction.GetValues() {
			if address == p.GetAddressString() {
				balance += value.GetValue()
				balanceBCH += value.GetValueBCH()
			}
		}
	}
	p.Balance = balance
	p.BalanceBCH = balanceBCH
	return nil
}

func GetProfiles() ([]*Profile, error) {
	pkHashes, err := db.GetUniqueMemoAPkHashes()
	if err != nil {
		return nil, jerr.Get("error getting unique pk hashes", err)
	}
	var profiles []*Profile
	for _, pkHash := range pkHashes {
		profile, err := GetProfile(pkHash)
		if err != nil {
			return nil, jerr.Get("error getting profile for hash", err)
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func GetProfile(pkHash []byte) (*Profile, error) {
	var name string
	var nameTx []byte
	memoSetName, err := db.GetNameForPkHash(pkHash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting MemoSetName for hash", err)
	}
	if memoSetName != nil {
		name = memoSetName.Name
		nameTx = memoSetName.TxHash
	}
	profile := &Profile{
		Name:   name,
		PkHash: pkHash,
		NameTx: nameTx,
	}
	if profile.Name == "" {
		profile.Name = fmt.Sprintf("Profile %.6s", profile.GetAddressString())
	}
	return profile, nil
}

func GetProfileAndSetBalances(pkHash []byte) (*Profile, error) {
	pf, err := GetProfile(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting profile", err)
	}
	err = pf.SetBalances()
	if err != nil {
		return nil, jerr.Get("error setting balances", err)
	}
	return pf, nil
}
