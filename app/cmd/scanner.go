package cmd

import (
	"fmt"
	"github.com/memocash/memo/app/bitcoin/scanner"
	"github.com/spf13/cobra"
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
