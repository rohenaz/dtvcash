package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/jchavannes/btcd/txscript"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/spf13/cobra"
)

var parseTransactionCmd = &cobra.Command{
	Use:   "parse-transaction",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) < 1 {
			return jerr.New("not enough arguments")
		}
		var txHex = args[0]
		rawHex, err := hex.DecodeString(txHex)
		if err != nil {
			jerr.Get("error decoding hex", err).Print()
			return nil
		}
		msg := wire.NewMsgTx(1)

		pr := bytes.NewBuffer(rawHex)
		err = msg.BtcDecode(pr, 0)

		str, err := txscript.DisasmString(msg.TxOut[1].PkScript)
		if err != nil {
			jerr.Get("error disassembling data", err).Print()
			return nil
		}
		fmt.Printf("disassembled output: %s\n", str)
		pushedData, err := txscript.PushedData(msg.TxOut[1].PkScript)
		if err != nil {
			jerr.Get("error getting push data", err).Print()
			return nil
		}
		fmt.Printf("pushedData: %#v\n", pushedData)
		return nil
	},
}
