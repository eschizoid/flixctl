package sonarr

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (s *Sonarr) get(endpoint string, params url.Values) (*http.Response, error) {
	endpointURL, err := url.Parse(endpoint)

	if err != nil {
		return &http.Response{}, err
	}

	if params == nil {
		params = endpointURL.Query()
	}

	params.Set("apikey", s.apiKey)

	endpointURL.RawQuery = params.Encode()

	requestURL := appendEndpoint(s.baseURL.String(), endpointURL.String())

	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		return &http.Response{}, err
	}

	return s.HTTPClient.Do(req)
}

func (s *Sonarr) put(endpoint string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)

	if err != nil {
		return &http.Response{}, err
	}

	endpointURL, err := url.Parse(endpoint)

	if err != nil {
		return &http.Response{}, err
	}

	params := endpointURL.Query()

	params.Set("apikey", s.apiKey)

	endpointURL.RawQuery = params.Encode()

	requestURL := appendEndpoint(s.baseURL.String(), endpointURL.String())

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(body))
	if err != nil {
		return &http.Response{}, err
	}

	return s.HTTPClient.Do(req)
}

func (s Sonarr) post(query string, body []byte) (*http.Response, error) {
	endpointURL, err := url.Parse(query)

	if err != nil {
		return &http.Response{}, err
	}

	client := http.Client{
		Timeout: time.Duration(s.Timeout) * time.Second,
	}

	requestURL := appendEndpoint(s.baseURL.String(), endpointURL.String())

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(body))

	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func (s *Sonarr) del(endpoint string, params url.Values) (*http.Response, error) {
	endpointURL, err := url.Parse(endpoint)

	if err != nil {
		return &http.Response{}, err
	}

	if params == nil {
		params = endpointURL.Query()
	}

	params.Set("apikey", s.apiKey)

	endpointURL.RawQuery = params.Encode()

	requestURL := appendEndpoint(s.baseURL.String(), endpointURL.String())

	req, err := http.NewRequest("DELETE", requestURL, nil)

	if err != nil {
		return &http.Response{}, err
	}

	return s.HTTPClient.Do(req)
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
