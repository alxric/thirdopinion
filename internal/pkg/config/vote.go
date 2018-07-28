package config

// Vote defines the vote we add to the database
type Vote struct {
	ArgumentID int `json:"argument_id"`
	Person     int `json:"person"`
	User       int `json:"user"`
}
