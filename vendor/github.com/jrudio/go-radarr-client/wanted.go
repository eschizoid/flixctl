package radarr

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Wanted ...
type Wanted struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Records  []struct {
		Added             string `json:"added"`
		AlternativeTitles []struct {
			ID         int    `json:"id"`
			Language   string `json:"language"`
			MovieID    int    `json:"movieId"`
			SourceID   int    `json:"sourceId"`
			SourceType string `json:"sourceType"`
			Title      string `json:"title"`
			VoteCount  int    `json:"voteCount"`
			Votes      int    `json:"votes"`
		} `json:"alternativeTitles"`
		CleanTitle string        `json:"cleanTitle"`
		Downloaded bool          `json:"downloaded"`
		FolderName string        `json:"folderName"`
		Genres     []interface{} `json:"genres"`
		HasFile    bool          `json:"hasFile"`
		ID         int           `json:"id"`
		Images     []struct {
			CoverType string `json:"coverType"`
			URL       string `json:"url"`
		} `json:"images"`
		ImdbID              string `json:"imdbId"`
		InCinemas           string `json:"inCinemas"`
		IsAvailable         bool   `json:"isAvailable"`
		LastInfoSync        string `json:"lastInfoSync"`
		MinimumAvailability string `json:"minimumAvailability"`
		Monitored           bool   `json:"monitored"`
		Overview            string `json:"overview"`
		Path                string `json:"path"`
		PathState           string `json:"pathState"`
		PhysicalRelease     string `json:"physicalRelease"`
		ProfileID           int    `json:"profileId"`
		QualityProfileID    int    `json:"qualityProfileId"`
		Ratings             struct {
			Value float64 `json:"value"`
			Votes int     `json:"votes"`
		} `json:"ratings"`
		Runtime               int           `json:"runtime"`
		SecondaryYearSourceID int           `json:"secondaryYearSourceId"`
		SizeOnDisk            int           `json:"sizeOnDisk"`
		SortTitle             string        `json:"sortTitle"`
		Status                string        `json:"status"`
		Studio                string        `json:"studio"`
		Tags                  []interface{} `json:"tags"`
		Title                 string        `json:"title"`
		TitleSlug             string        `json:"titleSlug"`
		TmdbID                int           `json:"tmdbId"`
		Website               string        `json:"website"`
		Year                  int           `json:"year"`
		YouTubeTrailerID      string        `json:"youTubeTrailerId"`
	} `json:"records"`
	SortDirection string `json:"sortDirection"`
	SortKey       string `json:"sortKey"`
	TotalRecords  int    `json:"totalRecords"`
}

// GetWantedMissing returns a filtered wanted list of missing media
func (c Client) GetWantedMissing() (Wanted, error) {
	const endpoint = "/api/wanted/missing"

	params := make(url.Values, 1)

	params.Set("page", "1")
	params.Set("pageSize", "50")
	params.Set("sortKey", "title")
	params.Set("sortDir", "asc")
	params.Set("filterKey", "monitored")
	params.Set("filterValue", "true")
	params.Set("filterType", "equal")

	var wanted Wanted

	resp, err := c.get(endpoint, nil)

	if err != nil {
		return wanted, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return wanted, errors.New(err.Error())
	}

	err = json.NewDecoder(resp.Body).Decode(&wanted)

	return wanted, err
}
