package psql

import (
	"thirdopinion/internal/pkg/config"

	"golang.org/x/crypto/bcrypt"
)

// Register a new user in the database
func Register(r *config.Register) (string, error) {
	db, err := openConn()
	if err != nil {
		return "", err
	}
	defer db.Close()
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	stmt, err := db.Prepare("INSERT into users (email, pw_hash) VALUES ($1, $2)")
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(r.Email, string(hashedPW))
	if err != nil {
		return "", err
	}
	return "Registered", nil
}
