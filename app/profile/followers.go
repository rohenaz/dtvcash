package profile

import (
	"bytes"
	"github.com/rohenaz/dtvcash/app/bitcoin/wallet"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/btcsuite/btcutil"
	"github.com/jchavannes/jgo/jerr"
)

type Follower struct {
	Name       string
	PkHash     []byte
	SelfPkHash []byte
	Reputation *Reputation
}

func (f *Follower) GetAddressString() string {
	address, err := btcutil.NewAddressPubKeyHash(f.PkHash, &wallet.MainNetParamsOld)
	if err != nil {
		return ""
	}
	return address.String()
}

func GetFollowing(selfPkHash []byte, pkHash []byte, offset int) ([]*Follower, error) {
	memoFollows, err := db.GetFollowersForPkHash(pkHash, offset)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting memo follows for hash", err)
	}

	var following []*Follower
MemoFollow:
	for _, memoFollow := range memoFollows {
		var name = "Unknown"
		memoSetName, err := db.GetNameForPkHash(memoFollow.FollowPkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for pk hash", err)
		}
		if memoSetName != nil {
			name = memoSetName.Name
		}
		for i, follower := range following {
			if bytes.Equal(follower.PkHash, memoFollow.FollowPkHash) {
				if memoFollow.Unfollow {
					following = append(following[:i], following[i+1:]...)
				}
				// Check if follower is already set, don't add again
				continue MemoFollow
			}
		}
		if ! memoFollow.Unfollow {
			following = append(following, &Follower{
				Name:       name,
				SelfPkHash: selfPkHash,
				PkHash:     memoFollow.FollowPkHash,
			})
		}
	}
	return following, nil
}

func GetFollowers(selfPkHash []byte, pkHash []byte, offset int) ([]*Follower, error) {
	memoFollows, err := db.GetFollowingForPkHash(pkHash, offset)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting memo follows for hash", err)
	}
	var followers []*Follower
	for _, memoFollow := range memoFollows {
		var name = "Unknown"
		memoSetName, err := db.GetNameForPkHash(memoFollow.PkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for pk hash", err)
		}
		if memoSetName != nil {
			name = memoSetName.Name
		}
		followers = append(followers, &Follower{
			Name:       name,
			PkHash:     memoFollow.PkHash,
			SelfPkHash: selfPkHash,
		})
	}
	return followers, nil
}

func AttachReputationToFollowers(followers []*Follower) error {
	for _, follower := range followers {
		reputation, err := GetReputation(follower.SelfPkHash, follower.PkHash)
		if err != nil {
			return jerr.Get("error getting reputation", err)
		}
		follower.Reputation = reputation
	}
	return nil
}
