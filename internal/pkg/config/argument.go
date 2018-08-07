package config

import "time"

// Argument defines a full argument
type Argument struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	CreationTime time.Time  `json:"creation_time"`
	Opinions     []*Opinion `json:"opinions"`
	Votes        Votes      `json:"votes"`
	UserID       int64      `json:"user_id"`
}

// Opinion defines an individual opinion
type Opinion struct {
	ID           int       `json:"id"`
	Person       int       `json:"person"`
	CreationTime time.Time `json:"creation_time"`
	Text         string    `json:"text"`
}

// Votes contains the votes an argument has gotten
type Votes struct {
	Person1 int64 `json:"person_1"`
	Person2 int64 `json:"person_2"`
}

// CreationResult is the struct used to communicate creation results
// between the websocket and the client
type CreationResult struct {
	Error    string `json:"error"`
	Position string `json:"position"`
}
