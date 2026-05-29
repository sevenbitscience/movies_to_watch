package types

// Represent a movie in my database
type Movie struct {
	ID 		int 		`json:"id"`
	Title	string		`json:"title"`
	TMDB_ID	int			`json:"tmdb_id"`
	Year	*int		`json:"release_year"`
	Genres	[]string	`json:"genre"`
	Status	string		`json:"status"`
}

