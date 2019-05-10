package sonarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Sonarr contains fields needed to make API calls to a Sonarr server
type Sonarr struct {
	baseURL    *url.URL
	apiKey     string
	HTTPClient http.Client
	// Timeout in seconds -- default 5
	Timeout int
}

const (
	calendarEndpoint     = "calendar"
	diskSpaceEndpoint    = "diskspace"
	episodeEndpoint      = "episode"
	episodeFileEndpoint  = "episodefile"
	seriesEndpoint       = "series"
	systemStatusEndpoint = "system/status"
	tagEndpoint          = "tag"
)

// New creates a new Sonarr client instance.
func New(apiURL, apiKey string) (*Sonarr, error) {
	if apiURL == "" {
		return &Sonarr{}, errors.New("apiURL is required")
	}

	if apiKey == "" {
		return &Sonarr{}, errors.New("apiKey is required")
	}

	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}

	baseURL, err := url.Parse(apiURL)

	if err != nil {
		return &Sonarr{}, fmt.Errorf("Failed to parse baseURL: %v", err)
	}

	return &Sonarr{
		baseURL:    baseURL,
		apiKey:     apiKey,
		HTTPClient: http.Client{},
	}, nil
}

// GetCalendar retrieves info about when episodes were/will be downloaded.
// If start and end are not provided, retrieves episodes airing today and tomorrow.
func (s *Sonarr) GetCalendar(start, end string) ([]Calendar, error) {
	params := make(url.Values)
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}
	var results []Calendar
	res, err := s.get(calendarEndpoint, params)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&results)
	return results, err
}

// GetDiskSpace retrieves info about the disk space remaining on the server.
func (s *Sonarr) GetDiskSpace() ([]DiskSpace, error) {
	var results []DiskSpace
	res, err := s.get(diskSpaceEndpoint, nil)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&results)
	return results, err
}

// GetEpisodes retrieves all Episodes for the given series ID.
func (s *Sonarr) GetEpisodes(seriesID int) ([]Episode, error) {
	var results []Episode
	if seriesID <= 0 {
		return results, errors.New("seriesID must be a positive integer")
	}
	params := make(url.Values)
	params.Set("seriesId", strconv.Itoa(seriesID))
	res, err := s.get(episodeEndpoint, params)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&results)
	return results, err
}

// GetEpisode retrieves the Episode with the given ID.
func (s *Sonarr) GetEpisode(episodeID int) (*Episode, error) {
	results := &Episode{}
	if episodeID <= 0 {
		return results, errors.New("episodeID must be a positive integer")
	}
	episodeURL := fmt.Sprintf("%s/%s", episodeEndpoint, strconv.Itoa(episodeID))
	res, err := s.get(episodeURL, nil)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(results)
	return results, err
}

// UpdateEpisode updates the given Episode. Currently, the API only supports
// updating the "Monitored" status. Any other changes are ignored.
// This should be an Episode you have previously retrieved with GetEpisodes()
// or GetEpisode(). The updated Episode is returned.
func (s *Sonarr) UpdateEpisode(ep *Episode) (*Episode, error) {
	results := &Episode{}
	episodeURL := fmt.Sprintf("%s/%s", episodeEndpoint, strconv.Itoa(ep.ID))
	res, err := s.put(episodeURL, ep)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(results)
	return results, err
}

// GetEpisodeFiles retrieves all EpisodeFiles for the given series ID.
func (s *Sonarr) GetEpisodeFiles(seriesID int) ([]EpisodeFile, error) {
	var results []EpisodeFile
	if seriesID <= 0 {
		return results, errors.New("seriesID must be a positive integer")
	}
	params := make(url.Values)
	params.Set("seriesId", strconv.Itoa(seriesID))
	res, err := s.get(episodeFileEndpoint, params)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&results)
	return results, err
}

// GetEpisodeFile retrieves the EpisodeFile with the given ID.
func (s *Sonarr) GetEpisodeFile(episodeFileID int) (*EpisodeFile, error) {
	results := &EpisodeFile{}
	if episodeFileID <= 0 {
		return results, errors.New("episodeFileID must be a positive integer")
	}
	episodeFileURL := fmt.Sprintf("%s/%s", episodeFileEndpoint, strconv.Itoa(episodeFileID))
	res, err := s.get(episodeFileURL, nil)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(results)
	return results, err
}

// DeleteEpisodeFile deletes the EpisodeFile with the given ID.
// This also deletes the media file from disk!
func (s *Sonarr) DeleteEpisodeFile(episodeFileID int) (*EpisodeFile, error) {
	results := &EpisodeFile{}
	if episodeFileID <= 0 {
		return results, errors.New("episodeFileID must be a positive integer")
	}
	episodeFileURL := fmt.Sprintf("%s/%s", episodeFileEndpoint, strconv.Itoa(episodeFileID))
	res, err := s.del(episodeFileURL, nil)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(results)
	return results, err
}

// GetAllSeries retrieves all Series.
func (s *Sonarr) GetAllSeries() ([]Series, error) {
	var results []Series
	res, err := s.get(seriesEndpoint, nil)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return results, errors.New(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&results)
	return results, err
}

// GetSeries retrieves the Series with the given ID.
func (s *Sonarr) GetSeries(seriesID int) (*Series, error) {
	results := &Series{}
	if seriesID <= 0 {
		return results, errors.New("seriesID must be a positive integer")
	}
	seriesURL := fmt.Sprintf("%s/%s", seriesEndpoint, strconv.Itoa(seriesID))
	res, err := s.get(seriesURL, nil)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return results, errors.New(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(results)
	return results, err
}

// GetSeriesFromTVDB retrieves the Series with the given ID.
func (s *Sonarr) GetSeriesFromTVDB(seriesID int) (*Series, error) {
	const endpoint = "/api/series/lookup"

	params := url.Values{}

	params.Set("term", fmt.Sprintf("tvdb:%d", seriesID))

	matchedSeries := &Series{}

	if seriesID <= 0 {
		return matchedSeries, errors.New("seriesID must be a positive integer")
	}

	res, err := s.get(endpoint, params)

	if err != nil {
		return matchedSeries, err
	}

	defer res.Body.Close()

	// handle non-200 status code
	if res.StatusCode != http.StatusOK {
		return matchedSeries, errors.New(res.Status)
	}

	var results []Series

	err = json.NewDecoder(res.Body).Decode(&results)

	for _, show := range results {
		if show.TvdbID == seriesID {
			matchedSeries = &show
		}
	}

	if matchedSeries.TvdbID != seriesID {
		return matchedSeries, errors.New("invalid series id: " + strconv.Itoa(seriesID))
	}

	return matchedSeries, err
}

// UpdateSeries updates the given Series.
// This should be a Series you have previously retrieved with GetAllSeries()
// or GetSeries(). The updated Series is returned.
func (s *Sonarr) UpdateSeries(ser *Series) (*Series, error) {
	results := &Series{}
	seriesURL := fmt.Sprintf("%s/%s", seriesEndpoint, strconv.Itoa(ser.ID))
	res, err := s.put(seriesURL, ser)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(results)
	return results, err
}

// AddSeries adds a movie to your wanted list
func (s Sonarr) AddSeries(series Series) []error {
	const endpoint = "/api/series"

	// check required fields
	if series.Title == "" {
		return []error{errors.New("title is required")}
	}

	if series.QualityProfileID == 0 {
		return []error{errors.New("quality profile id needs to be set")}
	}

	if series.TitleSlug == "" {
		return []error{errors.New("title slug is required")}
	}

	if len(series.Images) == 0 {
		return []error{errors.New("an array of images is required")}
	}

	if series.TvdbID == 0 {
		return []error{errors.New("tvdbid is required")}
	}

	if series.Path == "" && series.RootFolderPath == "" {
		return []error{errors.New("either a path or rootFolderPath is required")}
	}

	requestPayload, err := json.Marshal(series)

	if err != nil {
		return []error{err}
	}

	resp, err := s.post(endpoint, requestPayload)

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
			case ErrorSeriesExists.Error():
				newErr = ErrorSeriesExists
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

// DeleteSeries deletes the Series with the given ID.
// If deleteFiles is true, the series folder and all files will be deleted too.
func (s *Sonarr) DeleteSeries(seriesID int, deleteFiles bool) (*Series, error) {
	results := &Series{}
	if seriesID <= 0 {
		return results, errors.New("seriesID must be a positive integer")
	}
	params := make(url.Values)
	if deleteFiles {
		params.Set("deleteFiles", "true")
	}
	seriesURL := fmt.Sprintf("%s/%s", seriesEndpoint, strconv.Itoa(seriesID))
	res, err := s.del(seriesURL, params)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(results)
	return results, err
}

// GetSystemStatus retrieves system information about the Sonarr server.
func (s *Sonarr) GetSystemStatus() (*SystemStatus, error) {
	results := &SystemStatus{}
	res, err := s.get(systemStatusEndpoint, nil)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&results)
	return results, err
}

// GetTags retrieves all Tags that have been applied to any series.
func (s *Sonarr) GetTags() ([]Tag, error) {
	var results []Tag
	res, err := s.get(tagEndpoint, nil)
	if err != nil {
		return results, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&results)
	return results, err
}

// Search searches for media via tvdb (or Sonarr's default search engine)
func (s *Sonarr) Search(title string) ([]SearchResults, error) {
	// const searchOnlineEndpoint = "/api/series/lookup?term="
	const searchOnlineEndpoint = "/api/series/lookup"

	var results []SearchResults

	resp, err := s.get(searchOnlineEndpoint, url.Values{
		"term": []string{title},
	})

	if err != nil {
		return results, err
	}

	defer resp.Body.Close()

	// sonarr returns a json object with a field `error` for an invalid api key
	if resp.StatusCode != http.StatusOK {
		return results, errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&results)

	return results, err
}
