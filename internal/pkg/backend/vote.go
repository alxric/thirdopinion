package backend

import (
	"encoding/json"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"

	"github.com/gorilla/websocket"
)

// Vote writes the supplied vote to Postgres
func Vote(ws *websocket.Conn, newVote *config.Vote) ([]byte, error) {
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
