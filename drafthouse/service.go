package drafthouse

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

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
	response := ResponseMovies{market.Movies(dayFilter, cinemaFilter)}
	return c.JSON(http.StatusOK, response)
}

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
