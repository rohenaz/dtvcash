package cmd

import (
	"fmt"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/rohenaz/dtvcash/app/html-parser"
	"github.com/spf13/cobra"
	"log"
)

var fixLeadingCharsCmd = &cobra.Command{
	Use:   "fix-leading-chars",
	RunE: func(c *cobra.Command, args []string) error {
		var diffCount int
		var checked int
		for i := 1; i < 10000; i++ {
			memoPost, err := db.GetMemoPostById(uint(i))
			if err != nil {
				if db.IsRecordNotFoundError(err) {
					continue
				}
				log.Fatal(err)
			}
			checked++
			var newMessage string
			if len(memoPost.ParentTxHash) > 0 {
				newMessage = html_parser.EscapeWithEmojis(string(memoPost.PkScript[38:]))
			} else {
				if len(memoPost.PkScript) > 81 {
					newMessage = html_parser.EscapeWithEmojis(string(memoPost.PkScript[6:]))
				} else {
					newMessage = html_parser.EscapeWithEmojis(string(memoPost.PkScript[5:]))
				}
			}
			if newMessage != memoPost.Message {
				diffCount++
				memoPost.Message = newMessage
				err = memoPost.Save()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		fmt.Printf("Checked %d posts, updated %d\n", checked, diffCount)
		return nil
	},
}
