package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/jroyal/drafthouse-seat-finder/drafthouse"
	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", "World")
	})

	e.GET("/movies", drafthouse.HandleGetMovies)
	e.GET("/movies/:film-slug", drafthouse.HandleGetSingleMovie)
	e.Logger.Fatal(e.Start("localhost:8080"))
}
