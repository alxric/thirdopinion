package psql

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"strings"
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
	m := make(map[int64]*config.Argument)
	var rows *sql.Rows
	var err error
	switch filter {
	case "":
		rows, err = db.Query(`SELECT
		t.arg_id, t.arg_title, t.arg_create_time,
		SUM(t.person1) p1votes, SUM(t.person2) p2votes
		FROM
			(SELECT a.arg_id, a.arg_title, a.arg_create_time,
			CASE WHEN v.person = 1 THEN 1 ELSE 0 END person1,
			CASE WHEN v.person = 2 THEN 1 ELSE 0 END person2
			FROM arguments a
			LEFT JOIN votes v ON a.arg_id = v.arg_id
			WHERE v.person IN (1,2) OR v.person IS NULL) t
		GROUP BY
		t.arg_id, t.arg_title,t.arg_create_time`)
	case "specificPost":
		rows, err = db.Query(`SELECT
		t.arg_id, t.arg_title, t.arg_create_time,
		SUM(t.person1) p1votes, SUM(t.person2) p2votes
		FROM
			(SELECT a.arg_id, a.arg_title, a.arg_create_time,
			CASE WHEN v.person = 1 THEN 1 ELSE 0 END person1,
			CASE WHEN v.person = 2 THEN 1 ELSE 0 END person2
			FROM arguments a
			LEFT JOIN votes v ON a.arg_id = v.arg_id
			WHERE (v.person IN (1,2) or v.person IS NULL)
			AND a.arg_id= $1) t
		WHERE t.arg_id = $2
		GROUP BY
		t.arg_id, t.arg_title,t.arg_create_time`,
			filterVals.(int64), filterVals.(int64))
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var title string
	var argID, p1votes, p2votes int
	var argTime time.Time
	var argIDs int64slice
	var arguments []*config.Argument
	for rows.Next() {
		err := rows.Scan(&argID, &title, &argTime, &p1votes, &p2votes)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if _, ok := m[int64(argID)]; !ok {
			argIDs = append(argIDs, int64(argID))
			m[int64(argID)] = &config.Argument{
				Title:        title,
				ID:           int64(argID),
				CreationTime: argTime,
				Opinions:     []*config.Opinion{},
				Votes: config.Votes{
					Person1: int64(p1votes),
					Person2: int64(p2votes),
				},
			}
		}
	}
	err = appendOpinions(db, argIDs, m)
	if err != nil {
		return nil, err
	}
	for _, id := range argIDs {
		arguments = append(arguments, m[int64(id)])
	}

	return arguments, nil
}

func appendOpinions(db *sql.DB, argIDs int64slice, m map[int64]*config.Argument) error {
	argIDsSQL, err := argIDs.Value()
	if err != nil {
		return err
	}
	rows, err := db.Query(`SELECT arg_id, person, text, opinion_id, opinion_timestamp
	FROM opinions o
	WHERE o.arg_id = ANY ($1)`, argIDsSQL)
	if err != nil {
		return err
	}
	defer rows.Close()
	var person, opinionID, argID int
	var text string
	var opinionTime time.Time
	for rows.Next() {
		err := rows.Scan(&argID, &person, &text, &opinionID, &opinionTime)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		o := &config.Opinion{
			Person:       person,
			Text:         text,
			CreationTime: opinionTime,
			ID:           opinionID,
		}
		m[int64(argID)].Opinions = append(m[int64(argID)].Opinions, o)
	}
	return nil
}

type int64slice []int64

func (a int64slice) Value() (driver.Value, error) {
	ints := make([]string, len(a))
	for i, v := range a {
		ints[i] = strconv.FormatInt(v, 10)
	}
	return "{" + strings.Join(ints, ",") + "}", nil
}
