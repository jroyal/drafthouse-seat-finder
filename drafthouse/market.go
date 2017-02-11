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

func (m *Market) Movies(date time.Time, cinemaFilter string) []string {
	targetDay := m.getDate(date)
	movies := targetDay.getMovies(cinemaFilter)
	sort.Strings(movies)
	return movies
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
