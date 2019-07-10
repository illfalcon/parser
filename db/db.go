package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func dbConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Prepare() {
	db := dbConn()
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS landings (
      url TEXT NOT NULL,
      hash TEXT NOT NULL,
      name TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS webpages (
      url TEXT NOT NULL,
      hash TEXT NOT NULL,
      parsed INTEGER NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS events (
      url TEXT NOT NULL,
      hash TEXT NOT NULL,
      article TEXT NOT NULL,
      intro TEXT NOT NULL,
      approved INTEGER,
      intent TEXT,
      probability REAL
    )`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
}
