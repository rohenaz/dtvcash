package profile

import (
	"bytes"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/db"
)

type Poll struct {
	Question   *db.MemoPollQuestion
	Votes      []*db.MemoPollVote
	SelfPkHash []byte
}

type Option struct {
	Name        string
	Votes       int
	UniqueVotes int
	Satoshis    int64
}
func (p *Poll) IsMulti() bool {
	return p.Question.PollType != memo.CodePollTypeSingle
}

func (p *Poll) CanVote() bool {
	if p.IsMulti() {
		return true
	}
	for _, vote := range p.Votes {
		if bytes.Equal(vote.PkHash, p.SelfPkHash) {
			return false
		}
	}
	return true
}

func (p *Poll) GetOptions() []Option {
	var options []Option
	for _, dbOption := range p.Question.Options {
		var voteCount int
		var uniqueVoteCount int
		var satoshis int64
		var previousVotes [][]byte
		for _, vote := range p.Votes {
			if bytes.Equal(vote.OptionTxHash, dbOption.TxHash) {
				var hasVoted bool
				for _, previousVote := range previousVotes {
					if bytes.Equal(vote.PkHash, previousVote) {
						hasVoted = true
					}
				}
				voteCount++
				if ! hasVoted {
					uniqueVoteCount++
					previousVotes = append(previousVotes, vote.PkHash)
				}
				if bytes.Equal(vote.TipPkHash, dbOption.PkHash) {
					satoshis += vote.TipAmount
				}
			}
		}
		var option = Option{
			Name:        dbOption.Option,
			Votes:       voteCount,
			UniqueVotes: uniqueVoteCount,
			Satoshis:    satoshis,
		}
		options = append(options, option)
	}
	return options
}
