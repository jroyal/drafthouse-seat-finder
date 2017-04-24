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

type SimpleCinema struct {
	CinemaName string
	CinemaSlug string
}

func (c *Cinema) GetFilmNames() []string {
	filmset := utils.NewStringSet()
	for i := 0; i < len(c.Films); i++ {
		filmset.Add(c.Films[i].FilmName)
	}
	return filmset.ToSlice()
}

func (c *Cinema) GetFilms() []SimpleFilm {
	filmset := make([]SimpleFilm, len(c.Films))
	for i, film := range c.Films {
		filmset[i] = SimpleFilm{
			FilmName: film.FilmName,
			FilmSlug: film.FilmSlug,
			FilmYear: film.FilmYear,
			FilmID:   film.FilmID,
		}
	}
	return filmset
}

func (c *Cinema) GetFilmTimes(filmSlug string) []string {
	var filmTimes []string
	for _, film := range c.Films {
		if film.FilmSlug == filmSlug {
			filmTimes = film.GetFilmTimes()
			break
		}
	}
	return filmTimes
}

func (c *Cinema) GetFilmSessions(filmSlug string) []FilmSession {
	var filmSessions []FilmSession
	for _, film := range c.Films {
		if film.FilmSlug == filmSlug {
			filmSessions = film.GetFilmSessions()
			break
		}
	}

	for i := 0; i < len(filmSessions); i++ {
		film := &filmSessions[i]
		film.CinemaID = c.CinemaID
		film.CinemaName = c.CinemaName
	}
	return filmSessions
}
