package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"igor.local/go/financ/src/binance"
	"igor.local/go/financ/src/libs"
	"igor.local/go/financ/src/yahoo"
)

func main() {
	initDB()
	// USD()
	Forex("USD RUB", "USD/RUB")
	yahoo.FromChart("BZQ20.NYM", "BRENT")
	binance.Crypto()
	yahoo.FromChart("ZM", "")
	yahoo.FromChart("NVDA", "")
	yahoo.FromChart("^GSPC", "S&P500")
}

func initDB() {
	database, err := sql.Open("sqlite3", "./fin.db")

	if err != nil {
		panic(err)
	}

	statement, err := database.Prepare(`
	CREATE TABLE IF NOT EXISTS history
		(
			id INTEGER PRIMARY KEY,
			symbol CHARACTER,
			value DECIMAL(10,2),
			first DATE
		)
	`)

	if err != nil {
		panic(err)
	}

	statement.Exec()
	// statement, _ = database.Prepare("INSERT INTO history (symbol, value, first) VALUES (?, ?, ?)")
	// statement.Exec("BTCUSDT", 9400.03, "2020-01-02")
	// rows, _ := database.Query("SELECT * FROM history")
	// var id int
	// var symbol string
	// var value float64
	// var first string
	// for rows.Next() {
	// 	rows.Scan(&id, &symbol, &value, &first)
	// 	fmt.Println(id, ": ", symbol, " ", value, " ", first)
	// }
}

func USD() {
	type Curr struct {
		Code string  `json:code`
		Rate float32 `json:rate`
	}

	var msg []Curr

	resp, err := http.Get("https://marketdata-marketplace.moex.com/api/securities?category=currencies")

	defer resp.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &msg)
	for _, v := range msg {
		if v.Code == "USD/RUB" {
			fmt.Printf("%-10s: %.2f\n", "USDRUB", v.Rate)
			break
		}
	}
}

func Forex(symbol string, alias string) {
	var name string

	if alias != "" {
		name = alias
	} else {
		name = symbol
	}
	type Message struct {
		SiteID         string   `json:"siteId"`
		ProductType    string   `json:"productType"`
		MarketType     string   `json:"marketType"`
		Products       []string `json:"products"`
		RequiredFields []string `json:"requiredFields"`
	}

	message := &Message{
		SiteID:         "forex.web.us",
		ProductType:    "FX",
		MarketType:     "FX",
		Products:       []string{symbol},
		RequiredFields: []string{"Bid"},
	}

	ress, _ := json.Marshal(message)

	url := "https://www.forex.com/_Srvc/feeds/LiveRates.asmx/GetProductRates"
	method := "POST"

	payload := strings.NewReader(string(ress))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("content-type", "application/json; charset=UTF-8")
	req.Header.Add("Cookie", "ForexNetworkPool_15-sitecore_SCD_pool_443=MKMNACAK; forex_us#lang=en")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var data map[string][]map[string]string

	json.Unmarshal(body, &data)

	val, _ := strconv.ParseFloat(data["d"][0]["Bid"], 64)

	libs.Print(&name, &val)
}
