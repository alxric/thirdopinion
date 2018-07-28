package psql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbUser     = os.Getenv("POSTGRES_USER")
	dbPassword = os.Getenv("POSTGRES_PASSWORD")
	dbName     = os.Getenv("POSTGRES_DATABASE")
)

func openConn() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
