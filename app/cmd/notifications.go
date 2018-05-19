package cmd

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/notify"
	"github.com/spf13/cobra"
)

var addNotificationsCmd = &cobra.Command{
	Use:   "notifications",
	RunE: func(c *cobra.Command, args []string) error {
		for offset := uint(0); offset < 100; offset += 25 {
			likes, err := db.GetRecentLikes(offset)
			fmt.Printf("Found %d likes\n", len(likes))
			if err != nil {
				jerr.Get("error getting recent likes", err).Print()
				return nil
			}
			for _, like := range likes {
				err := notify.AddLikeNotification(like)
				if err != nil {
					jerr.Get("error adding like notification", err).Print()
					return nil
				}
			}
		}
		fmt.Println("All done")
		return nil
	},
}
