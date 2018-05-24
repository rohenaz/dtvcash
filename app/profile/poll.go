package profile

import (
	"bytes"
	"github.com/memocash/memo/app/db"
)

type Poll struct {
	Question *db.MemoPollQuestion
	Votes    []*db.MemoPollVote
}

type Option struct {
	Name     string
	Votes    int
	Satoshis int64
}

func (p *Poll) GetOptions() []Option {
	var options []Option
	for _, dbOption := range p.Question.Options {
		var voteCount int
		var satoshis int64
		for _, vote := range p.Votes {
			if bytes.Equal(vote.OptionTxHash, dbOption.TxHash) {
				voteCount++
				if bytes.Equal(vote.TipPkHash, dbOption.PkHash) {
					satoshis += vote.TipAmount
				}
			}
		}
		var option = Option{
			Name:     dbOption.Option,
			Votes:    voteCount,
			Satoshis: satoshis,
		}
		options = append(options, option)
	}
	return options
}
