package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"

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
	switch wsr.Method {
	case "NewArgument":
		res, err := newArgument(ws, wsr.Argument, msg)
		if err != nil {
			return err
		}
		err = ws.WriteMessage(websocket.TextMessage, []byte(res))
		if err != nil {
			return err
		}
	case "GetArguments":
		b, err := getArguments(ws, wsr.Argument)
		if err != nil {
			return err
		}
		err = ws.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			return err
		}
	case "Vote":
		b, err := vote(ws, wsr.Vote)
		if err != nil {
			return err
		}
		err = ws.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func getArguments(ws *websocket.Conn, argument *config.Argument) ([]byte, error) {
	var filter string
	if argument.ID != 0 {
		filter = "specificPost"
	}
	m, err := psql.View(filter, argument.ID)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func newArgument(ws *websocket.Conn, argument *config.Argument, msg []byte) (string, error) {
	cr := verifyArgument(argument)
	if cr.Error != "" {
		b, err := json.Marshal(cr)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Unknown error"))
			return "", err
		}
		ws.WriteMessage(websocket.TextMessage, b)
		return "", err
	}
	res, err := psql.Create(msg)
	if err != nil {
		return "", err
	}
	return res, nil
}

func verifyArgument(arg *config.Argument) config.CreationResult {
	if len(arg.Title) < 5 {
		return config.CreationResult{
			Error:    "Title too short",
			Position: "title",
		}
	}
	if len(arg.Opinions) <= 1 {
		return config.CreationResult{
			Error:    "You need at least two opinions for an argument",
			Position: "argument",
		}
	}
	return config.CreationResult{
		Error: "",
	}

}

func vote(ws *websocket.Conn, newVote *config.Vote) ([]byte, error) {
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
