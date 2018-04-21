package main_node

import (
	"fmt"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
	"github.com/jchavannes/jgo/jerr"
)

func setKeys(n *Node) {
	allKeys, err := db.GetAllKeys()
	if err != nil {
		fmt.Println(jerr.Get("error getting keys from db", err))
		return
	}
	n.Keys = allKeys
}

func saveKeys(n *Node) {
	for _, key := range n.Keys {
		err := key.Save()
		if err != nil {
			fmt.Println(jerr.Get("error updating key", err))
			return
		}
	}
}

func getScriptAddresses(n *Node) []*wallet.Address {
	if len(n.scriptAddresses) == 0 {
		for _, key := range n.Keys {
			address := key.GetAddress()
			n.scriptAddresses = append(n.scriptAddresses, &address)
		}
	}
	return n.scriptAddresses
}

func getKeyFromScriptAddress(n *Node, address *wallet.Address) *db.Key {
	for _, key := range n.Keys {
		if address.GetEncoded() == key.GetAddress().GetEncoded() {
			return key
		}
	}
	return nil
}
