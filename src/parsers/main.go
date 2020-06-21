package parsers

import (
	"fmt"
	"github.com/nuwak/go_finance/src/db/services"
	"github.com/nuwak/go_finance/src/libs/mat"
	"log"
	"time"
)

func GetQuote() {
	p := fmt.Println
	weekday := time.Now().Weekday()
	now := time.Now()
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	openNY := time.Date(now.Year(), now.Month(), now.Day(), 16, 30, 0, 0, location)

	ForexGet("USD RUB", "USD/RUB")
	YahooGet("BZQ20.NYM", "BRENT")
	YahooGet("ROSN.ME", "ROSN.ME")
	YahooGet("AFKS.ME", "AFKS.ME")
	BinanceGet()
	YahooGet("ZM", "")
	YahooGet("NVDA", "")

	if int(weekday) != 0 && int(weekday) != 6 && openNY.After(now) {
		p("future")
		YahooGet("ES=F", "S&P500")
		YahooGet("NQ=F", "NDAQ")
	} else {
		p("now")
		YahooGet("^GSPC", "S&P500")
		YahooGet("^IXIC", "NDAQ")
	}
}

func PrintQuote() {
	fmt.Printf(
		"%52g | %8g | %8g \n",
		mat.Round(services.Portfolio().Total[services.USD]["profitPercent"]),
		mat.Round(services.Portfolio().Total[services.USD]["valueDiff"]),
		mat.Round(services.Portfolio().Total[services.USD]["volume"]),
	)

	fmt.Printf(
		"%52g | %8g | %8g \n",
		mat.Round(services.Portfolio().Total[services.RUB]["profitPercent"]),
		mat.Round(services.Portfolio().Total[services.RUB]["valueDiff"]),
		mat.Round(services.Portfolio().Total[services.RUB]["volume"]),
	)
}
