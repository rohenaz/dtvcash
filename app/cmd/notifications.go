package cmd

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/notify"
	"github.com/spf13/cobra"
)

var addLikeNotificationsCmd = &cobra.Command{
	Use:   "like-notifications",
	RunE: func(c *cobra.Command, args []string) error {
		for offset := uint(0); offset < 100000; offset += 25 {
			likes, err := db.GetRecentLikes(offset)
			fmt.Printf("Found %d likes\n", len(likes))
			if err != nil {
				jerr.Get("error getting recent likes", err).Print()
				return nil
			}
			for _, like := range likes {
				err := notify.AddLikeNotification(like, false)
				if err != nil {
					jerr.Get("error adding like notification", err).Print()
				}
			}
			if len(likes) != 25 {
				break
			}
		}
		fmt.Println("All done")
		return nil
	},
}

var addReplyNotificationsCmd = &cobra.Command{
	Use:   "reply-notifications",
	RunE: func(c *cobra.Command, args []string) error {
		for offset := uint(0); offset < 100000; offset += 25 {
			posts, err := db.GetRecentReplyPosts(offset)
			fmt.Printf("Found %d posts\n", len(posts))
			if err != nil {
				jerr.Get("error getting recent posts", err).Print()
				return nil
			}
			for _, post := range posts {
				if len(post.ParentTxHash) == 0 {
					continue
				}
				err := notify.AddReplyNotification(post, false)
				if err != nil {
					jerr.Get("error adding reply notification", err).Print()
				}
			}
			if len(posts) != 25 {
				break
			}
		}
		fmt.Println("All done")
		return nil
	},
}
