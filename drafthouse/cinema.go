package drafthouse

type Cinema struct {
	CinemaID          string `json:"CinemaId"`
	CinemaName        string `json:"CinemaName"`
	CinemaTimeZoneATE string `json:"CinemaTimeZoneATE"`
	MarketName        string `json:"MarketName"`
	CinemaSlug        string `json:"CinemaSlug"`
	MarketSlug        string `json:"MarketSlug"`
	Films             []Film `json:"Films"`
}

func (c *Cinema) GetFilmNames() []string {
	filmset := make(map[string]struct{})
	for i := 0; i < len(c.Films); i++ {
		filmset[c.Films[i].FilmName] = struct{}{}
	}

	i := 0
	films := make([]string, len(filmset))
	for k := range filmset {
		films[i] = k
		i++
	}

	return films
}
