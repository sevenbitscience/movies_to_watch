# Movies to Watch

Track movies that you want to watch in the future.

The Go server was written by me, a human, but the clients were written by
Google's Gemini.

The server was hacked together in a couple days, and is my first attempt at
creating a REST API.

This project makes use of [The Movie Database](https://www.themoviedb.org/) and
requires an API key from them.

The server provides a REST API, with endpoints documented in the
`specifications_by_gemini.md` file.

## Configuration

Settings for the server should be provided in a .env file.

```
# For this server
MOVIES_SERVER_URL=":5000"

# For TMDB API
TMDB_API_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

# For the SQlite database
MOVIES_DATABASE_PATH=./movies.db
```

