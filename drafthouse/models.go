package drafthouse

type ResponseFilms struct {
	Films []SimpleFilm `json:"films"`
}

type ResponseMovieTimes struct {
	Times map[string][]string `json:"times"`
}

type MarketResponse struct {
	Market Market `json:"Market"`
}

type IndexTemplate struct {
	BaseUrl  string
	IndexUrl string
	Dates    []string
	Films    []SimpleFilm
	Cinemas  []SimpleCinema
}

type SeatPickerTemplate struct {
	BaseUrl  string
	IndexUrl string
	Cinemas  map[string][]FilmSession
}
