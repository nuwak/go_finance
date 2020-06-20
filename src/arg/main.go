package arg

import (
	"flag"
	"github.com/nuwak/go_finance/src/db/services"
	"github.com/nuwak/go_finance/src/parsers"
	"log"
)

func Args() {
	get := flag.Bool("get", false, "get quote")
	buy := flag.String("buy", "", "symbol")
	sell := flag.String("sell", "", "symbol")
	price := flag.Float64("p", 0, "price")
	volume := flag.Float64("v", 0, "volume")
	flag.Parse()

	if *get == true {
		parsers.GetQuote()
		parsers.PrintQuote()
	}

	if *buy != "" && (*price == 0 || *volume == 0) {
		log.Fatal("need price or volume")
	}

	if *buy != "" {
		item, err := services.Symbols().Get(buy)
		if err != nil {
			log.Fatal("No symbol")
		}

		services.Portfolio().Buy(buy, price, volume, &item.Currency)
	}

	if *sell != "" {
		item, err := services.Symbols().Get(sell)
		if err != nil {
			log.Fatal("No symbol")
		}

		services.Portfolio().Sell(&item.Symbol, price, volume)
	}
}
