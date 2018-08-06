package backend

import (
	"encoding/json"
	"errors"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"
)

// Vote writes the supplied vote to Postgres
func Vote(newVote *config.Vote) ([]byte, error) {
	if newVote.User == 0 {
		return nil, errors.New("Could not detect user ID. Try logging in again")
	}
	res, err := psql.Vote(newVote)
	if err != nil {
		return nil, err
	}
	resp := &config.WSResponse{
		Msg: res,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UserVotes list all votes for the supplied user key
func UserVotes(sessionKey string) (map[int]int, error) {
	if sessionKey == "" {
		return nil, errors.New("Invalid session key")
	}
	votes, err := psql.UserVotes(sessionKey)
	if err != nil {
		return nil, err
	}
	return votes, nil
}
