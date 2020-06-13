
func FetchGet(url string, msg *interface{}) {
}

func Crypto() {
	type Item struct {
		Symbol  string  `json:"s"`
		Current float64 `json:"c"`
	}

	type Currency struct {
		Data []Item `json:data`
	}

	var msg Currency

	url := "https://www.binance.com/exchange-api/v1/public/asset-service/product/get-products"

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

	count := 0

	for _, v := range msg.Data {
		if contains([]interface{}{"BTCUSDT", "ETHUSDT", "LTCUSDT"}, v.Symbol) {
			fmt.Printf("%-10s: %.2f\n", v.Symbol, v.Current)
			count++
		}

		if count == 3 {
			break
		}
	}

}