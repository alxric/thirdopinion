package psql

import "thirdopinion/internal/pkg/config"

// Vote adds a new vote to the database
func Vote(vote *config.Vote) (string, error) {
	db, err := openConn()
	if err != nil {
		return "", err
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT into votes (arg_id, person, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(vote.ArgumentID, vote.Person, 0)
	if err != nil {
		return "", err
	}
	return "Voted", nil
}
