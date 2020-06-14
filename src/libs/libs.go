package libs

import (
	"fmt"
	"math"
	"time"

	"github.com/nuwak/go_finance/src/config"
	"github.com/nuwak/go_finance/src/db"

	_ "github.com/mattn/go-sqlite3"
)

func Contains(arr []interface{}, str interface{}) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func Print(name *string, value *float64) {
	var yesterdayVal float64
	var todayVal float64
	var diff float64
	var change float64

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	db.DB.
		QueryRow("SELECT value FROM history where symbol = ? and first = ? order by first desc limit 1", *name, yesterday).
		Scan(&yesterdayVal)
	db.DB.
		QueryRow("SELECT value FROM history where symbol = ? and first = ? order by first desc limit 1", *name, today).
		Scan(&todayVal)

	if todayVal == 0 {
		statement, _ := db.DB.Prepare("INSERT INTO history (symbol, value, first) VALUES (?, ?, ?)")
		statement.Exec(*name, *value, time.Now().Format("2006-01-02"))
	}

	if yesterdayVal == 0 {
		if todayVal == 0 {
			fmt.Printf("%-10s: %g\n", *name, *value)
		} else {
			diff = math.Round((*value-todayVal)*100) / 100
			change = (math.Round(diff / *value * 10000)) / 100
			fmt.Printf("%-10s: %-10g  | %8g | %6g \n", *name, *value, diff, change)
		}
	} else {
		diff = math.Round((*value-yesterdayVal)*100) / 100
		change = (math.Round(diff / *value * 10000)) / 100
		if portfolioItem, ok := config.Portfolio[*name]; ok {
			portfolioDiff := math.Round((*value-portfolioItem.Price)*100) / 100
			portfolioProfit := math.Round(portfolioDiff/portfolioItem.Price*10000) / 100

			fmt.Printf(
				"%-10s: %-10g  | %8g | %6g | %6g | %8g\n",
				*name,
				*value,
				diff,
				change,
				portfolioProfit,
				portfolioDiff,
			)
		} else {
			fmt.Printf("%-10s: %-10g  | %8g | %6g \n", *name, *value, diff, change)
		}
	}
}
