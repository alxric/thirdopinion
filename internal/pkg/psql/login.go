package psql

import (
	"errors"
	"thirdopinion/internal/pkg/config"

	"golang.org/x/crypto/bcrypt"
)

// Login a user
func Login(r *config.Register) (*config.User, error) {
	u := &config.User{
		Email: r.Email,
	}
	db, err := openConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT u.user_id, u.pw_hash
	FROM users u
	WHERE u.email=$1`, r.Email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var password string
	var userID int
	if !rows.Next() {
		return nil, errors.New("Invalid credentials")
	}
	err = rows.Scan(&userID, &password)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword(
		[]byte(password), []byte(r.Password),
	); err != nil {
		return nil, errors.New("Invalid credentials")
	}
	u.ID = userID
	return u, nil
}
