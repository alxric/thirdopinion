package psql

import (
	"database/sql"
	"thirdopinion/internal/pkg/config"

	log "github.com/mgutz/logxi/v1"
)

// Create a new argument entry
func Create(arg *config.Argument) (string, error) {
	db, err := openConn()
	if err != nil {
		return "", err
	}
	defer db.Close()
	lastInsertID, err := createArgument(db, arg.Title)
	if err != nil {
		return "", err
	}
	for _, opinion := range arg.Opinions {
		err := createOpinion(db, opinion, lastInsertID)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return "Argument created", nil
}

func createOpinion(db *sql.DB, opinion *config.Opinion, argumentID int) error {
	stmt, err := db.Prepare("INSERT into opinions (arg_id, person, text) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(argumentID, opinion.Person, opinion.Text)
	if err != nil {
		return err
	}
	return nil
}

func createArgument(db *sql.DB, title string) (int, error) {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO arguments (arg_title) VALUES($1) returning arg_id", title).Scan(&lastInsertID)
	return lastInsertID, err
}
