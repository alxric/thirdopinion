package config

import "time"

// Argument defines a full argument
type Argument struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	CreationTime time.Time  `json:"creation_time"`
	Opinions     []*Opinion `json:"opinions"`
}

// Opinion defines an individual opinion
type Opinion struct {
	ID           int       `json:"id"`
	Person       int       `json:"person"`
	CreationTime time.Time `json:"creation_time"`
	Text         string    `json:"text"`
}

// CreationResult is the struct used to communicate creation results
// between the websocket and the client
type CreationResult struct {
	Error    string `json:"error"`
	Position string `json:"position"`
}
