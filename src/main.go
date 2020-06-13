package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"igor.local/go/financ/src/binance"
	"igor.local/go/financ/src/yahoo"
)

func main() {
	// USD()
	Forex("USD RUB", "USD/RUB")
	yahoo.FromChart("BZQ20.NYM", "BRENT")
	binance.Crypto()
	yahoo.FromChart("ZM", "")
	yahoo.FromChart("NVDA", "")
	yahoo.FromChart("^GSPC", "S&P500")
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

func Yahoo(symbol string, alias string) {

	var name string

	if alias != "" {
		name = alias
	} else {
		name = symbol
	}

	type Symbol struct {
		Chart struct {
			Result []struct {
				Meta struct {
					RegularMarketPrice float32 `json:"regularMarketPrice"`
				} `json:"meta"`
			} `json:"result"`
		} `json:"chart"`
	}

	query := "https://query1.finance.yahoo.com/v8/finance/chart/?symbol=%s&period1=%s&period2=%s&interval=1m"

	req := fmt.Sprintf(
		query,
		url.QueryEscape(symbol),
		strconv.FormatInt(time.Now().Unix()-200, 10),
		strconv.FormatInt(time.Now().Unix()-30, 10),
	)

	resp, err := http.Get(req)

	defer resp.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var msg Symbol
	json.Unmarshal(body, &msg)

	fmt.Printf("%-10s: %g\n", name, msg.Chart.Result[0].Meta.RegularMarketPrice)
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

	fmt.Printf("%-10s: %g\n", name, val)
}

func contains(arr []interface{}, val interface{}) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}
