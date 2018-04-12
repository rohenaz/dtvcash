package cmd

import (
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/web/server"
	"github.com/jchavannes/jgo/jerr"
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

var randomTaskCmd = &cobra.Command{
	Use:   "random-task",
	Short: "",
	RunE: func(c *cobra.Command, args []string) error {
		keys, err := db.GetAllKeys()
		if err != nil {
			return jerr.Get("error getting keys from db", err)
		}
		for _, key := range keys {
			key.PkHash = key.GetAddress().GetScriptAddress()
			err = key.Save()
			if err != nil {
				return jerr.Get("error saving key", err)
			}
		}
		return nil
	},
}

func Execute() {
	memoCmd.AddCommand(webCmd)
	memoCmd.AddCommand(randomTaskCmd)
	memoCmd.Execute()
}

func init() {
	webCmd.Flags().Bool(FlagInsecure, false, "Allow session cookie over unencrypted HTTP")
	webCmd.Flags().Bool(FlagDebugMode, false, "Debug mode")
}
