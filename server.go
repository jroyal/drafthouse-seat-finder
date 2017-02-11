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
	e.GET("/movies", drafthouse.HandleGetMovies)
	e.GET("/movies/:film-slug", drafthouse.HandleGetSingleMovie)
	e.Logger.Fatal(e.Start("localhost:8080"))
}
