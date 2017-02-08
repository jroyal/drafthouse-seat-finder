package drafthouse

import "github.com/jroyal/drafthouse-seat-finder/utils"

type Market struct {
	FeedGenerated     string `json:"FeedGenerated"`
	SessionsGenerated string `json:"SessionsGenerated"`
	MarketID          string `json:"MarketId"`
	MarketName        string `json:"MarketName"`
	MarketSlug        string `json:"MarketSlug"`
	Dates             []Date `json:"Dates"`
}

func (m *Market) MoviesShowingToday() []string {
	movies := utils.NewStringSet()
	today := m.Dates[0]

	for _, cinema := range today.Cinemas {
		movies.AddSlice(cinema.GetFilmNames())
	}
	return movies.ToSlice()
}
