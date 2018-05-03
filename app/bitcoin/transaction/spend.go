package transaction

import (
	"git.jasonc.me/main/bitcoin/bitcoin/wallet"
)

type SpendOutput struct {
	Address wallet.Address
	Amount  int64
	Type    SpendOutputType
	RefData []byte
	Data    []byte
}
