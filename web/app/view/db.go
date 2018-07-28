package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"thirdopinion/internal/pkg/config"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
)

type messageEvent struct {
	*js.Object
	Data *js.Object `js:"data"`
}

func fetchPost(argumentID string) (chan *messageEvent, error) {
	url := "ws://localhost:8081/ws"
	ws, err := websocketjs.New(url)
	if err != nil {
		log.Fatal(err)
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
			ID: iArgumentID,
		},
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
	return ch, nil
}

func voteDB(buttonID string) (chan *messageEvent, error) {
	url := "ws://localhost:8081/ws"
	ws, err := websocketjs.New(url)
	if err != nil {
		return nil, err
	}
	idVals := strings.Split(buttonID, "_")
	if len(idVals) != 3 {
		return nil, fmt.Errorf("Invalid button ID")
	}
	argID, err := strconv.Atoi(idVals[1])
	if err != nil {
		return nil, err
	}
	person, err := strconv.Atoi(idVals[2])
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
