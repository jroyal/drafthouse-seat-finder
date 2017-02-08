package drafthouse

import "github.com/jroyal/drafthouse-seat-finder/utils"

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
	filmset := utils.NewStringSet()
	for i := 0; i < len(c.Films); i++ {
		filmset.Add(c.Films[i].FilmName)
	}
	return filmset.ToSlice()
}
