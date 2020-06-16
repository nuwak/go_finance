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
	Total        map[string]float64
}

var portfolio *PortfolioStruct

func Portfolio() *PortfolioStruct {
	if portfolio == nil {
		portfolio = &PortfolioStruct{
			QValue: "select avg(open_price) as price, sum(volume) as volume from portfolio where symbol = ? and is_close = 0",
			Total:  make(map[string]float64),
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

func (portfolio *PortfolioStruct) CalcProfitItem(prev *PortfolioItem, current *float64) map[string]float64 {
	res := make(map[string]float64)
	res["priceDiff"] = *current - prev.Price
	res["volume"] = prev.Volume * prev.Price
	res["currentValue"] = prev.Volume * *current
	res["valueDiff"] = res["currentValue"] - res["volume"]
	res["profitPercent"] = res["valueDiff"] / res["volume"] * 100

	return res
}

func (portfolio *PortfolioStruct) CalcTotal(res map[string]float64) {
	portfolio.Total["valueDiff"] += res["valueDiff"]
	portfolio.Total["volume"] += res["volume"]
	portfolio.Total["profitPercent"] = portfolio.Total["valueDiff"] / portfolio.Total["volume"] * 100
}
