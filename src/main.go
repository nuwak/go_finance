package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/nuwak/go_finance/src/arg"
	"github.com/nuwak/go_finance/src/db"
)

func main() {
	db.InitDB()
	arg.Args()
}
