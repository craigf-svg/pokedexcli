package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

type LocationResult struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) FetchLocations(url *string) (LocationResult, error) {
	getURL := baseURL + "/location-area"
	if url != nil && *url != "" {
		getURL = *url
	}

	resp, err := c.httpClient.Get(getURL)
	if err != nil {
		return LocationResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationResult{}, err
	}

	var locations LocationResult
	if err := json.Unmarshal(body, &locations); err != nil {
		return LocationResult{}, err
	}

	return locations, nil
}
