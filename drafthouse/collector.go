package drafthouse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	tmdb "github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
)

type Collector struct {
	Client http.Client
	Cache  *Cache
}

type MetaDataResults struct {
	PosterURL   string
	Description string
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

func (c *Collector) GetFilmMetaData(filmName string, filmYear string) *MetaDataResults {
	// Things to replace quickly in order to get better answers
	r := regexp.MustCompile("\\([0-9]{4}\\)|2D|3D|\\(Subtitled\\)|\\(Dubbed\\)")
	filmName = r.ReplaceAllString(filmName, "")

	cacheKey := fmt.Sprintf("film_meta_%s", filmName)
	if data, exists := c.Cache.Get(cacheKey); exists {
		return data.(*MetaDataResults)
	}

	tmdbClient := tmdb.Init("872b1c79febd6e43d7e17f8bffb898ff")
	log.WithFields(log.Fields{
		"filmName": filmName,
		"filmYear": filmYear,
	}).Info("Getting Film Meta Data")
	options := map[string]string{}

	results, _ := tmdbClient.SearchMovie(filmName, options)
	metaResults := &MetaDataResults{}
	if len(results.Results) > 0 {
		posterBaseURL := "https://image.tmdb.org/t/p/w185"
		metaResults.Description = results.Results[0].Overview
		metaResults.PosterURL = fmt.Sprintf("%s/%s", posterBaseURL, results.Results[0].PosterPath)
	} else {
		// Check the drafthouse for the information
	}
	c.Cache.SetWithExpiration(cacheKey, metaResults, 168*time.Hour)
	// TODO: Handle the error
	return metaResults
}
