package main

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

func writeDB(arg *config.Argument) chan *messageEvent {
	url := "ws://localhost:8081/ws"
	ws, err := websocketjs.New(url)
	if err != nil {
		log.Fatal(err)
	}
	wsr := &config.WSRequest{
		Method:   "NewArgument",
		Argument: arg,
	}
	b, err := json.Marshal(wsr)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan *messageEvent, 1)

	onOpen := func(ev *js.Object) {
		err := ws.Send(b) // Send a binary frame.
		if err != nil {
			return
		}
	}

	onMessage := func(ev *js.Object) {
		msg := &messageEvent{Object: ev}
		ch <- msg
	}
	ws.AddEventListener("open", false, onOpen)
	ws.AddEventListener("message", false, onMessage)
	return ch
}
