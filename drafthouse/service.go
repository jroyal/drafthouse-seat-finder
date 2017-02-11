package drafthouse

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func HandleGetMovies(c echo.Context) error {

	day := c.QueryParam("day")
	if day == "" {
		day = time.Now().Format(apiFormat)
	}
	dayFilter, _ := time.Parse(apiFormat, day)
	return c.JSON(http.StatusOK, GetMovies(dayFilter))
}
