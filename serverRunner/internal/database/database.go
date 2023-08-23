package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func Init() {
	// TODO(jaegyu): load the dbpath from a config file at
	// /etc/server-runner/sr.conf

	// ! TODO(jaegyu): change `app` with app name
	db, err := sql.Open("sqlite3", "/opt/app/data/app.db")
	if err != nil {
		log.Fatal("Could not open db: ", err)
	}

	Db = db
}

func Sync() {
	instanceTable := `
	CREATE TABLE IF NOT EXISTS instances (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		dir TEXT 
	);
	`

	_, err := Db.Exec(instanceTable)
	if err != nil {
		log.Fatal("Could not create db: ", err)
	}
}
