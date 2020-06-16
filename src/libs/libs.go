package libs

import (
	"fmt"
	"github.com/nuwak/go_finance/src/db/services"
	"math"

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

func Print(symbol *string, value *float64) {
	var diff float64
	var change float64
	total := make(map[string]float64)

	yesterdayVal := services.History().GetValue(symbol, false)
	todayVal := services.History().GetValue(symbol, true)

	if todayVal == 0 {
		services.History().AddValue(symbol, value)
	}

	if yesterdayVal == 0 {
		if todayVal == 0 {
			fmt.Printf("%-10s: %g\n", *symbol, *value)
		} else {
			diff = math.Round((*value-todayVal)*100) / 100
			change = (math.Round(diff / *value * 10000)) / 100
			fmt.Printf("%-10s: %-10g  | %8g | %6g \n", *symbol, *value, diff, change)
		}
	} else {
		diff = math.Round((*value-yesterdayVal)*100) / 100
		change = (math.Round(diff / *value * 10000)) / 100
		if portfolioItem, err := services.Portfolio().GetValue(symbol); err == nil {
			portfolio := services.Portfolio().CalcProfitItem(portfolioItem, value)
			services.Portfolio().CalcTotal(portfolio)

			total["profitPercent"] += portfolio["profitPercent"]
			total["valueDiff"] += portfolio["valueDiff"]
			total["volume"] += portfolio["volume"]

			fmt.Printf(
				"%-10s: %-10g  | %8g | %6g | %6g | %8g | %8g | %8g\n",
				*symbol,
				*value,
				diff,
				change,
				Round(portfolio["profitPercent"]),
				Round(portfolio["valueDiff"]),
				Round(portfolio["volume"]),
				Round(portfolio["priceDiff"]),
			)
		} else {
			fmt.Printf("%-10s: %-10g  | %8g | %6g \n", *symbol, *value, diff, change)
		}
	}
}

func Round(val float64) float64 {
	return math.Round(val*100) / 100
}
