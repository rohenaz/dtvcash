package cache

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
)

type Reputation struct {
	TrustedFollowers int
	TotalFollowing   int
	DirectFollow     bool
}

func GetReputation(selfPkHash []byte, pkHash []byte) (*Reputation, error) {
	var reputation Reputation
	err := GetItem(getReputationName(selfPkHash, pkHash), &reputation)
	if err != nil {
		return nil, jerr.Get("error getting reputation", err)
	}
	return &reputation, nil
}

func SetReputation(selfPkHash []byte, pkHash []byte, reputation *Reputation) error {
	err := SetItemWithExpiration(getReputationName(selfPkHash, pkHash), reputation, 10 * 60)
	if err != nil {
		return jerr.Get("error setting reputation", err)
	}
	return nil
}

func ClearReputation(selfPkHash []byte, pkHash []byte) error {
	err := DeleteItem(getReputationName(selfPkHash, pkHash))
	if err != nil {
		return jerr.Get("error clearing reputation", err)
	}
	return nil
}

func getReputationName(selfPkHash []byte, pkHash []byte) string {
	return fmt.Sprintf("reputation-%x-%x", selfPkHash, pkHash)
}
