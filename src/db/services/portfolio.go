package services

import (
	"github.com/nuwak/go_finance/src/db"
)

type PortfolioItem struct {
	Price  float64
	Volume float64
}

type PortfolioStruct struct {
	QValue       string
	QInsertValue string
}

var portfolio *PortfolioStruct

func Portfolio() *PortfolioStruct {
	if portfolio == nil {
		portfolio = &PortfolioStruct{
			QValue: "select avg(open_price) as price, sum(volume) as volume from portfolio where symbol = ? and is_close = 0",
		}
	}

	return portfolio
}

func (portfolio *PortfolioStruct) GetValue(symbol *string) (*PortfolioItem, error) {
	item := &PortfolioItem{}

	err := db.DB.QueryRow(portfolio.QValue, *symbol).Scan(&item.Price, &item.Volume)
	if err != nil {
		return item, err
	}

	return item, nil
}
