package main

import (
	"net/http"

	"github.com/jroyal/drafthouse-seat-finder/drafthouse"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Drafthouse Seat Finder")
	})
	e.GET("/moviesshowingtoday", func(c echo.Context) error {
		return c.JSON(http.StatusOK, drafthouse.GetMoviesShowingToday())
	})
	e.Logger.Fatal(e.Start(":8080"))
}
