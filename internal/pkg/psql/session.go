package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"thirdopinion/internal/pkg/config"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// How old can a session be before we invalidate it?
const maxSession = time.Duration(1 * time.Hour)

// UpdateSession updates the database session
func UpdateSession(u *config.User) error {
	db, err := openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	err = sessionExists(db, u)
	if err != nil {
		return err
	}
	err = generateSessionKey(u)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare(`UPDATE user_sessions
	SET session_key=$1, last_seen=(now() at time zone 'utc')
	WHERE session_id=$2`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(u.SessionKey, u.SessionID)
	if err != nil {
		return err
	}
	return nil
}

// ValidateSession to see if it is okay
func ValidateSession(u *config.User) error {
	db, err := openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT us.user_id, us.session_id, us.last_seen, u.email
	FROM user_sessions us
	JOIN users u ON u.user_id=us.user_id
	WHERE us.session_key=$1`, u.SessionKey)
	if err != nil {
		return err
	}
	defer rows.Close()
	switch rows.Next() {
	case true:
		err = rows.Scan(&u.ID, &u.SessionID, &u.LastSeen, &u.Email)
		if err != nil {
			return err
		}
		if time.Now().Sub(u.LastSeen) > maxSession {
			go deleteSession(u)
			return errors.New("Session too old")
		}
		go lastSeen(u)
	case false:
		return errors.New("Invalid session")
	}
	return nil
}

func sessionExists(db *sql.DB, u *config.User) error {
	rows, err := db.Query(`SELECT u.session_id
	FROM user_sessions u
	WHERE u.user_id=$1`, u.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var sessionID int
	switch rows.Next() {
	case true:
		err = rows.Scan(&sessionID)
		if err != nil {
			return err
		}
	case false:
		err = db.QueryRow(`INSERT INTO user_sessions
		(user_id)
		VALUES($1)
		returning session_id`, u.ID).Scan(&sessionID)
		if err != nil {
			return err
		}
	}
	u.SessionID = sessionID
	return nil
}

func generateSessionKey(u *config.User) error {
	var sessionGen strings.Builder
	sessionGen.WriteString(u.Email)
	sessionGen.WriteString(strconv.Itoa(u.ID))
	sessionGen.WriteString(strconv.Itoa(u.SessionID))
	skey, err := bcrypt.GenerateFromPassword(
		[]byte(sessionGen.String()), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.SessionKey = string(skey)
	return nil
}

func lastSeen(u *config.User) error {
	db, err := openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(`UPDATE user_sessions
	SET last_seen=(now() at time zone 'utc')
	WHERE session_id=$1`)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec(u.SessionID)
	if err != nil {
		return err
	}
	return nil
}

func deleteSession(u *config.User) error {
	db, err := openConn()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(`DELETE FROM user_sessions
	WHERE session_id=$1`)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = stmt.Exec(u.SessionID)
	if err != nil {
		return err
	}
	return nil
}
