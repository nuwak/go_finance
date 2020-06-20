package parsers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nuwak/go_finance/src/libs"
)

type Item struct {
	Symbol  string  `json:"s"`
	Current float64 `json:"c"`
}

type Currency struct {
	Data []Item `json:data`
}

func FetchGet(url string, msg *Currency) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &msg)
}

func BinanceGet() {
	var msg Currency
	url := "https://www.binance.com/exchange-api/v1/public/asset-service/product/get-products"

	FetchGet(url, &msg)

	count := 0

	for _, v := range msg.Data {
		if libs.Contains([]interface{}{"BTCUSDT", "ETHUSDT", "LTCUSDT"}, v.Symbol) {

			libs.Print(&v.Symbol, &v.Current)
			count++
		}

		if count == 3 {
			break
		}
	}

}
