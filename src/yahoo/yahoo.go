package yahoo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nuwak/go_finance/src/libs"
)

type Api struct {
	urlTemplate string
	symbol      string
	alias       string
	name        string
	value       float64
	data        struct {
		Chart struct {
			Result []struct {
				Meta struct {
					RegularMarketPrice float64 `json:"regularMarketPrice"`
				} `json:"meta"`
			} `json:"result"`
		} `json:"chart"`
	}
}

func NewApi(symbol string, alias string) Api {
	api := Api{}
	api.urlTemplate = "https://query1.finance.yahoo.com/v8/finance/chart/?symbol=%s&period1=%s&period2=%s&interval=1m"
	api.symbol = symbol
	api.alias = alias
	api.name = api.Name()

	return api
}

func (api *Api) getURL() string {

	return fmt.Sprintf(
		api.urlTemplate,
		url.QueryEscape(api.symbol),
		strconv.FormatInt(time.Now().Unix()-200, 10),
		strconv.FormatInt(time.Now().Unix()-30, 10),
	)
}

func (api *Api) Name() string {
	if api.alias != "" {
		return api.alias
	}

	return api.symbol
}

func (api *Api) getData() {
	req := api.getURL()
	resp, err := http.Get(req)

	defer resp.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &api.data)

	api.value = api.Value()
}

func (api *Api) Value() float64 {
	return api.data.Chart.Result[0].Meta.RegularMarketPrice
}

func FromChart(symbol string, alias string) {
	api := NewApi(symbol, alias)
	api.getData()

	libs.Print(&api.name, &api.value)
}
