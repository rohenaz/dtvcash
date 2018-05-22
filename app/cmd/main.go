package cmd

import (
	"github.com/spf13/cobra"
)

var memoCmd = &cobra.Command{
	Use:   "memo",
	Short: "Run Memo app",
}

func Execute() {
	memoCmd.AddCommand(webCmd)
	memoCmd.AddCommand(mainNodeCmd)
	memoCmd.AddCommand(scannerCmd)
	memoCmd.AddCommand(fixPostEmojisCmd)
	memoCmd.AddCommand(fixNameEmojisCmd)
	memoCmd.AddCommand(viewPostCmd)
	memoCmd.AddCommand(fixLeadingCharsCmd)
	memoCmd.AddCommand(addLikeNotificationsCmd)
	memoCmd.AddCommand(addReplyNotificationsCmd)
	memoCmd.AddCommand(parseTransactionCmd)
	memoCmd.Execute()
}
