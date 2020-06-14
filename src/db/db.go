package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("sqlite3", "./fin.db")

	if err != nil {
		panic(err)
	}

	statement, err := DB.Prepare(`
	CREATE TABLE IF NOT EXISTS history
		(
			id INTEGER PRIMARY KEY,
			symbol CHARACTER,
			value DECIMAL(10,2),
			first DATE
		)
	`)

	if err != nil {
		panic(err)
	}

	statement.Exec()
}
