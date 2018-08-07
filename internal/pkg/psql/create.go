package psql

import (
	"database/sql"
	"thirdopinion/internal/pkg/config"

	log "github.com/mgutz/logxi/v1"
)

// Create a new argument entry
func Create(arg *config.Argument) (*config.WSResponse, error) {
	db, err := openConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	lastInsertID, err := createArgument(db, arg)
	if err != nil {
		return nil, err
	}
	for _, opinion := range arg.Opinions {
		err := createOpinion(db, opinion, lastInsertID)
		if err != nil {
			log.Error(err.Error())
		}
	}
	resp := &config.WSResponse{
		Msg:        "Argument created",
		ArgumentID: lastInsertID,
	}
	return resp, nil
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

func createArgument(db *sql.DB, arg *config.Argument) (int, error) {
	var lastInsertID int
	err := db.QueryRow(`INSERT INTO arguments
	(arg_title, user_id)
	VALUES($1, $2)
	returning arg_id`, arg.Title, arg.UserID).Scan(&lastInsertID)
	return lastInsertID, err
}
