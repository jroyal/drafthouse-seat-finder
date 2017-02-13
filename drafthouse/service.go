package drafthouse

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// HandleIndex is the handler for GET /
func HandleIndex(c echo.Context) error {
	market := getMarketInfo()

	dayFilter, _ := time.Parse(apiFormat, time.Now().Format(apiFormat))
	cinemaFilter := ""
	indexTemplate := IndexTemplate{
		Dates:   market.GetDates(),
		Films:   market.GetSimpleFilms(dayFilter, cinemaFilter),
		Cinemas: market.GetCinemas(),
	}
	log.WithFields(log.Fields{
		"path":   c.Path(),
		"scheme": c.Scheme(),
		"method": c.Request().Method,
	}).Info("Request Received")

	return c.Render(http.StatusOK, "index", indexTemplate)
}

// HandleSeats is the handler for POST /seats
func HandleSeats(c echo.Context) error {
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

	market := getMarketInfo()
	seatTemplate := SeatPickerTemplate{
		Films: market.GetFilmSessions(film, dayFilter, cinemaFilter),
	}

	return c.Render(http.StatusOK, "seats", seatTemplate)
}

// HandleGetMovies is the handler for GET /movies
func HandleGetMovies(c echo.Context) error {

	day := c.QueryParam("day")
	if day == "" {
		day = time.Now().Format(apiFormat)
	}
	dayFilter, _ := time.Parse(apiFormat, day)
	cinemaFilter := c.QueryParam("cinema")
	log.WithFields(log.Fields{
		"path":      c.Path(),
		"dayFilter": day,
		"cinemas":   cinemaFilter,
		"scheme":    c.Scheme(),
		"method":    c.Request().Method,
	}).Info("Request Received")
	market := getMarketInfo()
	response := ResponseMovies{market.GetFilmNames(dayFilter, cinemaFilter)}
	return c.JSON(http.StatusOK, response)
}

// HandleGetSingleMovie is the handler for GET /movies/:film-slug
func HandleGetSingleMovie(c echo.Context) error {

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
	market := getMarketInfo()
	response := ResponseMovieTimes{market.GetFilmTimes(filmSlug, dayFilter, cinemaFilter)}
	return c.JSON(http.StatusOK, response)
}
