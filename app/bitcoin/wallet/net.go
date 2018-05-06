package wallet

import (
	chainCfgOld "github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/jchavannes/bchutil"
	"github.com/jchavannes/btcd/chaincfg"
	"github.com/jchavannes/btcd/txscript"
)

var MainNetParams = chaincfg.MainNetParams
var MainNetParamsOld = chainCfgOld.MainNetParams

const SigHashForkID txscript.SigHashType = 0x40

func init() {
	MainNetParams.Net = bchutil.MainnetMagic
	MainNetParamsOld.Net = wire.BitcoinNet(bchutil.MainnetMagic)
}
