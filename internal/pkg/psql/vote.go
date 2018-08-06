package psql

import (
	"thirdopinion/internal/pkg/config"
)

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
	_, err = stmt.Exec(vote.ArgumentID, vote.Person, vote.User)
	if err != nil {
		return "", err
	}
	return "Voted", nil
}

// UserVotes fetches votes for a specific user
func UserVotes(sessionKey string) (map[int]int, error) {
	db, err := openConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT arg_id, person FROM votes WHERE user_id=(
		SELECT user_id FROM user_sessions WHERE session_key=$1)`, sessionKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	votes := make(map[int]int)
	var argumentID, person int
	for rows.Next() {
		err := rows.Scan(&argumentID, &person)
		if err != nil {
			return nil, err
		}
		votes[argumentID] = person
	}
	return votes, nil

}
