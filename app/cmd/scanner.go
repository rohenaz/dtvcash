package cmd

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/scanner"
	"github.com/spf13/cobra"
	"strconv"
)

var scannerCmd = &cobra.Command{
	Use:   "scanner",
	RunE: func(c *cobra.Command, args []string) error {
		scanner.Node.Start()
		scanner.Node.Peer.WaitForDisconnect()
		fmt.Println("Disconnected.")
		return nil
	},
}

var scanRecentCmd = &cobra.Command{
	Use:   "scan-recent",
	Short: "scan-recent [num_blocks_back]",
	RunE: func(c *cobra.Command, args []string) error {
		var numBlocksBack = 100
		if len(args) == 1 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				return jerr.Get("error parsing num blocks back", err)
			}
			numBlocksBack = i
		}
		scanner.Node.NumBlocksBack = uint(numBlocksBack)
		scanner.Node.Start()
		scanner.Node.Peer.WaitForDisconnect()
		fmt.Println("Disconnected.")
		return nil
	},
}
