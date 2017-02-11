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

	return c.JSON(http.StatusOK, GetMovies(dayFilter, cinemaFilter))
}
