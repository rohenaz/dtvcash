package node

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/cpacia/bchutil"
)

var mainNetParams = &chaincfg.MainNetParams

func init() {
	mainNetParams.Net = wire.BitcoinNet(bchutil.MainnetMagic)
}
