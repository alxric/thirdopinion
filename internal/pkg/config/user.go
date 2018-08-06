package config

import "time"

// User defines a visiting user
type User struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	SessionKey string    `json:"session_key"`
	SessionID  int       `json:"session_id"`
	LastSeen   time.Time `json:"last_seen"`
	LoggedIn   bool      `json:"logged_in"`
}
