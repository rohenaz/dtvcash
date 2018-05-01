package profile

import (
	"bytes"
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
)

type Reputation struct {
	TrustedFollowers int
	TotalFollowing   int
	DirectFollow     bool
}

func (r Reputation) GetPercentString() string {
	return fmt.Sprintf("%.2f", float32(r.TrustedFollowers)/float32(r.TotalFollowing)*100)
}

func GetReputation(selfPkHash []byte, pkHash []byte) (*Reputation, error) {
	trustedUsers, err := db.GetFollowersForPkHash(selfPkHash)
	if err != nil {
		return nil, jerr.Get("error getting trustedUsers", err)
	}
	followersToCheck, err := db.GetFollowingForPkHash(pkHash)
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
	for _, deDupedTrustedUser := range deDupedTrustedUsers {
		name, err := db.GetNameForPkHash(deDupedTrustedUser.FollowPkHash)
		if err != nil {
			return nil, jerr.Get("error getting name for pk hash", err)
		}
		fmt.Printf("deDupedTrustedUser name: %s\n", name.Name)
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
	return &Reputation{
		TrustedFollowers: len(trustedFollowers),
		TotalFollowing:   len(deDupedTrustedUsers),
		DirectFollow:     directFollow,
	}, nil
}
