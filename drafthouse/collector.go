package drafthouse

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Collector struct {
	Client http.Client
	Cache  *Cache
}

func (c *Collector) getJson(url string, target interface{}) error {
	r, err := c.Client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (c *Collector) GetMarketInfo() Market {
	if data, exists := c.Cache.Get("market"); exists {
		return data.(Market)
	}

	log.WithField("market", "austin").Info("Getting Feed Information")
	resp := &MarketResponse{}
	c.getJson(austinMarketFeed, resp)
	c.Cache.Set("market", resp.Market)
	return resp.Market
}

func (c *Collector) GetFilmSeats(film FilmSession) SeatResponse {
	cacheKey := fmt.Sprintf("film_session_%s_%s", film.CinemaID, film.SessionID)
	if data, exists := c.Cache.Get(cacheKey); exists {
		return data.(SeatResponse)
	}
	var seatResponse SeatResponse
	log.WithFields(log.Fields{
		"filmName":    film.FilmName,
		"sessionTime": film.SessionTime,
		"sessionID":   film.SessionID,
	}).Info("Getting Film Information")
	url := fmt.Sprintf("https://drafthouse.com/s/vista/wsVistaWebClient/RESTData.svc/cinemas/%s/sessions/%s/seat-plan", film.CinemaID, film.SessionID)
	c.getJson(url, &seatResponse)
	c.Cache.Set(cacheKey, seatResponse)
	return seatResponse
}
