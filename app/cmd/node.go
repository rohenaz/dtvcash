package cmd

import (
	"fmt"
	"github.com/rohenaz/dtvcash/app/bitcoin/main-node"
	"github.com/spf13/cobra"
	"os"
)

var mainNodeCmd = &cobra.Command{
	Use:   "main-node",
	RunE: func(c *cobra.Command, args []string) error {
		main_node.Start()
		main_node.WaitForDisconnect()
		fmt.Println("Disconnected.")
		os.Exit(1)
		return nil
	},
}
