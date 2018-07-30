package main

import (
	"encoding/json"
	"strconv"
	"thirdopinion/internal/pkg/config"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
)

type messageEvent struct {
	*js.Object
	Data *js.Object `js:"data"`
	Code *js.Object `js:"code"`
}

func fetchPost(argumentID string) (chan *messageEvent, error) {
	url := "ws://localhost:8081/ws"
	ws, err := websocketjs.New(url)
	if err != nil {
		return nil, err
	}
	var iArgumentID int
	switch argumentID {
	case "":
		iArgumentID = 0
	default:
		iArgumentID, err = strconv.Atoi(argumentID)
		if err != nil {
			return nil, err
		}
	}
	wsr := &config.WSRequest{
		Method: "GetArguments",
		Argument: &config.Argument{
			ID: int64(iArgumentID),
		},
	}
	b, err := json.Marshal(wsr)
	if err != nil {
		return nil, err
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

	onError := func(ev *js.Object) {
	}

	onClose := func(ev *js.Object) {
		msg := &messageEvent{Object: ev}
		ch <- msg
	}

	ws.AddEventListener("open", false, onOpen)
	ws.AddEventListener("message", false, onMessage)
	ws.AddEventListener("error", false, onError)
	ws.AddEventListener("close", false, onClose)
	return ch, nil
}

func voteDB(argID int, person int) (chan *messageEvent, error) {
	url := "ws://localhost:8081/ws"
	ws, err := websocketjs.New(url)
	if err != nil {
		return nil, err
	}
	wsr := &config.WSRequest{
		Method: "Vote",
		Vote: &config.Vote{
			ArgumentID: argID,
			Person:     person,
		},
	}

	b, err := json.Marshal(wsr)
	if err != nil {
		return nil, err
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
	return ch, nil
}
