package backend

import (
	"encoding/json"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"

	"github.com/gorilla/websocket"
)

// Login sends the login request to the database
func Login(ws *websocket.Conn, r *config.Register) ([]byte, error) {
	resp := &config.WSResponse{}

	res, err := psql.Login(r)
	switch err {
	case nil:
		resp.Msg = res
	default:
		resp.Error = err.Error()
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}
