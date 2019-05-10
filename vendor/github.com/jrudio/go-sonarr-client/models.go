package sonarr

import (
	"errors"
	"time"
)

// Series stored on the Sonarr server.
type Series struct {
	Title      string `json:"title"`
	AddOptions struct {
		IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
		IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
		SearchForMissingEpisodes   bool `json:"searchForMissingEpisodes"`
	} `json:"addOptions"`
	AlternateTitles []struct {
		Title        string `json:"title"`
		SeasonNumber int    `json:"seasonNumber"`
	} `json:"alternateTitles"`
	SortTitle         string    `json:"sortTitle"`
	SeasonCount       int       `json:"seasonCount"`
	TotalEpisodeCount int       `json:"totalEpisodeCount"`
	EpisodeCount      int       `json:"episodeCount"`
	EpisodeFileCount  int       `json:"episodeFileCount"`
	SizeOnDisk        int       `json:"sizeOnDisk"`
	Status            string    `json:"status"`
	Overview          string    `json:"overview"`
	PreviousAiring    time.Time `json:"previousAiring"`
	Network           string    `json:"network"`
	AirTime           string    `json:"airTime"`
	Images            []struct {
		CoverType string `json:"coverType"`
	} `json:"images"`
	Seasons []struct {
		SeasonNumber int  `json:"seasonNumber"`
		Monitored    bool `json:"monitored"`
		Statistics   struct {
			PreviousAiring    time.Time `json:"previousAiring"`
			EpisodeFileCount  int       `json:"episodeFileCount"`
			EpisodeCount      int       `json:"episodeCount"`
			TotalEpisodeCount int       `json:"totalEpisodeCount"`
			SizeOnDisk        int       `json:"sizeOnDisk"`
			PercentOfEpisodes int       `json:"percentOfEpisodes"`
		} `json:"statistics"`
	} `json:"seasons"`
	Year              int       `json:"year"`
	Path              string    `json:"path"`
	ProfileID         int       `json:"profileId"`
	SeasonFolder      bool      `json:"seasonFolder"`
	Monitored         bool      `json:"monitored"`
	UseSceneNumbering bool      `json:"useSceneNumbering"`
	Runtime           int       `json:"runtime"`
	TvdbID            int       `json:"tvdbId"`
	TvRageID          int       `json:"tvRageId"`
	TvMazeID          int       `json:"tvMazeId"`
	FirstAired        time.Time `json:"firstAired"`
	LastInfoSync      time.Time `json:"lastInfoSync"`
	SeriesType        string    `json:"seriesType"`
	CleanTitle        string    `json:"cleanTitle"`
	ImdbID            string    `json:"imdbId"`
	TitleSlug         string    `json:"titleSlug"`
	Certification     string    `json:"certification"`
	Genres            []string  `json:"genres"`
	Tags              []int     `json:"tags"`
	Added             time.Time `json:"added"`
	Ratings           struct {
		Votes int     `json:"votes"`
		Value float32 `json:"value"`
	} `json:"ratings"`
	RootFolderPath   string `json:"rootFolderPath"`
	QualityProfileID int    `json:"qualityProfileId"`
	ID               int    `json:"id"`
}

// Episode of a Series.
type Episode struct {
	SeriesID                 int       `json:"seriesId"`
	EpisodeFileID            int       `json:"episodeFileID"`
	SeasonNumber             int       `json:"seasonNumber"`
	EpisodeNumber            int       `json:"episodeNumber"`
	Title                    string    `json:"title"`
	AirDate                  string    `json:"airDate"`
	AirDateUTC               time.Time `json:"airDateUTC"`
	Overview                 string    `json:"overview"`
	HasFile                  bool      `json:"hasFile"`
	Monitored                bool      `json:"monitored"`
	UnverifiedSceneNumbering bool      `json:"unverifiedSceneNumbering"`
	ID                       int       `json:"id"`
}

// Quality of a file.
type Quality struct {
	Quality struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"quality"`
	Revision struct {
		Version int `json:"version"`
		Real    int `json:"real"`
	} `json:"revision"`
	Proper bool `json:"proper"`
}

// EpisodeFile of an Episode. Represents a file stored on disk.
type EpisodeFile struct {
	SeriesID            int     `json:"seriesId"`
	SeasonNumber        int     `json:"seasonNumber"`
	RelativePath        string  `json:"relativePath"`
	Path                string  `json:"path"`
	Size                int     `json:"size"`
	DateAdded           string  `json:"dateAdded"`
	SceneName           string  `json:"sceneName"`
	Quality             Quality `json:"quality"`
	QualityCutoffNotMet bool    `json:"qualityCutoffNotMet"`
	ID                  int     `json:"id"`
}

// Queue item currently being downloaded.
type Queue struct {
	Series                Series  `json:"series"`
	Episode               Episode `json:"episode"`
	Quality               Quality `json:"quality"`
	Size                  int     `json:"size"`
	Title                 string  `json:"title"`
	SizeLeft              int     `json:"sizeLeft"`
	Status                string  `json:"status"`
	TrackedDownloadStatus string  `json:"trackedDownloadStatus"`
	StatusMessages        []struct {
		Title    string   `json:"title"`
		Messages []string `json:"messages"`
	} `json:"statusMessages"`
	DownloadID string `json:"downloadId"`
	Protocol   string `json:"protocol"`
	ID         int    `json:"id"`
}

// Calendar entry for a past or upcoming airing.
type Calendar struct {
	SeriesID                 int       `json:"seriesId"`
	EpisodeFileID            int       `json:"episodeFileId"`
	SeasonNumber             int       `json:"seasonNumber"`
	EpisodeNumber            int       `json:"episodeNumber"`
	Title                    string    `json:"title"`
	AirDate                  string    `json:"airDate"`
	AirDateUTC               time.Time `json:"airDateUtc"`
	HasFile                  bool      `json:"hasFile"`
	Monitored                bool      `json:"monitored"`
	AbsoluteEpisodeNumber    int       `json:"absoluteEpisodeNumber"`
	Series                   Series    `json:"series"`
	UnverifiedSceneNumbering bool      `json:"unverifiedSceneNumbering"`
}

// DiskSpace remaining on each drive mounted on the server.
type DiskSpace struct {
	Path       string `json:"path"`
	Label      string `json:"label"`
	FreeSpace  int    `json:"freeSpace"`
	TotalSpace int    `json:"totalSpace"`
}

// Tag used to tag Series.
type Tag struct {
	Label string `json:"label"`
	ID    int    `json:"id"`
}

// SystemStatus of the server.
type SystemStatus struct {
	Version           string `json:"version"`
	BuildTime         string `json:"buildTime"`
	IsDebug           bool   `json:"isDebug"`
	IsProduction      bool   `json:"isProduction"`
	IsAdmin           bool   `json:"isAdmin"`
	IsUserInteractive bool   `json:"isUserInteractive"`
	StartupTime       string `json:"startupTime"`
	AppData           string `json:"appData"`
	OsName            string `json:"osName"`
	OsVersion         string `json:"osVersion"`
	IsMonoRuntime     bool   `json:"isMonoRuntime"`
	IsMono            bool   `json:"isMono"`
	IsLinux           bool   `json:"isLinux"`
	IsOsx             bool   `json:"isOsx"`
	IsWindows         bool   `json:"isWindows"`
	Branch            string `json:"branch"`
	Authentication    string `json:"forms"`
	SqlliteVersion    string `json:"sqlliteVersion"`
	URLBase           string `json:"urlBase"`
	RuntimeVersion    string `json:"runtimeVersion"`
	RuntimeName       string `json:"runtimeName"`
}

// SearchResults are results of Sonarr's online series lookup
// not the results of what's already located on the Sonarr server
type SearchResults struct {
	Added         string   `json:"added"`
	AirTime       string   `json:"airTime"`
	Certification string   `json:"certification"`
	CleanTitle    string   `json:"cleanTitle"`
	FirstAired    string   `json:"firstAired"`
	Genres        []string `json:"genres"`
	Images        []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	ImdbID           string `json:"imdbId"`
	Monitored        bool   `json:"monitored"`
	Network          string `json:"network"`
	Overview         string `json:"overview"`
	ProfileID        int    `json:"profileId"`
	QualityProfileID int    `json:"qualityProfileId"`
	Ratings          struct {
		Value float64 `json:"value"`
		Votes int     `json:"votes"`
	} `json:"ratings"`
	RemotePoster string `json:"remotePoster"`
	Runtime      int    `json:"runtime"`
	SeasonCount  int    `json:"seasonCount"`
	SeasonFolder bool   `json:"seasonFolder"`
	Seasons      []struct {
		Monitored    bool `json:"monitored"`
		SeasonNumber int  `json:"seasonNumber"`
	} `json:"seasons"`
	SeriesType        string   `json:"seriesType"`
	SortTitle         string   `json:"sortTitle"`
	Status            string   `json:"status"`
	Tags              []string `json:"tags"`
	Title             string   `json:"title"`
	TitleSlug         string   `json:"titleSlug"`
	TvMazeID          int      `json:"tvMazeId"`
	TvRageID          int      `json:"tvRageId"`
	TvdbID            int      `json:"tvdbId"`
	UseSceneNumbering bool     `json:"useSceneNumbering"`
	Year              int      `json:"year"`
}

// ErrorMessage radarr's error message struct
type ErrorMessage struct {
	PropertyName                      string      `json:"propertyName"`
	Message                           string      `json:"errorMessage"`
	AttemptedValue                    interface{} `json:"attemptedValue"`
	FormattedMessageArguments         []string    `json:"formattedMessageArguments"`
	FormattedMessagePlaceholderValues struct {
		PropertyName  string      `json:"propertyName"`
		PropertyValue interface{} `json:"propertyValue"`
	} `json:"formattedMessagePlaceholderValues"`
}

// ErrorSeriesExists error when trying to add an already added series
var ErrorSeriesExists = errors.New("This series has already been added")

// ErrorPathAlreadyConfigured path exists for another series
var ErrorPathAlreadyConfigured = errors.New("Path is already configured for another series")
