package libs

import (
	"fmt"
	"github.com/nuwak/go_finance/src/db/services"
	"github.com/nuwak/go_finance/src/libs/mat"
	"math"
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

			fmt.Printf(
				"%-10s: %-10g  | %8g | %6g | %6g | %8g | %8g | %8g\n",
				*symbol,
				*value,
				diff,
				change,
				mat.Round(portfolio.ProfitPercent),
				mat.Round(portfolio.ValueDiff),
				mat.Round(portfolio.Volume),
				mat.Round(portfolio.PriceDiff),
			)
		} else {
			fmt.Printf("%-10s: %-10g  | %8g | %6g \n", *symbol, *value, diff, change)
		}
	}
}
