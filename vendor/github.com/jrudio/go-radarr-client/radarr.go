package radarr

import (
	"errors"
	"net/url"
	"strings"
)

// Client ...
type Client struct {
	URL    *url.URL
	APIKey string
	// Timeout in seconds -- default 5
	Timeout int
}

// New creates a client to make api calls to Radarr
func New(host, apiKey string) (Client, error) {
	var client Client

	if host == "" {
		return client, errors.New("radarr url is required")
	}

	if apiKey == "" {
		return client, errors.New("radarr api key is required")
	}

	// remove the trailing slash
	if strings.HasSuffix(host, "/") {
		host = host[0 : len(host)-1]
	}

	hostURL, err := url.Parse(host)

	if err != nil {
		return client, err
	}

	client.URL = hostURL
	client.APIKey = apiKey
	client.Timeout = 5

	return client, nil
}
