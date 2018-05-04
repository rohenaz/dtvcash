package watcher

import (
	"git.jasonc.me/main/memo/app/db"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"time"
)

type Item struct {
	Socket     *web.Socket
	Topic      string
	LastPostId uint
	Error      chan error
}

var items []*Item

func RegisterSocket(socket *web.Socket, topic string, lastPostId uint) error {
	var item = &Item{
		Socket:     socket,
		Topic:      topic,
		LastPostId: lastPostId,
		Error:      make(chan error),
	}
	items = append(items, item)
	return <-item.Error
}

func init() {
	go func() {
		for {
			var topics = make(map[string]uint)
			for _, item := range items {
				_, ok := topics[item.Topic]
				if !ok {
					topics[item.Topic] = item.LastPostId
				}
				if item.LastPostId < topics[item.Topic] {
					topics[item.Topic] = item.LastPostId
				}
			}
			for topic, lastPostId := range topics {
				recentPosts, err := db.GetRecentPostsForTopic(topic, lastPostId)
				if err != nil && !db.IsRecordNotFoundError(err) {
					for i := 0; i < len(items); i++ {
						var item = items[i]
						if item.Topic == topic {
							item.Error <- jerr.Get("error getting recent post for topic", err)
							items = append(items[:i], items[i+1:]...)
							i--
						}
					}
				}
				if len(recentPosts) > 0 {
					for _, recentPost := range recentPosts {
						txHash := recentPost.GetTransactionHashString()
						for i := 0; i < len(items); i++ {
							var item = items[i]
							if item.Topic == topic && item.LastPostId < recentPost.Id {
								item.LastPostId = recentPost.Id
								err = item.Socket.WriteJSON(txHash)
								if err != nil {
									go func(item *Item, err error) {
										item.Error <- jerr.Get("error writing to socket", err)
									}(item, err)
									items = append(items[:i], items[i+1:]...)
									i--
								}
							}
						}
					}
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()
}
