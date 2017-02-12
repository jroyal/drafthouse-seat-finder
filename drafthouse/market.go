package drafthouse

import (
	"sort"
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

func (m *Market) GetSimpleFilms(date time.Time, cinemaFilter string) []SimpleFilm {
	targetDay := m.getDate(date)
	films := targetDay.GetFilms(cinemaFilter)
	sort.Sort(films)
	return films
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
