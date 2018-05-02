package profile

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/cache"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/cpacia/bchutil"
	"github.com/jchavannes/jgo/jerr"
	"regexp"
	"strings"
)

type Profile struct {
	Name           string
	PkHash         []byte
	NameTx         []byte
	Profile        string
	ProfileTx      []byte
	Self           bool
	SelfPkHash     []byte
	Balance        int64
	BalanceBCH     float64
	hasBalance     bool
	FollowerCount  uint
	FollowingCount uint
	Followers      []*Follower
	Following      []*Follower
	Reputation     *Reputation
	CanFollow      bool
	CanUnfollow    bool
}

func (p Profile) IsSelf() bool {
	return bytes.Equal(p.PkHash, p.SelfPkHash)
}

func (p Profile) HasBalance() bool {
	return p.hasBalance
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

func (p Profile) GetCashAddressString() string {

	addr, err := btcutil.NewAddressPubKeyHash(p.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		return ""
	}
	cashAddr, err := bchutil.NewCashAddressPubKeyHash(addr.ScriptAddress(), &wallet.MainNetParamsOld)
	if err != nil {
		return ""
	}
	return cashAddr.String()
}

func (p *Profile) SetBalances() error {
	bal, err := cache.GetBalance(p.PkHash)
	if err == nil {
		p.Balance = bal
		p.BalanceBCH = float64(bal) * 1e-8
		p.hasBalance = true
		return nil
	} else if ! cache.IsMissError(err) {
		jerr.Get("error getting balance from cache", err).Print()
	}
	outs, err := db.GetTransactionOutputsForPkHash(p.PkHash)
	if err != nil {
		return jerr.Get("error getting outs", err)
	}
	var balance int64
	var balanceBCH float64

	for _, out := range outs {
		if out.TxnInHashString != "" {
			continue
		}
		balance += out.Value
		balanceBCH += out.ValueInBCH()
	}
	p.Balance = balance
	p.BalanceBCH = balanceBCH
	p.hasBalance = true
	err = cache.SetBalance(p.PkHash, p.Balance)
	if err != nil {
		jerr.Get("error setting balance in cache", err).Print()
	}
	return nil
}

func (p *Profile) SetFollowing() error {
	following, err := GetFollowing(p.PkHash)
	if err != nil {
		return jerr.Get("error getting following for hash", err)
	}
	p.Following = following
	return nil
}

func (p *Profile) SetFollowerCount() error {
	cnt, err := db.GetFollowerCountForPkHash(p.PkHash)
	if err != nil {
		return jerr.Get("error getting follower count for hash", err)
	}
	p.FollowerCount = cnt
	return nil
}

func (p *Profile) SetFollowingCount() error {
	cnt, err := db.GetFollowingCountForPkHash(p.PkHash)
	if err != nil {
		return jerr.Get("error getting following count for hash", err)
	}
	p.FollowingCount = cnt
	return nil
}

func (p *Profile) SetCanFollow() error {
	canFollow, err := CanFollow(p.PkHash, p.SelfPkHash)
	if err != nil {
		return jerr.Get("error getting can follow", err)
	}
	p.CanFollow = canFollow
	p.CanUnfollow = !canFollow
	return nil
}

func (p *Profile) SetReputation() error {
	reputation, err := GetReputation(p.SelfPkHash, p.PkHash)
	if err != nil {
		return jerr.Get("error getting reputation", err)
	}
	p.Reputation = reputation
	return nil
}

func (p Profile) GetText() string {
	if p.Profile == "" {
		return "Not set"
	}
	var re = regexp.MustCompile(`(http[s]?://[^\s]*)`)
	s := re.ReplaceAllString(p.Profile, `<a href="$1" target="_blank">$1</a>`)
	s = strings.Replace(s, "\n", "<br/>", -1)
	return s
}

func GetProfiles(selfPkHash []byte) ([]*Profile, error) {
	pkHashes, err := db.GetUniqueMemoAPkHashes()
	if err != nil {
		return nil, jerr.Get("error getting unique pk hashes", err)
	}
	var profiles []*Profile
	for _, pkHash := range pkHashes {
		profile, err := GetProfile(pkHash, selfPkHash)
		if err != nil {
			return nil, jerr.Get("error getting profile for hash", err)
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func GetProfile(pkHash []byte, selfPkHash []byte) (*Profile, error) {
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
		Name:       name,
		PkHash:     pkHash,
		NameTx:     nameTx,
		SelfPkHash: selfPkHash,
	}
	if profile.Name == "" {
		profile.Name = fmt.Sprintf("Profile %.6s", profile.GetAddressString())
	}
	memoSetProfile, err := db.GetProfileForPkHash(pkHash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting MemoSetProfile for hash", err)
	}
	if memoSetProfile != nil {
		profile.Profile = memoSetProfile.Profile
		profile.ProfileTx = memoSetProfile.TxHash
	}
	return profile, nil
}

func GetProfileAndSetBalances(pkHash []byte, selfPkHash []byte) (*Profile, error) {
	pf, err := GetProfile(pkHash, selfPkHash)
	if err != nil {
		return nil, jerr.Get("error getting profile", err)
	}
	err = pf.SetBalances()
	if err != nil {
		return nil, jerr.Get("error setting balances", err)
	}
	return pf, nil
}

func CanFollow(pkHash []byte, selfPkHash []byte) (bool, error) {
	isFollowing, err := db.IsFollowing(selfPkHash, pkHash)
	if err != nil {
		return false, jerr.Get("error determining is follower from db", err)
	}
	return !isFollowing, nil
}
