package cmd

import (
	"git.jasonc.me/main/memo/web/server"
	"github.com/jchavannes/jgo/jlog"
	"github.com/spf13/cobra"
)

const (
	FlagInsecure  = "insecure"
	FlagDebugMode = "debug"
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
		if debugMode {
			jlog.SetLogLevel(jlog.DEBUG)
		}
		server.Run(sessionCookieInsecure)
		return nil
	},
}

func Execute() {
	memoCmd.AddCommand(webCmd)
	memoCmd.Execute()
}

func init() {
	webCmd.Flags().Bool(FlagInsecure, false, "Allow session cookie over unencrypted HTTP")
	webCmd.Flags().Bool(FlagDebugMode, false, "Debug mode")
}
