package profile

import (
	"bytes"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/cache"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/util"
	"regexp"
	"strings"
	"time"
)

type Post struct {
	Name       string
	Memo       *db.MemoPost
	Parent     *Post
	Likes      []*Like
	HasLiked   bool
	SelfPkHash []byte
	ReplyCount uint
	Replies    []*Post
	Reputation *Reputation
	ShowMedia  bool
	Poll       *Poll
}

func (p Post) IsSelf() bool {
	if len(p.SelfPkHash) == 0 {
		return false
	}
	return bytes.Equal(p.SelfPkHash, p.Memo.PkHash)
}

func (p Post) IsLoggedIn() bool {
	return len(p.SelfPkHash) > 0
}

func (p Post) GetTotalTip() int64 {
	var totalTip int64
	for _, like := range p.Likes {
		totalTip += like.Amount
	}
	return totalTip
}

func (p Post) GetMessage() string {
	var msg = p.Memo.Message
	if p.ShowMedia {
		msg = addYoutubeVideos(msg)
		msg = addImgurImages(msg)
		msg = addGiphyImages(msg)
	}
	msg = addLinks(msg)
	return msg
}

func (p Post) IsPoll() bool {
	if !p.Memo.IsPoll || p.Poll == nil {
		return false
	}
	numOptions := len(p.Poll.Question.Options)
	if numOptions >= 2 && int(p.Poll.Question.NumOptions) == numOptions {
		return true
	}
	return false
}

func addYoutubeVideos(msg string) string {
	var re = regexp.MustCompile(`(http[s]?://youtu\.be/)([A-Za-z0-9_\-\?=]+)`)
	msg = re.ReplaceAllString(msg, `<div class="video-container"><iframe frameborder="0" src="https://www.youtube.com/embed/$2"></iframe></div>`)
	re = regexp.MustCompile(`(http[s]?://y2u\.be/)([A-Za-z0-9_\-\?=]+)`)
	msg = re.ReplaceAllString(msg, `<div class="video-container"><iframe frameborder="0" src="https://www.youtube.com/embed/$2"></iframe></div>`)
	re = regexp.MustCompile(`(http[s]?://(www\.)?youtube\.com/watch\?v=)([A-Za-z0-9_\-\?=]+)`)
	msg = re.ReplaceAllString(msg, `<div class="video-container"><iframe frameborder="0" src="https://www.youtube.com/embed/$3"></iframe></div>`)
	return msg
}

func addImgurImages(msg string) string {
	// Album link
	if strings.Contains(msg, "imgur.com/a/") || strings.Contains(msg, "imgur.com/gallery/") {
		return msg
	}
	containsRex := regexp.MustCompile(`\.jpg|\.jpeg|\.png|\.gif|\.gifv`)
	if strings.Contains(msg, ".mp4") {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?imgur\.com/)([^\s]*)`)
		msg = re.ReplaceAllString(msg, `<div class="video-container"><video controls><source src="https://i.imgur.com/$3" type="video/mp4"></video></iframe></div>`)
	} else if !containsRex.MatchString(msg) {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?imgur\.com/)([^\s]*)`)
		msg = re.ReplaceAllString(msg, `<a href="https://i.imgur.com/$3.jpg" target="_blank"><img class="imgur" src="https://i.imgur.com/$3.jpg"/></a>`)
	} else {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?imgur\.com/)([^\s]*)`)
		msg = re.ReplaceAllString(msg, `<a href="https://i.imgur.com/$3.jpg" target="_blank"><img class="imgur" src="https://i.imgur.com/$3"/></a>`)
	}
	return msg
}

func addGiphyImages(msg string) string {
	if strings.Contains(msg, "giphy.com/gifs/") {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?giphy.com/gifs/[a-z-]*-([A-Za-z0-9]+))`)
		msg = re.ReplaceAllString(msg, `<a href="https://i.giphy.com/$3.gif" target="_blank"><img class="imgur" src="https://i.giphy.com/$3.gif"/></a>`)
	} else {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?giphy\.com/)([^\s]*)`)
		msg = re.ReplaceAllString(msg, `<a href="https://i.giphy.com/$3" target="_blank"><img class="imgur" src="https://i.giphy.com/$3"/></a>`)
	}
	return msg
}

func addLinks(msg string) string {
	var re = regexp.MustCompile(`(^|\s)(http[s]?://[^\s]*)`)
	s := re.ReplaceAllString(msg, `$1<a href="$2" target="_blank">$2</a>`)
	return strings.Replace(s, "\n", "<br/>", -1)
}

func (p Post) GetTimeString(timezone string) string {
	if p.Memo.BlockId != 0 {
		if p.Memo.Block != nil {
			timeLayout := "2006-01-02 15:04:05 MST"
			if len(timezone) > 0 {
				displayLocation, err := time.LoadLocation(timezone)
				if err != nil {
					jerr.Get("error finding location", err).Print()
					return p.Memo.Block.Timestamp.Format(timeLayout)
				}
				return p.Memo.Block.Timestamp.In(displayLocation).Format(timeLayout)
			} else {
				return p.Memo.Block.Timestamp.Format(timeLayout)
			}
		} else {
			return "Unknown"
		}
	}
	return "Unconfirmed"
}

func (p Post) GetTimeAgo() string {
	if p.Memo.Block != nil && p.Memo.Block.Timestamp.Before(p.Memo.CreatedAt) {
		ts := p.Memo.Block.Timestamp
		return util.GetTimeAgo(ts)
	} else {
		return util.GetTimeAgo(p.Memo.CreatedAt)
	}
}

func (p Post) GetLastLikeId() uint {
	var lastLikeId uint
	for _, like := range p.Likes {
		if like.Id > lastLikeId {
			lastLikeId = like.Id
		}
	}
	return lastLikeId
}

func GetPostsFeed(selfPkHash []byte, offset uint) ([]*Post, error) {
	dbPosts, err := db.GetPostsFeedForPkHash(selfPkHash, offset)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var foundPkHashes [][]byte
	for _, dbPost := range dbPosts {
		for _, foundPkHash := range foundPkHashes {
			if bytes.Equal(foundPkHash, dbPost.PkHash) {
				continue
			}
		}
		foundPkHashes = append(foundPkHashes, dbPost.PkHash)
	}
	names := make(map[string]string)
	for _, pkHash := range foundPkHashes {
		setName, err := db.GetNameForPkHash(pkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for hash", err)
		}
		if setName == nil {
			continue
		}
		names[string(pkHash)] = setName.Name
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		cnt, err := db.GetPostReplyCount(dbPost.TxHash)
		if err != nil {
			return nil, jerr.Get("error getting post reply count", err)
		}
		post := &Post{
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
			ReplyCount: cnt,
		}
		name, ok := names[string(dbPost.PkHash)]
		if ok {
			post.Name = name
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsForHash(pkHash []byte, selfPkHash []byte, offset uint) ([]*Post, error) {
	var name = ""
	setName, err := db.GetNameForPkHash(pkHash)
	if err != nil {
		return nil, jerr.Get("error getting name for hash", err)
	}
	if setName != nil {
		name = setName.Name
	}
	dbPosts, err := db.GetPostsForPkHash(pkHash, offset)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		cnt, err := db.GetPostReplyCount(dbPost.TxHash)
		if err != nil {
			return nil, jerr.Get("error getting post reply count", err)
		}
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
			ReplyCount: cnt,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostByTxHashWithReplies(txHash []byte, selfPkHash []byte, offset uint) (*Post, error) {
	memoPost, err := db.GetMemoPost(txHash)
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	setName, err := db.GetNameForPkHash(memoPost.PkHash)
	if err != nil {
		return nil, jerr.Get("error getting name for hash", err)
	}
	var name = ""
	if setName != nil {
		name = setName.Name
	}
	cnt, err := db.GetPostReplyCount(txHash)
	if err != nil {
		return nil, jerr.Get("error getting post reply count", err)
	}
	post := &Post{
		Name:       name,
		Memo:       memoPost,
		SelfPkHash: selfPkHash,
		ReplyCount: cnt,
	}
	err = AttachRepliesToPost(post, offset)
	if err != nil {
		return nil, jerr.Get("error attaching replies to post", err)
	}
	return post, nil
}

func GetPostByTxHash(txHash []byte, selfPkHash []byte) (*Post, error) {
	memoPost, err := db.GetMemoPost(txHash)
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	setName, err := db.GetNameForPkHash(memoPost.PkHash)
	if err != nil {
		return nil, jerr.Get("error getting name for hash", err)
	}
	var name = ""
	if setName != nil {
		name = setName.Name
	}
	cnt, err := db.GetPostReplyCount(txHash)
	if err != nil {
		return nil, jerr.Get("error getting post reply count", err)
	}
	post := &Post{
		Name:       name,
		Memo:       memoPost,
		SelfPkHash: selfPkHash,
		ReplyCount: cnt,
	}
	return post, nil
}

func AttachRepliesToPost(post *Post, offset uint) error {
	replyMemoPosts, err := db.GetPostReplies(post.Memo.TxHash, offset)
	if err != nil {
		return jerr.Get("error getting post replies", err)
	}
	var replies []*Post
	for _, reply := range replyMemoPosts {
		setName, err := db.GetNameForPkHash(reply.PkHash)
		if err != nil {
			return jerr.Get("error getting name for reply hash", err)
		}
		var name = ""
		if setName != nil {
			name = setName.Name
		}
		cnt, err := db.GetPostReplyCount(reply.TxHash)
		if err != nil {
			return jerr.Get("error getting post reply count", err)
		}
		replies = append(replies, &Post{
			Name:       name,
			Memo:       reply,
			SelfPkHash: post.SelfPkHash,
			ReplyCount: cnt,
		})
	}
	post.Replies = replies
	return nil
}

func GetRecentPosts(selfPkHash []byte, offset uint) ([]*Post, error) {
	dbPosts, err := db.GetRecentPosts(offset)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		var name string
		setName, err := db.GetNameForPkHash(dbPost.PkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for hash", err)
		}
		if setName != nil {
			name = setName.Name
		}
		cnt, err := db.GetPostReplyCount(dbPost.TxHash)
		if err != nil {
			return nil, jerr.Get("error getting post reply count", err)
		}
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
			ReplyCount: cnt,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetTopPostsNamedRange(selfPkHash []byte, offset uint, timeRange string, personalized bool) ([]*Post, error) {
	var timeStart time.Time
	switch timeRange {
	case TimeRange1Hour:
		timeStart = time.Now().Add(-1 * time.Hour)
	case TimeRange24Hours:
		timeStart = time.Now().Add(-24 * time.Hour)
	case TimeRange7Days:
		timeStart = time.Now().Add(-24 * 7 * time.Hour)
	case TimeRangeAll:
		timeStart = time.Now().Add(-24 * 365 * 10 * time.Hour)
	}
	return GetTopPosts(selfPkHash, offset, timeStart, time.Time{}, personalized)
}

func GetTopPosts(selfPkHash []byte, offset uint, timeStart time.Time, timeEnd time.Time, personalized bool) ([]*Post, error) {
	var dbPosts []*db.MemoPost
	var err error
	if personalized {
		dbPosts, err = db.GetPersonalizedTopPosts(selfPkHash, offset, timeStart, timeEnd)
		if err != nil {
			return nil, jerr.Get("error getting posts for hash", err)
		}
	} else {
		dbPosts, err = db.GetTopPosts(offset, timeStart, timeEnd)
		if err != nil {
			return nil, jerr.Get("error getting posts for hash", err)
		}
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		var name string
		setName, err := db.GetNameForPkHash(dbPost.PkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for hash", err)
		}
		if setName != nil {
			name = setName.Name
		}
		cnt, err := db.GetPostReplyCount(dbPost.TxHash)
		if err != nil {
			return nil, jerr.Get("error getting post reply count", err)
		}
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
			ReplyCount: cnt,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsForTopic(tag string, selfPkHash []byte, offset uint) ([]*Post, error) {
	dbPosts, err := db.GetPostsForTopic(tag, offset)
	if err != nil {
		return nil, jerr.Get("error getting posts for hash", err)
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		var name string
		setName, err := db.GetNameForPkHash(dbPost.PkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for hash", err)
		}
		if setName != nil {
			name = setName.Name
		}
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetOlderPostsForTopic(tag string, selfPkHash []byte, firstPostId uint) ([]*Post, error) {
	dbPosts, err := db.GetOlderPostsForTopic(tag, firstPostId)
	if err != nil {
		return nil, jerr.Get("error getting posts", err)
	}
	var posts []*Post
	for _, dbPost := range dbPosts {
		var name string
		setName, err := db.GetNameForPkHash(dbPost.PkHash)
		if err != nil && ! db.IsRecordNotFoundError(err) {
			return nil, jerr.Get("error getting name for hash", err)
		}
		if setName != nil {
			name = setName.Name
		}
		post := &Post{
			Name:       name,
			Memo:       dbPost,
			SelfPkHash: selfPkHash,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func AttachReplyCountToPosts(posts []*Post) error {
	var txHashes [][]byte
	for _, post := range posts {
		txHashes = append(txHashes, post.Memo.TxHash)
	}
	txHashCounts, err := db.GetPostReplyCounts(txHashes)
	if err != nil {
		return jerr.Get("error getting post reply counts", err)
	}
	for _, txHashCount := range txHashCounts {
		for _, post := range posts {
			if bytes.Equal(post.Memo.TxHash, txHashCount.TxHash) {
				post.ReplyCount = txHashCount.Count
			}
		}
	}
	return nil
}

func AttachParentToPosts(posts []*Post) error {
	for _, post := range posts {
		if len(post.Memo.ParentTxHash) == 0 {
			continue
		}
		parentPost, err := db.GetMemoPost(post.Memo.ParentTxHash)
		if err != nil {
			jerr.Get("error getting memo post parent", err).Print()
			continue
		}
		setName, err := db.GetNameForPkHash(parentPost.PkHash)
		if err != nil {
			return jerr.Get("error getting name for reply hash", err)
		}
		var name = ""
		if setName != nil {
			name = setName.Name
		}
		post.Parent = &Post{
			Name:       name,
			Memo:       parentPost,
			SelfPkHash: post.SelfPkHash,
		}
	}
	return nil
}

func SetShowMediaForPosts(posts []*Post, userId uint) error {
	if userId == 0 {
		for _, post := range posts {
			post.ShowMedia = true
		}
		return nil
	}
	settings, err := cache.GetUserSettings(userId)
	if err != nil {
		return jerr.Get("error getting user settings", err)
	}
	if settings.Integrations == db.SettingIntegrationsAll {
		for _, post := range posts {
			post.ShowMedia = true
		}
	}
	return nil
}

func AttachPollsToPosts(posts []*Post) error {
	for _, post := range posts {
		if post.Memo.IsPoll {
			question, err := db.GetMemoPollQuestion(post.Memo.TxHash)
			if err != nil {
				return jerr.Get("error getting memo poll question", err)
			}
			numOptions := len(question.Options)
			if numOptions < 2 || int(question.NumOptions) != numOptions {
				continue
			}
			post.Poll = &Poll{
				Question:   question,
				SelfPkHash: post.SelfPkHash,
			}
			var optionHashes [][]byte
			for _, option := range question.Options {
				optionHashes = append(optionHashes, option.TxHash)
			}
			single := question.PollType == memo.CodePollTypeSingle
			votes, err := db.GetVotesForOptions(optionHashes, single)
			if err != nil {
				if db.IsRecordNotFoundError(err) {
					continue
				}
				return jerr.Get("error getting votes for options", err)
			}
			post.Poll.Votes = votes
		}
	}
	return nil
}
