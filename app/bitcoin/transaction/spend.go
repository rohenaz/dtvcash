package transaction

import (
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
)

type SpendOutput struct {
	Address   wallet.Address
	Amount    int64
	Type      SpendOutputType
	ReplyHash []byte
	Data      []byte
}
