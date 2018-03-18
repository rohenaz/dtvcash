package node

import (
	"fmt"
	"git.jasonc.me/main/bitcoin/wallet"
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
)

func SetKeys(n *Node) {
	allKeys, err := db.GetAllKeys()
	if err != nil {
		fmt.Println(jerr.Get("error getting keys from db", err))
		return
	}
	n.Keys = allKeys
}

func GetScriptAddresses(n *Node) []*wallet.Address {
	if len(n.scriptAddresses) == 0 {
		for _, key := range n.Keys {
			address := key.GetAddress()
			n.scriptAddresses = append(n.scriptAddresses, &address)
		}
	}
	return n.scriptAddresses
}
