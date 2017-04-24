package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jroyal/drafthouse-seat-finder/drafthouse"
	"github.com/labstack/echo"
	"github.com/namsral/flag"
	tmdb "github.com/ryanbradynd05/go-tmdb"
	log "github.com/sirupsen/logrus"
)

// Command Line Options
var (
	port       = flag.String("port", "8080", "Port to run the server on")
	local      = flag.Bool("local", false, "Run the server only on localhost")
	urlBase    = flag.String("urlBase", "", "For reverse proxy support")
	tmdbAPIKey = flag.String("tmdbAPIKey", "", "API key needed to talk to TMDB")
)

var cacheTTL = 300 * time.Second

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

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = t

	base := strings.Trim(*urlBase, "/")
	index := "/" + base
	if base != "" {
		base = "/" + base
	}

	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("baseUrl", base)
			c.Set("indexUrl", index)
			return h(c)
		}
	})

	collector := &drafthouse.Collector{
		Client:     http.Client{Timeout: 10 * time.Second},
		Cache:      drafthouse.NewCache(cacheTTL),
		TMDBClient: tmdb.Init(*tmdbAPIKey),
	}

	config := &drafthouse.DrafthouseServiceConfig{
		Index: index,
		Base:  base,
	}

	drafthouse.Service(e, collector, config)

	if *local {
		e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%s", *port)))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
	}

}
