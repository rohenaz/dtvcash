package cmd

import (
	"fmt"
	"git.jasonc.me/main/memo/app/bitcoin/main-node"
	"git.jasonc.me/main/memo/app/bitcoin/scanner"
	"git.jasonc.me/main/memo/app/db"
	"git.jasonc.me/main/memo/app/html-parser"
	"git.jasonc.me/main/memo/app/res"
	"git.jasonc.me/main/memo/web/server"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jlog"
	"github.com/spf13/cobra"
	"log"
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

var fixPostEmojisCmd = &cobra.Command{
	Use:   "fix-post-emojis",
	Short: "",
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
	Short: "",
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

var scanPostsCmd = &cobra.Command{
	Use:   "scan-posts",
	Short: "",
	RunE: func(c *cobra.Command, args []string) error {
		var foundCount int
		for i := 1; i < 10000; i++ {
			memoPost, err := db.GetMemoPostById(uint(i))
			if err != nil {
				if db.IsRecordNotFoundError(err) {
					fmt.Printf("all done\n")
					break
				}
				log.Fatal(err)
			}
			if len(memoPost.PkScript) > 80 {
				fmt.Printf("PkScript len: %d, Msglen: %d, Message: %s\n", len(memoPost.PkScript), len(memoPost.Message), memoPost.Message)
				foundCount++
			}
			if i%500 == 0 {
				fmt.Printf("Checked %d posts, found %d\n", i, foundCount)
			}
		}
		return nil
	},
}

var viewPostCmd = &cobra.Command{
	Use:   "view-post",
	Short: "",
	RunE: func(c *cobra.Command, args []string) error {
		hash, err := chainhash.NewHashFromStr("41b531d1821d13c48b2b879c0d44b2e02e858e625d6ba7312497b5cd33b95044")
		if err != nil {
			log.Fatal(err)
		}
		memoPost, err := db.GetMemoPost(hash.CloneBytes())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("PkScript len: %d, Msglen: %d, Message: %s\n", len(memoPost.PkScript), len(memoPost.Message), memoPost.Message)
		return nil
	},
}

var fixLeadingCharsCmd = &cobra.Command{
	Use:   "fix-leading-chars",
	Short: "",
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

func Execute() {
	memoCmd.AddCommand(webCmd)
	memoCmd.AddCommand(mainNodeCmd)
	memoCmd.AddCommand(scannerCmd)
	memoCmd.AddCommand(fixPostEmojisCmd)
	memoCmd.AddCommand(fixNameEmojisCmd)
	memoCmd.AddCommand(scanPostsCmd)
	memoCmd.AddCommand(viewPostCmd)
	memoCmd.AddCommand(fixLeadingCharsCmd)
	memoCmd.Execute()
}

func init() {
	webCmd.Flags().Bool(FlagInsecure, false, "Allow session cookie over unencrypted HTTP")
	webCmd.Flags().Bool(FlagDebugMode, false, "Debug mode")
	webCmd.Flags().Int(FlagAppendNum, 0, "Number appended to js and css files")
}
