package services

import (
	"github.com/nuwak/go_finance/src/db"
)

type PortfolioItem struct {
	Price    float64
	Volume   float64
	Currency Currency
}

type PortfolioStruct struct {
	QValue       string
	QInsertValue string
	Total        map[Currency]map[string]float64
}

var portfolio *PortfolioStruct

func Portfolio() *PortfolioStruct {
	if portfolio == nil {
		portfolio = &PortfolioStruct{
			QValue: "select avg(open_price) as price, sum(volume) as volume, currency from portfolio where symbol = ? and is_close = 0",
			Total:  make(map[Currency]map[string]float64),
		}
	}

	return portfolio
}

func (portfolio *PortfolioStruct) GetValue(symbol *string) (*PortfolioItem, error) {
	item := &PortfolioItem{}

	err := db.DB.QueryRow(portfolio.QValue, *symbol).Scan(&item.Price, &item.Volume, &item.Currency)
	if err != nil {
		return item, err
	}

	return item, nil
}

type Profit struct {
	PriceDiff     float64
	Volume        float64
	CurrentValue  float64
	ValueDiff     float64
	ProfitPercent float64
	Currency      Currency
}

func (portfolio *PortfolioStruct) CalcProfitItem(prev *PortfolioItem, current *float64) *Profit {

	profit := &Profit{}
	profit.Currency = prev.Currency
	profit.PriceDiff = *current - prev.Price
	profit.Volume = prev.Volume * prev.Price
	profit.CurrentValue = prev.Volume * *current
	profit.ValueDiff = profit.CurrentValue - profit.Volume
	profit.ProfitPercent = profit.ValueDiff / profit.Volume * 100

	return profit
}

func (portfolio *PortfolioStruct) CalcTotal(profit *Profit) {
	if len(portfolio.Total[profit.Currency]) == 0 {
		portfolio.Total[profit.Currency] = make(map[string]float64)
	}
	portfolio.Total[profit.Currency]["valueDiff"] += profit.ValueDiff
	portfolio.Total[profit.Currency]["volume"] += profit.Volume
	portfolio.Total[profit.Currency]["profitPercent"] =
		portfolio.Total[profit.Currency]["valueDiff"] / portfolio.Total[profit.Currency]["volume"] * 100
}
