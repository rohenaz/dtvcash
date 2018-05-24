package profile

import (
	"bytes"
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/db"
)

type Vote struct {
	Name    string
	Option  string
	Message string
	Tip     int64
}

func GetVotesForTxHash(txHash []byte) ([]*Vote, error) {
	question, err := db.GetMemoPollQuestion(txHash)
	if err != nil {
		return nil, jerr.Get("error getting memo poll question", err)
	}
	numOptions := len(question.Options)
	if numOptions < 2 || int(question.NumOptions) != numOptions {
		return nil, jerr.Get("invalid question", err)
	}
	var optionHashes [][]byte
	for _, option := range question.Options {
		optionHashes = append(optionHashes, option.TxHash)
	}
	single := question.PollType == memo.CodePollTypeSingle
	dbVotes, err := db.GetVotesForOptions(optionHashes, single)
	if err != nil {
		if db.IsRecordNotFoundError(err) {
			return []*Vote{}, nil
		}
		return nil, jerr.Get("error getting votes for options", err)
	}
	var namePkHashes [][]byte
	for _, dbVote := range dbVotes {
		namePkHashes = append(namePkHashes, dbVote.PkHash)
	}
	setNames, err := db.GetNamesForPkHashes(namePkHashes)
	if err != nil {
		return nil, jerr.Get("error getting set names for pk hashes", err)
	}
	var votes []*Vote
	for _, dbVote := range dbVotes {
		var name string
		for _, setName := range setNames {
			if bytes.Equal(dbVote.PkHash, setName.PkHash) {
				name = setName.Name
			}
		}
		if name == "" {
			name = fmt.Sprintf("%.10s", dbVote.GetAddressString())
		}
		var optionString string
		for _, option := range question.Options {
			if bytes.Equal(option.TxHash, dbVote.OptionTxHash) {
				optionString = option.Option
			}
		}
		votes = append(votes, &Vote{
			Name:    name,
			Message: dbVote.Message,
			Tip:     dbVote.TipAmount,
			Option:  optionString,
		})
	}
	return votes, nil
}
