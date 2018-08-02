package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"thirdopinion/internal/pkg/backend"
	"thirdopinion/internal/pkg/config"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WS defines the websocket handler
func WS(c echo.Context) error {
	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	if c.Request().Host != "localhost:8081" {
		return fmt.Errorf("Unauthorized request")
	}
	// Read
	_, msg, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	wsr := &config.WSRequest{}
	err = json.Unmarshal(msg, wsr)
	if err != nil {
		return err
	}
	var b []byte
	switch wsr.Method {
	case "NewArgument":
		b, err = backend.NewArgument(ws, wsr.Argument)
	case "GetArguments":
		b, err = backend.ListArguments(ws, wsr.Argument)
	case "Vote":
		b, err = backend.Vote(ws, wsr.Vote)
	case "Register":
		b, err = backend.Register(ws, wsr.Register)
	case "Login":
		b, err = backend.Login(ws, wsr.Register)
	}
	if err != nil {
		return err
	}
	err = ws.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return err
	}
	return nil
}
