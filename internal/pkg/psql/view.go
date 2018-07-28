package psql

import (
	"database/sql"
	"thirdopinion/internal/pkg/config"
	"time"

	log "github.com/mgutz/logxi/v1"
)

// View will list entries based on filter supplied
func View(filter string, filterVals interface{}) ([]*config.Argument, error) {
	db, err := openConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	arguments, err := listArgument(db, filter, filterVals)
	if err != nil {
		return nil, err
	}
	return arguments, nil
}

func listArgument(db *sql.DB, filter string, filterVals interface{}) ([]*config.Argument, error) {
	m := make(map[int]*config.Argument)
	var rows *sql.Rows
	var err error
	switch filter {
	case "":
		rows, err = db.Query(`SELECT
		a.arg_id,a.arg_title,a.arg_create_time,o.person,o.text,
		o.opinion_id,o.opinion_timestamp
		FROM opinions o
		LEFT JOIN arguments a ON o.arg_id=a.arg_id
	`)
	case "specificPost":
		rows, err = db.Query(`SELECT
		a.arg_id,a.arg_title,a.arg_create_time,o.person,o.text,
		o.opinion_id,o.opinion_timestamp
		FROM opinions o
		LEFT JOIN arguments a ON o.arg_id=a.arg_id
		WHERE a.arg_id = $1`, filterVals.(int))
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var title, text string
	var argID, opinionID, person int
	var argTime, opinionTime time.Time
	var argIDs []int
	var arguments []*config.Argument
	for rows.Next() {
		err := rows.Scan(&argID, &title, &argTime, &person, &text,
			&opinionID, &opinionTime)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if _, ok := m[argID]; !ok {
			argIDs = append(argIDs, argID)
			m[argID] = &config.Argument{
				Title:        title,
				ID:           argID,
				CreationTime: argTime,
				Opinions:     []*config.Opinion{},
			}
		}
		o := &config.Opinion{
			Person:       person,
			Text:         text,
			CreationTime: opinionTime,
			ID:           opinionID,
		}
		m[argID].Opinions = append(m[argID].Opinions, o)
	}
	for _, id := range argIDs {
		arguments = append(arguments, m[id])
	}

	return arguments, nil
}
