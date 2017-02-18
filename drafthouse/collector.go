package drafthouse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

var market = MarketResponse{}

var lastUpdated time.Time

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func isStale() bool {
	age := time.Since(lastUpdated)
	return age.Minutes() > 30
}

func getMarketInfo() Market {
	if market.Market.MarketID == "" || isStale() {
		log.WithField("market", "austin").Info("Getting Feed Information")
		getJson(austinMarketFeed, &market)
		lastUpdated = time.Now()
	}

	return market.Market
}

func getFilmSeats(film FilmSession) SeatResponse {
	var seatResponse SeatResponse
	log.WithFields(log.Fields{
		"filmName":    film.FilmName,
		"sessionTime": film.SessionTime,
		"sessionID":   film.SessionID,
	}).Info("Getting Film Information")
	url := fmt.Sprintf("https://drafthouse.com/s/vista/wsVistaWebClient/RESTData.svc/cinemas/%s/sessions/%s/seat-plan", film.CinemaID, film.SessionID)
	getJson(url, &seatResponse)
	return seatResponse
}
