package main

import (
	"flag"
	"html/template"
	"io"

	"github.com/jroyal/drafthouse-seat-finder/drafthouse"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// Command Line Options
var (
	local = flag.Bool("local", false, "Run the server only on localhost")
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	flag.Parse()
	e := echo.New()
	e.Static("/", "public")

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	// These are the two main routes used by the UI
	e.GET("/", drafthouse.HandleIndex)
	e.POST("/seats", drafthouse.HandleSeats)

	// These are fun convienience routes that I used for testing. Eventually I might clean these out
	e.GET("/films", drafthouse.HandleGetFilms)
	e.GET("/movies/:film-slug", drafthouse.HandleGetSingleMovie)

	if *local {
		e.Logger.Fatal(e.Start("localhost:8080"))
	} else {
		e.Logger.Fatal(e.Start(":8080"))
	}

}
