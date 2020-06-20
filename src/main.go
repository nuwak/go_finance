package main

import (
	"github.com/nuwak/go_finance/src/arg"
	"github.com/nuwak/go_finance/src/db"
)

func main() {
	db.InitDB()
	arg.Args()
}
