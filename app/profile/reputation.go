package profile

import (
	"bytes"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
)

type Reputation struct {
	rep *cache.Reputation
}

func (r Reputation) HasReputation() bool {
	return r.rep != nil
}

func (r Reputation) IsDirectFollow() bool {
	return r.rep.DirectFollow
}

func (r Reputation) GetTrustedFollowers() int {
	return r.rep.TrustedFollowers
}

func (r Reputation) GetTotalFollowing() int {
	return r.rep.TotalFollowing
}

func (r Reputation) GetPercentString() string {
	if r.rep.TotalFollowing == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.0f%%", float32(r.rep.TrustedFollowers)/float32(r.rep.TotalFollowing)*100)
}

func (r Reputation) GetPercentStringIncludingDirect() string {
	return r.GetPercentString()
}

func GetReputation(selfPkHash []byte, pkHash []byte) (*Reputation, error) {
	if len(selfPkHash) == 0 {
		return &Reputation{}, nil
	}
	cachedRep, err := cache.GetReputation(selfPkHash, pkHash)
	if err == nil {
		return &Reputation{
			rep: cachedRep,
		}, nil
	} else if ! cache.IsMissError(err) {
		return nil, jerr.Get("error getting reputation from cache", err)
	}

	trustedUsers, err := db.GetFollowersForPkHash(selfPkHash, -1)
	if err != nil {
		return nil, jerr.Get("error getting trustedUsers", err)
	}
	followersToCheck, err := db.GetFollowingForPkHash(pkHash, -1)
	if err != nil {
		return nil, jerr.Get("error getting followersToCheck", err)
	}
	var directFollow bool
	var deDupedTrustedUsers []*db.MemoFollow
TrustedFollowersDeDupeLoop:
	for _, trustedUser := range trustedUsers {
		if bytes.Equal(trustedUser.FollowPkHash, pkHash) {
			directFollow = true
		}
		for _, deDupedTrustedUser := range deDupedTrustedUsers {
			if bytes.Equal(deDupedTrustedUser.FollowPkHash, trustedUser.FollowPkHash) {
				continue TrustedFollowersDeDupeLoop
			}
		}
		deDupedTrustedUsers = append(deDupedTrustedUsers, trustedUser)
	}
	var trustedFollowers []*db.MemoFollow
TrustedFollowersLoop:
	for _, trustedUser := range deDupedTrustedUsers {
		for _, followerToCheck := range followersToCheck {
			if bytes.Equal(followerToCheck.PkHash, trustedUser.FollowPkHash) {
				trustedFollowers = append(trustedFollowers, followerToCheck)
				continue TrustedFollowersLoop
			}
		}
	}
	var rep = &cache.Reputation{
		TrustedFollowers: len(trustedFollowers),
		TotalFollowing:   len(deDupedTrustedUsers),
		DirectFollow:     directFollow,
	}
	err = cache.SetReputation(selfPkHash, pkHash, rep)
	if err != nil {
		jerr.Get("error saving reputation to cache", err).Print()
	}
	return &Reputation{
		rep: rep,
	}, nil
}
