package profile

import (
	"bytes"
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
				Name:   name,
				PkHash: memoFollow.FollowPkHash,
			})
		}
	}
	return following, nil
}

func GetFollowers(pkHash []byte) ([]*Follower, error) {
	memoFollows, err := db.GetFollowsForFollowPkHash(pkHash)
	if err != nil && ! db.IsRecordNotFoundError(err) {
		return nil, jerr.Get("error getting memo follows for hash", err)
	}

	var followers []*Follower
MemoFollow:
	for _, memoFollow := range memoFollows {
		var name = "Unknown"
		memoSetName, err := db.GetNameForPkHash(memoFollow.PkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for pk hash", err)
		}
		if memoSetName != nil {
			name = memoSetName.Name
		}

		for i, follower := range followers {
			if bytes.Equal(follower.PkHash, memoFollow.PkHash) {
				if memoFollow.Unfollow {
					followers = append(followers[:i], followers[i+1:]...)
				}
				// Check if follower is already set, don't add again
				continue MemoFollow
			}
		}
		if ! memoFollow.Unfollow {
			followers = append(followers, &Follower{
				Name:   name,
				PkHash: memoFollow.PkHash,
			})
		}
	}
	return followers, nil
}
