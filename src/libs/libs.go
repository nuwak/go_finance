package libs

import (
	"fmt"
	"math"
	"time"

	"main"

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
	// database, err := sql.Open("sqlite3", "./fin.db")
	// if err != nil {
	// 	panic(err)
	// }
	var yesterdayVal float64
	var todayVal float64

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	main.DB.
		QueryRow("SELECT value FROM history where symbol = ? and first = ? order by first desc limit 1", *name, yesterday).
		Scan(&yesterdayVal)
	database.
		QueryRow("SELECT value FROM history where symbol = ? and first = ? order by first desc limit 1", *name, today).
		Scan(&todayVal)

	if todayVal == 0 {
		statement, _ := database.Prepare("INSERT INTO history (symbol, value, first) VALUES (?, ?, ?)")
		statement.Exec(*name, *value, time.Now().Format("2006-01-02"))
	} else if yesterdayVal == 0 {
		fmt.Printf("%-10s: %g\n", *name, *value)
	} else {
		var diff float64 = math.Round((*value-yesterdayVal)*100) / 100
		var change float64 = (math.Round(diff / *value * 10000)) / 100
		fmt.Printf("%-10s: %-10g  | %8g | %6g \n", *name, *value, diff, change)
	}
}
