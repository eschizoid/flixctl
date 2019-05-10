package sonarr

import (
	"encoding/json"
	"errors"
	"net/http"
)

// RootFolder ...
type RootFolder struct {
	FreeSpace       int    `json:"freeSpace"`
	ID              int    `json:"id"`
	Path            string `json:"path"`
	UnmappedFolders []struct {
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"unmappedFolders"`
}

// Profile ...
type Profile struct {
	Cutoff struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"cutoff"`
	ID    int `json:"id"`
	Items []struct {
		Allowed bool `json:"allowed"`
		Quality struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"quality"`
	} `json:"items"`
	Language      string `json:"language"`
	Name          string `json:"name"`
	PreferredTags string `json:"preferredTags"`
}

// GetRootFolders returns available root folders
func (s Sonarr) GetRootFolders() ([]RootFolder, error) {
	const endpoint = "/api/rootfolder"
	var folders []RootFolder

	resp, err := s.get(endpoint, nil)

	if err != nil {
		return folders, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return folders, errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&folders)

	return folders, err
}

// GetProfiles returns all movie quality settings
func (s Sonarr) GetProfiles() ([]Profile, error) {
	const endpoint = "/api/profile"
	var profiles []Profile

	resp, err := s.get(endpoint, nil)

	if err != nil {
		return profiles, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return profiles, errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&profiles)

	return profiles, err
}
