package transaction

import (
	"git.jasonc.me/main/memo/app/bitcoin/wallet"
)

type SpendOutput struct {
	Address wallet.Address
	Amount  int64
	Type    SpendOutputType
	RefData []byte
	Data    []byte
}
