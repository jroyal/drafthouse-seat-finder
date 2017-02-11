package drafthouse

import "time"

func GetMovies(dayFilter time.Time, cinemaFilter string) ResponseMovies {
	market := getMarketInfo()
	return ResponseMovies{market.Movies(dayFilter, cinemaFilter)}
}
