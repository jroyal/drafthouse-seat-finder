package drafthouse

type Market struct {
	FeedGenerated     string `json:"FeedGenerated"`
	SessionsGenerated string `json:"SessionsGenerated"`
	MarketID          string `json:"MarketId"`
	MarketName        string `json:"MarketName"`
	MarketSlug        string `json:"MarketSlug"`
	Dates             []Date `json:"Dates"`
}

func (m *Market) MoviesShowingToday() []string {
	var movies []string
	today := m.Dates[0]

	for _, cinema := range today.Cinemas {
		movies = append(movies, cinema.GetFilmNames()...)
	}
	return movies
}
