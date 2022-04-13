package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetConnection() *sql.DB {
	if db != nil {
		return db
	}

	var err error

	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	return db
}

func MakeMigrations() error {
	db := GetConnection()

	q := `CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(64) NULL,
		description VARCHAR(200) NULL,
		create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		update_at TIMESTAMP NOT NULL );`
	_, err := db.Exec(q)

	if err != nil {
		return err
	}
	return nil
}
