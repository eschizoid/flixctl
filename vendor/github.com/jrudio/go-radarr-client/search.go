package radarr

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// GetMovieOptions change the params when using GetMovies
type GetMovieOptions struct {
	// monitored only
	// - filterKey = monitored
	// - filterValue = true
	// - filterType = equal
	// missing only
	// - filterKey = downloaded
	// - filterValue = false
	// - filterType = equal
	// released only
	// - filterKey = downloaded
	// - filterValue = false
	// - filterType = equal
	// announced only
	// - filterKey = status
	// - filterValue = inCinemas
	// - filterType = equal
	//
	// FilterKey can be 'monitored', 'downloaded', or 'status'
	FilterKey string
	// FilterValue can be 'false', 'true', or 'inCinemas' depending on FilterKey
	FilterValue string
	// FilterType can be 'equal'
	FilterType string
	// Page when pagination is on we can skip to a point in our results
	Page string
	// PageSize can be any number or '-1' to return everything without pagination
	PageSize string
	// SortKey can be 'sortTitle'
	SortKey string
	// SortDir can be 'asc' or 'desc'
	SortDir string
}

// Search uses Radarr's method of online movie lookup
func (c Client) Search(title string) ([]Movie, error) {
	params := url.Values{}

	params.Set("term", title)

	resp, err := c.get("/api/movie/lookup", params)

	if err != nil {
		return []Movie{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Movie{}, errors.New(resp.Status)
	}

	var results []Movie

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return results, err
	}

	return results, nil
}

// SearchOffline searches for movies already in Radarr's library
func (c Client) SearchOffline(title string) (Movie, error) {
	return Movie{}, nil
}

// GetMovie returns a movie via the movie database id
func (c Client) GetMovie(tmdbID int) (Movie, error) {
	const endpoint = "/api/movie/lookup/tmdb"

	params := url.Values{}

	params.Set("tmdbId", strconv.Itoa(tmdbID))

	matchedMovie := Movie{}

	resp, err := c.get(endpoint, params)

	if err != nil {
		return matchedMovie, err
	}

	defer resp.Body.Close()

	// handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return matchedMovie, errors.New(resp.Status)
	}

	var result Movie

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return matchedMovie, err
	}

	return result, nil
}

// GetMovieIMDB returns a movie via the internet movie database id
func (c Client) GetMovieIMDB(imdbID int) (Movie, error) {
	const endpoint = "/api/movie/lookup/imdb"

	params := url.Values{}

	params.Set("imdbId", strconv.Itoa(imdbID))

	var result Movie

	resp, err := c.get(endpoint, params)

	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}

// GetMovies returns all movies in radarr and in the wanted list
func (c Client) GetMovies(options GetMovieOptions) ([]Movie, error) {
	var movies []Movie

	const endpoint = "/api/movie"

	params := url.Values{}

	if options.Page == "" {
		options.Page = "1"
	}

	// return everything in results if not specified
	if options.PageSize == "" {
		options.Page = "-1"
	}

	if options.SortKey == "" {
		options.Page = "sortTitle"
	}

	if options.SortDir == "" {
		options.Page = "asc"
	}

	params.Set("page", options.Page)
	params.Set("pageSize", options.PageSize)
	params.Set("sortKey", options.SortKey)
	params.Set("sortDir", options.SortDir)
	params.Set("filterType", options.FilterType)

	if options.FilterKey != "" {
		params.Set("filterKey", options.FilterKey)
	}

	if options.FilterValue != "" {
		params.Set("filterValue", options.FilterValue)
	}

	if options.FilterType != "" {
		params.Set("filterType", options.FilterType)
	}

	resp, err := c.get(endpoint, params)

	if err != nil {
		return movies, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movies, errors.New(resp.Status)
	}

	if options.PageSize == "-1" {
		if err := json.NewDecoder(resp.Body).Decode(&movies); err != nil {
			return movies, err
		}
	} else {
		var library Library

		if err := json.NewDecoder(resp.Body).Decode(&library); err != nil {
			return movies, err
		}

		movies = library.Records
	}

	return movies, nil
}
