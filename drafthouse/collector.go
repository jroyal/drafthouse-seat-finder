package drafthouse

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

var market = MarketResponse{}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getMarketInfo() Market {
	if market.Market.MarketID == "" {
		log.WithField("market", "austin").Info("Getting Feed Information")
		getJson(austinMarketFeed, &market)
	}

	return market.Market
}
