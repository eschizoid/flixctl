package radarr

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// utils.go holds network utils and function helpers

func (c Client) get(query string, params url.Values) (*http.Response, error) {
	endpointURL, err := url.Parse(query)

	if err != nil {
		return &http.Response{}, err
	}

	if params == nil {
		params = endpointURL.Query()
	}

	endpointURL.RawQuery = params.Encode()

	client := http.Client{
		Timeout: time.Duration(c.Timeout) * time.Second,
	}

	requestURL := appendEndpoint(c.URL.String(), endpointURL.String())

	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func (c Client) post(query string, body []byte) (*http.Response, error) {
	endpointURL, err := url.Parse(query)

	if err != nil {
		return &http.Response{}, err
	}

	client := http.Client{
		Timeout: time.Duration(c.Timeout) * time.Second,
	}

	requestURL := appendEndpoint(c.URL.String(), endpointURL.String())

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(body))

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func (c Client) delete(query string, params url.Values) (*http.Response, error) {
	endpointURL, err := url.Parse(query)

	if err != nil {
		return &http.Response{}, err
	}

	if params == nil {
		params = endpointURL.Query()
	}

	endpointURL.RawQuery = params.Encode()

	client := http.Client{
		Timeout: time.Duration(c.Timeout) * time.Second,
	}

	requestURL := appendEndpoint(c.URL.String(), endpointURL.String())

	req, err := http.NewRequest("DELETE", requestURL, nil)

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func encodeURL(str string) (string, error) {
	u, err := url.Parse(str)

	if err != nil {
		return "", err
	}

	return u.String(), nil
}

// appendEndpoint checks for and applies a "/" to a url when necessary
func appendEndpoint(baseURL, endpoint string) string {
	// /api/series && http://192.168.1.25:8989/
	if strings.HasPrefix(endpoint, "/") && strings.HasSuffix(baseURL, "/") {
		endpoint = endpoint[1:]
	} else if !strings.HasPrefix(endpoint, "/") && !strings.HasSuffix(baseURL, "/") {
		// api/series && http://192.168.1.25:8989
		endpoint = "/" + endpoint
	}

	return baseURL + endpoint
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
