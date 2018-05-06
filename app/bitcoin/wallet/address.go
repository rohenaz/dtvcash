package wallet

import (
	"github.com/btcsuite/btcutil"
)

func GetAddress(pubKey []byte) Address {
	if len(pubKey) == 0 {
		return Address{}
	}
	addr, err := btcutil.NewAddressPubKey(pubKey, &MainNetParamsOld)
	if err != nil {
		//fmt.Println(jerr.Get("error getting address", err))
		return Address{}
	}
	address, err := btcutil.DecodeAddress(addr.EncodeAddress(), &MainNetParamsOld)
	if err != nil {
		//fmt.Printf("error decoding address: %v\n", err)
		return Address{}
	}
	return Address{
		address: address,
	}
}

func GetAddressFromString(addressString string) Address {
	address, err := btcutil.DecodeAddress(addressString, &MainNetParamsOld)
	if err != nil {
		//fmt.Printf("error decoding address: %v\n", err)
	}
	return Address{
		address: address,
	}
}

func GetAddressFromPkHash(pkHash []byte) Address {
	addr, err := btcutil.NewAddressPubKeyHash(pkHash, &MainNetParamsOld)
	if err != nil {
		//fmt.Println(jerr.Get("error getting address", err))
		return Address{}
	}
	address, err := btcutil.DecodeAddress(addr.EncodeAddress(), &MainNetParamsOld)
	if err != nil {
		//fmt.Printf("error decoding address: %v\n", err)
		return Address{}
	}
	return Address{
		address: address,
	}
}

type Address struct {
	address btcutil.Address
}

func (a Address) GetEncoded() string {
	if a.address == nil {
		return ""
	}
	return a.address.EncodeAddress()
}

func (a Address) GetAddress() btcutil.Address {
	return a.address
}

func (a Address) GetScriptAddress() []byte {
	if a.address == nil {
		return []byte{}
	}
	return a.address.ScriptAddress()
}
