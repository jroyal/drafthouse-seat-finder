package drafthouse

import (
	"sort"
	"sync"
	"time"
)

type Market struct {
	FeedGenerated     string `json:"FeedGenerated"`
	SessionsGenerated string `json:"SessionsGenerated"`
	MarketID          string `json:"MarketId"`
	MarketName        string `json:"MarketName"`
	MarketSlug        string `json:"MarketSlug"`
	Dates             []Date `json:"Dates"`
}

func (m *Market) GetSimpleFilms(date time.Time, collector *Collector, cinemaFilter string, preload bool) []SimpleFilm {
	targetDay := m.getDate(date)
	films := targetDay.GetFilms(cinemaFilter)
	sort.Sort(films)
	var wg sync.WaitGroup
	for i := range films {
		film := &films[i]
		if preload {
			if i == 0 {
				// Preload the first film
				getFilmMetaData(collector, film)
			} else {
				go getFilmMetaData(collector, film)
			}
		} else {
			wg.Add(1)
			go func(collector *Collector, film *SimpleFilm) {
				defer wg.Done()
				getFilmMetaData(collector, film)
			}(collector, film)
		}

	}
	if !preload {
		wg.Wait()
	}

	return films
}

func getFilmMetaData(collector *Collector, film *SimpleFilm) {
	results := collector.GetFilmMetaData(film.FilmName, film.FilmYear)
	film.FilmPosterURL = results.PosterURL
	film.FilmDescription = results.Description
}

func (m *Market) GetFilmNames(date time.Time, cinemaFilter string) []string {
	targetDay := m.getDate(date)
	movies := targetDay.GetFilmNames(cinemaFilter)
	sort.Strings(movies)
	return movies
}

func (m *Market) GetFilmTimes(filmSlug string, date time.Time, cinemaFilter string) map[string][]string {
	targetDay := m.getDate(date)
	filmTimes := targetDay.GetFilmTimes(filmSlug, cinemaFilter)
	return filmTimes
}

func (m *Market) GetFilmSessions(filmSlug string, date time.Time, cinemaFilter string) []FilmSession {
	targetDay := m.getDate(date)
	filmSessions := targetDay.GetFilmSessions(filmSlug, cinemaFilter)
	return filmSessions
}

func (m *Market) GetCinemas() []SimpleCinema {
	// I choose the next day to make sure we get all theaters
	date := m.Dates[1]
	return date.GetCinemas()
}

func (m *Market) GetDates() []string {
	var dates []string
	for _, date := range m.Dates {
		dates = append(dates, date.Date)
	}
	return dates
}

func (m *Market) getDate(day time.Time) Date {
	for _, date := range m.Dates {
		if date.convertToTime() == day {
			return date
		}
	}
	// TODO: Properly handle failing to find the date
	return Date{}
}
