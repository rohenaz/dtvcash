package cmd

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/rohenaz/dtvcash/app/db"
	"github.com/spf13/cobra"
	"log"
)

var viewPostCmd = &cobra.Command{
	Use:   "view-post",
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
