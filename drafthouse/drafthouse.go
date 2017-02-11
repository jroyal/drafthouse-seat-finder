package drafthouse

import "time"

func GetMovies(dayFilter time.Time) ResponseMoviesShowingToday {
	market := getMarketInfo()
	return ResponseMoviesShowingToday{market.Movies(dayFilter)}
}
