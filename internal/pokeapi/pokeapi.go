package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
	"fmt"

	"pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

type Client struct {
	httpClient http.Client
	cache			 *pokecache.Cache
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(5 * time.Second),
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
	getURL := baseURL + "/location-area?offset=0&limit=20"
	if url != nil && *url != "" {
		getURL = *url
	}

	fmt.Println("Get key:", getURL)

	if data, ok := c.cache.Get(getURL); ok {
			return decodeLocations(data)
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

	c.cache.Add(getURL, body)

	return decodeLocations(body)
}

func decodeLocations(data []byte) (LocationResult, error) {
    var locations LocationResult
    if err := json.Unmarshal(data, &locations); err != nil {
        return LocationResult{}, err
    }
    return locations, nil
}
