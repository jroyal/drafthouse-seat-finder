package drafthouse

import (
	"encoding/json"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getMarketInfo() Market {
	market := MarketResponse{}
	getJson(austinMarketFeed, &market)
	return market.Market
}
