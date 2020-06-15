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
    id        INTEGER
        primary key,
    symbol    CHARACTER,
    value     DECIMAL(10, 2),
    created_d DATE default current_date,
    created_dt DATE default current_timestamp
); 

CREATE TABLE IF NOT EXISTS portfolio
(
    id          INTEGER
        primary key autoincrement,
    symbol      CHARACTER      not null,
    open_price  DECIMAL(16, 8) not null,
    volume      DECIMAL(16, 8) not null,
    open_dt     DATETIME default current_timestamp,
    close_dt    DATETIME,
    is_close    BOOLEAN  default 0,
    close_price DECIMAL(16, 8),
    update_dt   DATETIME
);

	`)

	if err != nil {
		panic(err)
	}

	statement.Exec()
}
