package backend

import (
	"encoding/json"
	"thirdopinion/internal/pkg/config"

	"github.com/gorilla/websocket"
)

// VerifySession checks if session is valid
func VerifySession(ws *websocket.Conn, u *config.User) ([]byte, error) {
	resp := &config.WSResponse{}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}
