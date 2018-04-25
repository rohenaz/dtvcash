package cmd

import (
	"fmt"
	"git.jasonc.me/main/memo/app/bitcoin/main-node"
	"git.jasonc.me/main/memo/app/bitcoin/scanner"
	"git.jasonc.me/main/memo/app/res"
	"git.jasonc.me/main/memo/web/server"
	"github.com/jchavannes/jgo/jlog"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
)

const (
	FlagInsecure  = "insecure"
	FlagDebugMode = "debug"
	FlagAppendNum = "append-num"
)

var memoCmd = &cobra.Command{
	Use:   "memo",
	Short: "Run Memo app",
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run Memo web",
	RunE: func(c *cobra.Command, args []string) error {
		sessionCookieInsecure, _ := c.Flags().GetBool(FlagInsecure)
		debugMode, _ := c.Flags().GetBool(FlagDebugMode)
		appendNum, _ := c.Flags().GetInt(FlagAppendNum)
		if appendNum == 0 {
			appendNum = rand.Intn(1e5)
		}
		res.SetAppendNumber(appendNum)
		if debugMode {
			jlog.SetLogLevel(jlog.DEBUG)
		}
		server.Run(sessionCookieInsecure)
		return nil
	},
}

var mainNodeCmd = &cobra.Command{
	Use:   "main-node",
	Short: "",
	RunE: func(c *cobra.Command, args []string) error {
		main_node.Start()
		main_node.WaitForDisconnect()
		fmt.Println("Disconnected.")
		os.Exit(1)
		return nil
	},
}

var scannerCmd = &cobra.Command{
	Use:   "scanner",
	Short: "",
	RunE: func(c *cobra.Command, args []string) error {
		scanner.Node.Start()
		scanner.Node.Peer.WaitForDisconnect()
		fmt.Println("Disconnected.")
		return nil
	},
}

func Execute() {
	memoCmd.AddCommand(webCmd)
	memoCmd.AddCommand(mainNodeCmd)
	memoCmd.AddCommand(scannerCmd)
	memoCmd.Execute()
}

func init() {
	webCmd.Flags().Bool(FlagInsecure, false, "Allow session cookie over unencrypted HTTP")
	webCmd.Flags().Bool(FlagDebugMode, false, "Debug mode")
	webCmd.Flags().Int(FlagAppendNum, 0, "Number appended to js and css files")
}
