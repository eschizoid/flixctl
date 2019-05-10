package radarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Movie ...
type Movie struct {
	Added      string `json:"added"`
	AddOptions struct {
		IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
		IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
		SearchForMovie             bool `json:"searchForMovie"`
	} `json:"addOptions"`
	AlternativeTitles []struct {
		Language   string `json:"language"`
		MovieID    int    `json:"movieId"`
		SourceID   int    `json:"sourceId"`
		SourceType string `json:"sourceType"`
		Title      string `json:"title"`
		VoteCount  int    `json:"voteCount"`
		Votes      int    `json:"votes"`
	} `json:"alternativeTitles"`
	CleanTitle       string   `json:"cleanTitle"`
	Deleted          bool     `json:"deleted"`
	Downloaded       bool     `json:"downloaded"`
	ErrorMessage     string   `json:"error"`
	EpisodeCount     int      `json:"episodeCount"`
	EpisodeFileCount int      `json:"episodeFileCount"`
	FolderName       string   `json:"folderName"`
	Genres           []string `json:"genres"`
	HasFile          bool     `json:"hasFile"`
	ID               int      `json:"id"`
	ImdbID           string   `json:"imdbId"`
	Images           []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	InCinemas           string `json:"inCinemas"`
	IsAvailable         bool   `json:"isAvailable"`
	IsExisting          bool   `json:"isExisting"`
	MinimumAvailability string `json:"minimumAvailability"`
	Monitored           bool   `json:"monitored"`
	Overview            string `json:"overview"`
	Path                string `json:"path"`
	PathState           string `json:"pathState"`
	PhysicalRelease     string `json:"physicalRelease"`

	ProfileID        int `json:"profileId"`
	QualityProfileID int `json:"qualityProfileId"`
	Ratings          struct {
		Value float64 `json:"value"`
		Votes int     `json:"votes"`
	} `json:"ratings"`
	RemotePoster          string   `json:"remotePoster"`
	RootFolderPath        string   `json:"rootFolderPath"`
	Runtime               int      `json:"runtime"`
	SecondaryYearSourceID int      `json:"secondaryYearSourceId"`
	SizeOnDisk            int      `json:"sizeOnDisk"`
	SortTitle             string   `json:"sortTitle"`
	Saved                 bool     `json:"saved"`
	Status                string   `json:"status"`
	Studio                string   `json:"studio"`
	Tags                  []string `json:"tags"`
	Title                 string   `json:"title"`
	TitleSlug             string   `json:"titleSlug"`
	TmdbID                int      `json:"tmdbId"`
	Year                  int      `json:"year"`
	YouTubeTrailerID      string   `json:"youTubeTrailerId"`
	Website               string   `json:"website"`
}

// Library movies in wanted list and recognized by radarr
type Library struct {
	Page          int     `json:"page"`
	PageSize      int     `json:"pageSize"`
	SortKey       string  `json:"sortKey"`
	SortDirection string  `json:"sortDirection"`
	TotalRecords  int     `json:"totalRecords"`
	Records       []Movie `json:"records"`
}

// ErrorMovieExists error when trying to add an already added movie
var ErrorMovieExists = errors.New("This movie has already been added")

// ErrorPathAlreadyConfigured path exists for another movie
var ErrorPathAlreadyConfigured = errors.New("Path is already configured for another movie")

// AddMovie adds a movie to your wanted list
func (c Client) AddMovie(movie Movie) []error {
	const endpoint = "/api/movie"

	// check required fields
	if movie.Title == "" {
		return []error{errors.New("title is required")}
	}

	if movie.QualityProfileID == 0 {
		return []error{errors.New("quality profile id needs to be set")}
	}

	if movie.TitleSlug == "" {
		return []error{errors.New("title slug is required")}
	}

	if len(movie.Images) == 0 {
		return []error{errors.New("an array of images is required")}
	}

	if movie.TmdbID == 0 {
		return []error{errors.New("tmdbid is required")}
	}

	if movie.Path == "" && movie.RootFolderPath == "" {
		return []error{errors.New("either a path or rootFolderPath is required")}
	}

	requestPayload, err := json.Marshal(movie)

	if err != nil {
		return []error{err}
	}

	resp, err := c.post(endpoint, requestPayload)

	if err != nil {
		return []error{err}
	}

	defer resp.Body.Close()

	// return the bad request error messages
	if resp.StatusCode != http.StatusCreated && resp.StatusCode == http.StatusBadRequest {
		var errMessages []ErrorMessage
		var errs []error

		// :/ we couldn't decode the error message -- bad struct?
		if err := json.NewDecoder(resp.Body).Decode(&errMessages); err != nil {
			return []error{fmt.Errorf("unable to decode error message (bad request): %v", err)}
		}

		// turn ErrorMessage into Go error
		for _, err := range errMessages {
			var newErr error

			switch err.Message {
			case ErrorMovieExists.Error():
				newErr = ErrorMovieExists
			case ErrorPathAlreadyConfigured.Error():
				newErr = ErrorPathAlreadyConfigured
			default:
				newErr = fmt.Errorf(err.Message)
			}

			errs = append(errs, newErr)
		}

		return errs
	}

	if resp.StatusCode != http.StatusCreated {
		return []error{errors.New(resp.Status)}
	}

	return nil
}

// DeleteMovie removes a movie from your wanted list and/or local disk
// id is the id for the movie in the radarr library
func (c Client) DeleteMovie(id string, deleteFiles, addExclusion bool) error {
	const endpoint = "/api/movie/%s"

	params := make(url.Values, 1)

	params.Set("deleteFiles", strconv.FormatBool(deleteFiles))
	params.Set("addExclusion", strconv.FormatBool(addExclusion))

	resp, err := c.delete(fmt.Sprintf(endpoint, id), params)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

// DiscoverMovies returns a list of recommended movies
func (c Client) DiscoverMovies() ([]Movie, error) {
	const endpoint = "/api/movies/discover/recommendations"

	var movies []Movie

	resp, err := c.get(endpoint, nil)

	if err != nil {
		return movies, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movies, errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&movies)

	// does this endpoint always return 29 items?
	// movieCount := len(movies)

	// fmt.Printf("discover movies count: %d\n", movieCount)
	return movies, err
}
