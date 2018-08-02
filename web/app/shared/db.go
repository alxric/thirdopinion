package shared

import (
	"encoding/json"
	"log"
	"thirdopinion/internal/pkg/config"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
)

type messageEvent struct {
	*js.Object
	Data *js.Object `js:"data"`
}

// WriteDB sends the request to the websocket and returns a channel with
// the response
func WriteDB(wsr *config.WSRequest) chan *config.WSResponse {
	ch := make(chan *config.WSResponse)
	wsresp := &config.WSResponse{}

	url := "ws://localhost:8081/ws"
	ws, err := websocketjs.New(url)
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(wsr)
	if err != nil {
		log.Fatal(err)
	}

	onOpen := func(ev *js.Object) {
		err := ws.Send(b) // Send a binary frame.
		if err != nil {
			return
		}
	}

	onMessage := func(ev *js.Object) {
		msg := &messageEvent{Object: ev}
		err := json.Unmarshal([]byte(msg.Data.String()), wsresp)
		if err != nil {
			wsresp.Error = err.Error()
		}
		ch <- wsresp
	}
	ws.AddEventListener("open", false, onOpen)
	ws.AddEventListener("message", false, onMessage)

	return ch
}
