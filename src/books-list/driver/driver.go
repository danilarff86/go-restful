package driver

import (
	"log"

	"database/sql"
	"github.com/lib/pq"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB() *sql.DB {
	pgURL, err := pq.ParseURL("postgres://postgres:docker@localhost/postgres?sslmode=disable")
	logFatal(err)

	db, err := sql.Open("postgres", pgURL)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	return db
}
