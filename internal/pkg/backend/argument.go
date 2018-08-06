package backend

import (
	"encoding/json"
	"thirdopinion/internal/pkg/config"
	"thirdopinion/internal/pkg/psql"

	"github.com/gorilla/websocket"
)

// ListArguments lists arguments from the database
func ListArguments(argument *config.Argument) ([]byte, error) {
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

// NewArgument writes an argument to the database
func NewArgument(ws *websocket.Conn, argument *config.Argument) ([]byte, error) {
	cr := verifyArgument(argument)
	if cr.Error != "" {
		b, err := json.Marshal(cr)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Unknown error"))
			return nil, err
		}
		ws.WriteMessage(websocket.TextMessage, b)
		return nil, err
	}
	res, err := psql.Create(argument)
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
