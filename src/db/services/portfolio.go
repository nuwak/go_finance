package services

import (
	"errors"
	"fmt"
	"github.com/nuwak/go_finance/src/db"
	"github.com/nuwak/go_finance/src/libs/mat"
	"log"
)

type PortfolioItem struct {
	id       int
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

func (portfolio *PortfolioStruct) GetResult(symbol *string) (*PortfolioItem, error) {
	rows, err := db.DB.Query("select id, volume, open_price, currency from portfolio where symbol = ? and is_close = 0", symbol)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var total float64
	res := &PortfolioItem{}
	res.Price = 0
	for rows.Next() {
		row := &PortfolioItem{}
		if err := rows.Scan(&row.id, &row.Volume, &row.Price, &row.Currency); err != nil {
			return res, err
		}

		res.Currency = row.Currency
		res.Volume += row.Volume
		total += row.Volume * row.Price
	}

	if res.Volume == 0 {
		return res, errors.New("None item in portfolio")
	}
	res.Price = total / res.Volume

	return res, nil
}

func (portfolio *PortfolioStruct) Buy(symbol *string, price *float64, volume *float64, currency *Currency) {
	stmt, err := db.DB.Prepare("INSERT INTO portfolio(symbol, open_price, volume, currency) VALUES( ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(symbol, price, volume, currency)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.LastInsertId())
}

func (portfolio *PortfolioStruct) Sell(symbol *string, price *float64, volume *float64) {
	rows, err := db.DB.Query("select id, symbol, volume, open_price, open_dt, currency from main.portfolio where symbol = ? and is_close = 0", symbol)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	sell := struct {
		total    float64
		idEqual  int64
		idMore   int64
		volume   float64
		price    float64
		openDt   string
		currency Currency
	}{}

	for rows.Next() {
		row := struct {
			id       int64
			name     string
			volume   float64
			price    float64
			openDt   string
			currency Currency
		}{}
		if err := rows.Scan(&row.id, &row.name, &row.volume, &row.price, &row.openDt, &row.currency); err != nil {
			log.Fatal(err)
		}

		if row.volume == *volume {
			sell.idEqual = row.id
		} else if row.volume > *volume {
			sell.idMore = row.id
		}

		sell.volume = row.volume
		sell.price = row.price
		sell.openDt = row.openDt
		sell.currency = row.currency

		fmt.Println(row)

		sell.total += row.volume
	}

	if *volume > sell.total {
		log.Fatalln("Not enough Volume")
	}

	fmt.Println(sell)

	if sell.idEqual != 0 {
		stmt, err := db.DB.Prepare(`
		UPDATE portfolio SET 
			close_dt = datetime('now', '+3 hour'),
			is_close = 1,
			close_price = ?,
			update_dt =  datetime('now', '+3 hour')
		WHERE id = ?
	`)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(price, sell.idEqual)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Profit: ", mat.RoundAcc(*price-sell.price, 8))
		fmt.Println(res.RowsAffected())
	} else if sell.idMore != 0 {
		stmt, err := db.DB.Prepare(`
		UPDATE portfolio SET 
			close_dt = datetime('now', '+3 hour'),
			is_close = 1,
			close_price = ?,
			volume = ?,
			update_dt =  datetime('now', '+3 hour')
		WHERE id = ?
	`)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		volumeDiff := sell.volume - *volume
		fmt.Println(volumeDiff)
		res, err := stmt.Exec(price, volume, sell.idMore)
		if err != nil {
			log.Fatal(err)
		}

		stmt, err = db.DB.Prepare(`
			INSERT INTO portfolio(symbol, open_price, volume, currency, open_dt)
			VALUES( ?, ?, ?, ?, ?)
		`)
		if err != nil {
			log.Fatal(err)
		}

		res, err = stmt.Exec(symbol, sell.price, volumeDiff, sell.currency, sell.openDt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res)
	} else {
		//	itodo: описать случай когда ни одна запись целиком не может покрыть продажу и нужно закрывать несколько записей
	}
	//	itodo: описать случай когда продается все
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
