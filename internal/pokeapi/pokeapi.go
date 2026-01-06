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

type ExploreResult struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (c *Client) ExploreLocation(locationAreaName string) (ExploreResult, error) {
	if locationAreaName == "" {
		fmt.Println("Must type an area to explore")
		return ExploreResult{}, nil 
	}
	getURL := baseURL + "/location-area/" + locationAreaName
	fmt.Println(getURL)

	if data, ok := c.cache.Get(getURL); ok {
		return decodeExploreResult(data)
	}

	resp, err := c.httpClient.Get(getURL)
	if err != nil {
			return ExploreResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
			return ExploreResult{}, err
	}

	c.cache.Add(getURL, body)

	return decodeExploreResult(body)
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

func decodeExploreResult(data []byte) (ExploreResult, error) {
    var pokemon ExploreResult
    if err := json.Unmarshal(data, &pokemon); err != nil {
        return ExploreResult{}, err
    }
    return pokemon, nil
}
