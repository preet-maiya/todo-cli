package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type DB struct {
	DBFile string
	conn   *sql.DB
}

// NewDB is a factory method for DB struct
func NewDB(dbFile string) DB {
	conn, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	return DB{DBFile: dbFile, conn: conn}
}
