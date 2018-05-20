package cmd

import (
	"fmt"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/html-parser"
	"github.com/spf13/cobra"
	"log"
)

var fixPostEmojisCmd = &cobra.Command{
	Use:   "fix-post-emojis",
	RunE: func(c *cobra.Command, args []string) error {
		var diffCount int
		for i := 1; i < 10000; i++ {
			memoPost, err := db.GetMemoPostById(uint(i))
			if err != nil {
				if db.IsRecordNotFoundError(err) {
					fmt.Printf("all done\n")
					break
				}
				log.Fatal(err)
			}
			var newMessage string
			if len(memoPost.ParentTxHash) > 0 {
				newMessage = html_parser.EscapeWithEmojis(string(memoPost.PkScript[38:]))
			} else {
				newMessage = html_parser.EscapeWithEmojis(string(memoPost.PkScript[5:]))
			}
			if newMessage != memoPost.Message {
				diffCount++
				memoPost.Message = newMessage
				err = memoPost.Save()
				if err != nil {
					log.Fatal(err)
				}
			}
			if i%500 == 0 {
				fmt.Printf("Checked %d posts, updated %d\n", i, diffCount)
			}
		}
		return nil
	},
}

var fixNameEmojisCmd = &cobra.Command{
	Use:   "fix-name-emojis",
	RunE: func(c *cobra.Command, args []string) error {
		var diffCount int
		for i := 1; i < 10000; i++ {
			memoSetName, err := db.GetMemoSetNameById(uint(i))
			if err != nil {
				if db.IsRecordNotFoundError(err) {
					fmt.Printf("all done\n")
					break
				}
				log.Fatal(err)
			}
			newName := html_parser.EscapeWithEmojis(string(memoSetName.PkScript[5:]))
			if newName != memoSetName.Name {
				diffCount++
				memoSetName.Name = newName
				err = memoSetName.Save()
				if err != nil {
					log.Fatal(err)
				}
			}
			if i%500 == 0 {
				fmt.Printf("Checked %d names, updated %d\n", i, diffCount)
			}
		}
		return nil
	},
}
