package psql

import (
	"errors"
	"thirdopinion/internal/pkg/config"

	"golang.org/x/crypto/bcrypt"
)

// Login a user
func Login(r *config.Register) (string, error) {
	db, err := openConn()
	if err != nil {
		return "", err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT u.user_id, u.pw_hash
	FROM users u
	WHERE u.email=$1`, r.Email)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var password string
	var userID int
	if !rows.Next() {
		return "", errors.New("Email does not exist. Change this message Alex")
	}
	err = rows.Scan(&userID, &password)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword(
		[]byte(password), []byte(r.Password),
	); err != nil {
		return "", err
	}
	return "Logged in", nil
}
