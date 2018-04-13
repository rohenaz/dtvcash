package profile

import (
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
)

type Follower struct {
	Name   string
	PkHash []byte
}

func (f *Follower) GetAddressString() string {
	address, err := btcutil.NewAddressPubKeyHash(f.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		return ""
	}
	return address.String()
}

func GetFollowing(pkHash []byte) ([]*Follower, error) {
	memoFollows, err := db.GetFollowsForPkHash(pkHash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting memo follows for hash", err)
	}

	var following []*Follower
	for _, memoFollow := range memoFollows {
		var name = "Unknown"
		memoSetName, err := db.GetNameForPkHash(memoFollow.FollowPkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for pk hash", err)
		}
		if memoSetName != nil {
			name = memoSetName.Name
		}
		following = append(following, &Follower{
			Name: name,
			PkHash: memoFollow.FollowPkHash,
		})
	}
	return following, nil
}
