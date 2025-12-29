package main

import (
	"fmt" 
	"strings"
	"bufio"
	"os"
	"net/http"
	"io"
	"log"
	"encoding/json"
)

type config struct {
    next     *string
    previous *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}


func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cliCommands := map[string]cliCommand{
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
		"help": {
			  name:        "help",
		    description: "Displays a help message",
			  callback:    commandHelp,
		}, 
		"map": {
			  name:        "map",
        description: "Print next map of the area",
        callback:    commandMap,
		},
		"mapb": {
		    name:        "mapb",
			  description: "Print previous map of the area",
			  callback:    commandMapB,
		},
	}

	cfg := &config{}

	for {
		fmt.Printf("Pokedex > ")

		if !scanner.Scan() {
			return // or break, if input is closed
    }

    input := scanner.Text()
    clean := cleanInput(input)
		command := clean[0]
		
		if cli, ok := cliCommands[command]; ok {
			cli.callback(cfg)
		} else {
			fmt.Println("Unknown command")
		}
	}

	return
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")	
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	fmt.Println("Map: Print the current map page")
	fmt.Println("Mapb: Print the previous map page")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	/* for _, message := range cliCommands {
		fmt.Println(message.description)
	} */
	return nil
}

type locationResult struct {
    Count    int       `json:"count"`
    Next     *string   `json:"next"`
    Previous *string   `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}

func commandMap(cfg *config) error {
	fmt.Println("Map")
	getURL := "https://pokeapi.co/api/v2/location-area"
	if cfg.next != nil && *cfg.next != "" {
			getURL = *cfg.next
	}

	res, err := http.Get(getURL)

	if err != nil {
		log.Fatal(err)
	}
	locations := locationResult{}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	unmarshalErr := json.Unmarshal(body, &locations)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	// DEBUG PRINT NEXT
	if cfg.next == nil {
		fmt.Println("No next, you're on the last page")
	} else {
		fmt.Println(*cfg.next)
	}
	// DEBUG PRINT PREVIOUS
	if cfg.previous == nil {
		fmt.Println("No previous, you're on the first page")
	} else {
		fmt.Println(*cfg.previous)
	}

	for _, loc := range locations.Results {
    fmt.Println(loc.Name)
	}
	return nil
}

func commandMapB(cfg *config) error {
	fmt.Println("MapB")
	if cfg.previous == nil {
    fmt.Println("you're on the first page")
    return nil
	}
	getURL := ""
	if cfg.previous != nil && *cfg.previous != "" {
			getURL = *cfg.previous
	} else {
		fmt.Println("Previous url could not be found")
	}

	res, err := http.Get(getURL)

	if err != nil {
		log.Fatal(err)
	}
	locations := locationResult{}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	unmarshalErr := json.Unmarshal(body, &locations)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous
	// DEBUG PRINT NEXT
	if cfg.next == nil {
		fmt.Println("No next, you're on the last page")
	} else {
		fmt.Println(*cfg.next)
	}
	// DEBUG PRINT PREVIOUS
	if cfg.previous == nil {
		fmt.Println("No previous, you're on the first page")
	} else {
		fmt.Println(*cfg.previous)
	}

	for _, loc := range locations.Results {
    fmt.Println(loc.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	return strings.Fields(lower)
}
