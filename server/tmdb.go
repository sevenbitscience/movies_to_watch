package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/*
   "genre_ids": [
     28
   ],
   "id": 347807,
   "original_language": "hi",
   "original_title": "Fight Club: Members Only",
   "overview": "Four friends head off to Bombay and get involved in the mother and father of all gang wars.",
   "popularity": 2.26,
   "poster_path": "/aXFmWfWYCCxQTkCn7K86RvDiMHZ.jpg",
   "release_date": "2006-02-17",
   "title": "Fight Club: Members Only",
   "video": false,
   "vote_average": 4.5,
   "vote_count": 12
*/


type TmdbMovie struct {
	Adult				bool	`json:"adult"`
	Backdrop_path		string	`json:"backdrop_path"`
	Genre_ids			[]int	`json:"genre_ids"`
	Id					int		`json:"id"`
	Original_language	string	`json:"original_language"`
	Original_title		string	`json:"original_title"`
	Overview			string	`json:"overview"`
	Popularity			float64	`json:"popularity"`
	Poster_path			string	`json:"poster_path"`
	Release_date		string	`json:"release_date"`
	Title				string	`json:"title"`
	Video				bool	`json:"video"`
	Vote_average		float64	`json:"vote_average"`
	Vote_count			int		`json:"vote_count"`
}

type Genre struct {
	Id		int		`json:"id"`
	Name	string	`json:"name"`
}

var GenreTable map[int]string

func tmdbEntryToMovie(m *TmdbMovie) Movie {
	// Parse out genres and convert them to strings
	var genreStrings []string
	for _, v := range(m.Genre_ids) {
		genreStrings = append(genreStrings, lookupGenre(v))
	}

	// Create a new Movie object
	newMovie := Movie{
		Title: m.Title,
		TMDB_ID: m.Id,
		Genres: genreStrings,
	}

	// Parse the year from YYYY-MM-DD format
	splitDate := strings.Split(m.Release_date, "-")
	if len(splitDate) > 0 && splitDate[0] != "" {
		justYear, err := strconv.Atoi(splitDate[0])
		if err != nil {
			log.Printf("Couldn't parse year %s", splitDate[0])
		}
		if justYear != 0 {
			newMovie.Year = new(int)
			*newMovie.Year = justYear
		}
	}

	return newMovie
}

var apiKey string

func setApiKey(key string) {
	apiKey = "Bearer " + key
}

func readApi(url string) ([]byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Couldn't build the request to %s", url)
		return nil, err
	}

	// Add needed headers including API key
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", apiKey)

	log.Printf("Quering TMDB @ %s", url)
	// Actually make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Couldn't connect to %s, %s", url, err.Error())
		return nil, err
	}
	defer res.Body.Close()
	
	// Make sure the response was okay
	if res.StatusCode != http.StatusOK {
		log.Printf("Bad response from TMDB: %d %s", res.StatusCode, res.Status)
	}

	// Read the reply and pass it along
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to parse response from %s", url)
		return nil, err
	}

	return body, nil
}

// https://api.themoviedb.org/3/search/movie?query=this%20is%20my%20title&include_adult=true&language=en-US&page=1
 
// Get the top n results from tmdb
func findMovies(title string, movies *[]Movie) {
	baseURL, err := url.Parse("https://api.themoviedb.org/3/search/movie")
	if err != nil {
		log.Println("Error parsing URL in findMovies()")
		return 
	}

	// Set up the parameters for the search
	params := url.Values{}
	params.Add("query", title)
	params.Add("include_adult", "false")
	params.Add("page", "1")
	params.Add("language", "en-US")

	// Add that to the URL and process to a string
	baseURL.RawQuery = params.Encode()
	searchURL := baseURL.String()

	// Call the API with that query
	response, err := readApi(searchURL)
	if err != nil {
		return
	}

	// Parse that response JSON into an array of movie results
	var movieResults struct {
		Page			int			`json:"page"`
		Results			[]TmdbMovie	`json:"results"`
		TotalPages		int			`json:"total_pages"`
		TotalResults	int			`json:"total_results"`
	}

	err = json.Unmarshal(response, &movieResults)
	if err != nil {
		log.Println("Unable to extract list of movies")
		return
	}

	// Convert those to the simpler Movie struct
	for _, v := range(movieResults.Results) {
		*movies = append(*movies, tmdbEntryToMovie(&v))
	}
}


// Get the information about a movie by ID
func getMovie(tmdb_id int) Movie {
	var movie Movie
	return movie
}

// Get the name associated with a genre ID
// Should check a table of genres in the DB then ask TMDB API
func lookupGenre(g int) string {
	if GenreTable == nil {
		// GenreTable doesn't exist yet
		GenreTable = make(map[int]string)
	}

	// If the requested genre's name is in the table, return it
	name, exists := GenreTable[g]
	if exists == true {
		return name
	}

	// We need to ask the database for the table
	var movieGenreList struct {
		Genres []struct {
			Id		int		`json:"id"`
			Name	string	`json:"name"`
		}					`json:"genres"`
	}

	// Query TMDB for the list of genres and parse out the JSON
	response, err := readApi("https://api.themoviedb.org/3/genre/movie/list?language=en")
	if err != nil {
		log.Printf("Failed to read genres from API")
		return ""
	}

	err = json.Unmarshal(response, &movieGenreList)
	if err != nil {
		log.Println("Failed to parse response while fetching genres")
	}

	// Now enter the genres into the map
	for _, val := range(movieGenreList.Genres) {
		GenreTable[val.Id] = val.Name
	}

	// Now return the genre if we found it
	name, exists = GenreTable[g]
	if exists != true {
		return ""
	}
	return name
}
