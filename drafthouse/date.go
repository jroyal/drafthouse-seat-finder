package drafthouse

import (
	"strings"
	"time"

	"github.com/jroyal/drafthouse-seat-finder/utils"
	log "github.com/sirupsen/logrus"
)

type Date struct {
	DateID  string   `json:"DateId"`
	Date    string   `json:"Date"`
	Cinemas []Cinema `json:"Cinemas"`
}

func (d *Date) convertToTime() time.Time {
	date, err := time.Parse(dateFormat, d.Date)
	if err != nil {
		log.WithField("err", err).Error("Failed to properly convert date")
	}
	return date
}

func (d *Date) GetFilmNames(cinemaFilter string) []string {
	movies := utils.NewStringSet()
	for _, cinema := range d.filterCinemas(cinemaFilter) {
		movies.AddSlice(cinema.GetFilmNames())
	}
	return movies.ToSlice()
}

func (d *Date) GetFilmTimes(filmSlug string, cinemaFilter string) map[string][]string {
	filmTimes := map[string][]string{}
	for _, cinema := range d.filterCinemas(cinemaFilter) {
		times := cinema.GetFilmTimes(filmSlug)
		if times != nil {
			filmTimes[cinema.CinemaName] = times
		}
	}
	return filmTimes
}

func (d *Date) GetCinemas() []SimpleCinema {
	var cinemas []SimpleCinema

	for _, cinema := range d.Cinemas {
		cinemas = append(cinemas, SimpleCinema{
			CinemaName: cinema.CinemaName,
			CinemaSlug: cinema.CinemaSlug,
		})
	}
	return cinemas
}

func (d *Date) GetFilms(cinemaFilter string) SimpleFilms {
	var filmResult []SimpleFilm
	filmSet := utils.NewStringSet()
	for _, cinema := range d.filterCinemas(cinemaFilter) {
		films := cinema.GetFilms()
		for _, film := range films {
			if !filmSet.Contains(film.FilmName) {
				filmSet.Add(film.FilmName)
				filmResult = append(filmResult, film)
			}
		}
	}
	return filmResult
}

func (d *Date) GetFilmSessions(filmSlug string, cinemaFilter string) []FilmSession {
	var filmSessions []FilmSession
	for _, cinema := range d.filterCinemas(cinemaFilter) {
		filmSessions = append(filmSessions, cinema.GetFilmSessions(filmSlug)...)
	}
	return filmSessions
}

func (d *Date) filterCinemas(cinemaFilter string) []Cinema {
	if cinemaFilter == "" {
		return d.Cinemas
	}

	var results []Cinema
	cinemas := strings.Split(cinemaFilter, ",")
	filter := utils.NewStringSet()
	filter.AddSlice(cinemas)
	for _, cinema := range d.Cinemas {
		if filter.Contains(cinema.CinemaSlug) {
			results = append(results, cinema)
		}
	}
	return results
}
