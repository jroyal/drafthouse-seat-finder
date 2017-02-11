package drafthouse

type ResponseMovies struct {
	Movies []string `json:"movies"`
}

type ResponseMovieTimes struct {
	Times map[string][]string `json:"times"`
}

type MarketResponse struct {
	Market Market `json:"Market"`
}
