package transaction

import (
	"github.com/rohenaz/dtvcash/app/bitcoin/wallet"
)

type SpendOutput struct {
	Address wallet.Address
	Amount  int64
	Type    SpendOutputType
	RefData []byte
	Data    []byte
}
