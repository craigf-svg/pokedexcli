package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"pokedex/internal/pokeapi"
)

type config struct {
	next       *string
	previous   *string
	pokeClient *pokeapi.Client
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := &config{
		pokeClient: pokeapi.NewClient(5 * time.Second),
	}

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
		"explore": {
			name:         "explore",
			description:  "Explore a specific location",
			callback:     commandExplore,
		},
	}

	for {
		fmt.Printf("Pokedex > ")

		if !scanner.Scan() {
			return
		}

		input := scanner.Text()
		clean := cleanInput(input)
		arg := ""	

		if len(clean) == 0 {
			continue
		}

		command := clean[0]
		
		if len(clean) == 2 {
			arg = clean[1]
		}
		if cli, ok := cliCommands[command]; ok {
			if err := cli.callback(cfg, arg); err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(cfg *config, arg string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, arg string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("Map: Print the current map page")
	fmt.Println("Mapb: Print the previous map page")
	fmt.Println("Explore: Explore a specific area")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandExplore(cfg *config, arg string) error {
	location, err := cfg.pokeClient.ExploreLocation(arg)
	if err !=nil {
		return err
	}

	for _, encounter := range location.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandMap(cfg *config, arg string) error {
	locations, err := cfg.pokeClient.FetchLocations(cfg.next)
	if err != nil {
		return err
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapB(cfg *config, arg string) error {
	if cfg.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locations, err := cfg.pokeClient.FetchLocations(cfg.previous)
	if err != nil {
		return err
	}

	cfg.next = locations.Next
	cfg.previous = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	return strings.Fields(lower)
}
