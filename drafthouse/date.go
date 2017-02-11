package drafthouse

import (
	"time"

	"github.com/jroyal/drafthouse-seat-finder/utils"
	log "github.com/sirupsen/logrus"
)

type Date struct {
	DateID  string   `json:"DateId"`
	Date    string   `json:"Date"`
	Cinemas []Cinema `json:"Cinemas"`
}

func (d *Date) convertToTime() time.Time {
	date, err := time.Parse(dateFormat, d.Date)
	if err != nil {
		log.WithField("err", err).Error("Failed to properly convert date")
	}
	return date
}

func (d *Date) getMovies() []string {
	movies := utils.NewStringSet()
	for _, cinema := range d.Cinemas {
		movies.AddSlice(cinema.GetFilmNames())
	}
	return movies.ToSlice()
}
