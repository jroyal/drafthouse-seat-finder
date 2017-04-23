package drafthouse

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

type DrafthouseService struct {
	collector *Collector
}

type DrafthouseServiceConfig struct {
	Index string
	Base  string
}

// HandleIndex is the handler for GET /
func (s *DrafthouseService) HandleIndex(c echo.Context) error {
	market := s.collector.GetMarketInfo()

	dayFilter, _ := time.Parse(apiFormat, time.Now().Format(apiFormat))
	cinemaFilter := ""
	indexTemplate := IndexTemplate{
		BaseUrl:  c.Get("baseUrl").(string),
		IndexUrl: c.Get("indexUrl").(string),
		Dates:    market.GetDates(),
		Films:    market.GetSimpleFilms(dayFilter, cinemaFilter),
		Cinemas:  market.GetCinemas(),
	}
	log.WithFields(log.Fields{
		"path":   c.Path(),
		"scheme": c.Scheme(),
		"method": c.Request().Method,
	}).Info("Request Received")

	return c.Render(http.StatusOK, "index", indexTemplate)
}

// HandleSeats is the handler for POST /seats
func (s *DrafthouseService) HandleSeats(c echo.Context) error {
	req := c.Request()
	req.ParseForm()
	form := req.Form

	cinemas := form["cinemas"]
	film := c.FormValue("film")
	date := c.FormValue("date")
	dayFilter, _ := time.Parse(dateFormat, date)
	cinemaFilter := strings.Join(cinemas, ",")
	log.WithFields(log.Fields{
		"path":    c.Path(),
		"scheme":  c.Scheme(),
		"method":  c.Request().Method,
		"cinemas": cinemaFilter,
		"film":    film,
		"date":    dayFilter,
	}).Info("Request Recieved")

	market := s.collector.GetMarketInfo()
	filmSessions := market.GetFilmSessions(film, dayFilter, cinemaFilter)
	loadFilmSeats(filmSessions, c.Get("baseUrl").(string), s.collector)
	cinemaFilmSessionMap := sortFilmSessions(filmSessions)
	seatTemplate := SeatPickerTemplate{
		BaseUrl:  c.Get("baseUrl").(string),
		IndexUrl: c.Get("indexUrl").(string),
		Cinemas:  cinemaFilmSessionMap,
	}

	return c.Render(http.StatusOK, "seats", seatTemplate)
}

// HandleGetFilms is the handler for GET /films
func (s *DrafthouseService) HandleGetFilms(c echo.Context) error {

	day := c.QueryParam("day")
	if day == "" {
		day = time.Now().Format(dateFormat)
	}
	dayFilter, _ := time.Parse(dateFormat, day)
	cinemaFilter := c.QueryParam("cinema")
	log.WithFields(log.Fields{
		"path":      c.Path(),
		"dayFilter": day,
		"cinemas":   cinemaFilter,
		"scheme":    c.Scheme(),
		"method":    c.Request().Method,
	}).Info("Request Received")
	market := s.collector.GetMarketInfo()
	response := ResponseFilms{market.GetSimpleFilms(dayFilter, cinemaFilter)}
	return c.JSON(http.StatusOK, response)
}

// HandleGetSingleMovie is the handler for GET /movies/:film-slug
func (s *DrafthouseService) HandleGetSingleMovie(c echo.Context) error {

	filmSlug := c.Param("film-slug")

	day := c.QueryParam("day")
	if day == "" {
		day = time.Now().Format(apiFormat)
	}
	dayFilter, _ := time.Parse(apiFormat, day)
	cinemaFilter := c.QueryParam("cinema")
	log.WithFields(log.Fields{
		"film":      filmSlug,
		"path":      c.Path(),
		"dayFilter": day,
		"cinemas":   cinemaFilter,
		"scheme":    c.Scheme(),
		"method":    c.Request().Method,
	}).Info("Request Received")
	market := s.collector.GetMarketInfo()
	response := ResponseMovieTimes{market.GetFilmTimes(filmSlug, dayFilter, cinemaFilter)}
	return c.JSON(http.StatusOK, response)
}

// Service registers the routes for the drafthouse service
func Service(routes *echo.Echo, collector *Collector, config *DrafthouseServiceConfig) {
	s := DrafthouseService{
		collector: collector,
	}

	routes.Static(fmt.Sprintf("%s", config.Index), "public")

	// These are the two main routes used by the UI
	routes.GET(fmt.Sprintf("%s", config.Index), s.HandleIndex)
	routes.POST(fmt.Sprintf("%s/seats", config.Base), s.HandleSeats)

	// These are fun convienience routes that I used for testing. Eventually I might clean these out
	routes.GET(fmt.Sprintf("%s/films", config.Base), s.HandleGetFilms)
	routes.GET(fmt.Sprintf("%s/movies/:film-slug", config.Base), s.HandleGetSingleMovie)
}
