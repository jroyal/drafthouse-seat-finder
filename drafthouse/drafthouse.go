package drafthouse

func GetMoviesShowingToday() ResponseMoviesShowingToday {
	market := getMarketInfo()
	return ResponseMoviesShowingToday{market.MoviesShowingToday()}
}
