package cmd

import (
	"github.com/jchavannes/jgo/jlog"
	"github.com/memocash/memo/app/res"
	"github.com/memocash/memo/web/server"
	"github.com/spf13/cobra"
	"math/rand"
)

const (
	FlagInsecure  = "insecure"
	FlagDebugMode = "debug"
	FlagAppendNum = "append-num"
)

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

func init() {
	webCmd.Flags().Bool(FlagInsecure, false, "Allow session cookie over unencrypted HTTP")
	webCmd.Flags().Bool(FlagDebugMode, false, "Debug mode")
	webCmd.Flags().Int(FlagAppendNum, 0, "Number appended to js and css files")
}
