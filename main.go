package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"math/rand"
	"time"
	"errors"

	"pokedex/internal/pokeapi"
)

type config struct {
	next       *string
	previous   *string
	pokeClient *pokeapi.Client
	pokedex	   map[string]pokeapi.PokemonResult
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
		pokedex: make(map[string]pokeapi.PokemonResult),
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
		"catch": {
			name:         "catch",
			description:  "Catch a pokemon",
			callback:     commandCatch,
		},
		"inspect": {
			name:         "inspect",
			description:  "Inspect a caught pokemon",
			callback:     commandInspect,
		},
		"pokedex": {
			name:         "pokedex",
			description:  "View your pokedex",
			callback:     commandPokedex,
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
	fmt.Println("Catch: Catch a specific pokemon")
	fmt.Println("Inspect: Inspect a caught pokemon")
	fmt.Println("Pokedex: View your pokedex")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func commandPokedex(cfg *config, arg string) error {
	fmt.Println("Your Pokedex:")
	for _, val := range cfg.pokedex {
		fmt.Printf("- %s\n", val.Name)
	}

	return nil
}

func commandInspect(cfg *config, arg string) error {
	if arg == "" {
		return errors.New("Pass pokemon to inspect") 
	}
	pokemon, exists := cfg.pokedex[arg]

	if !exists {
		return errors.New("You haven't caught that Pokemon yet")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Height: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func commandCatch(cfg *config, arg string) error {
	fmt.Printf("catch %s\n", arg)	
	pokemon, err := cfg.pokeClient.CatchPokemon(arg)

	if err != nil {
		return err
	}

	fmt.Printf("Base Experience: %v\n", pokemon.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	catchChance := 100 - (pokemon.BaseExperience / 3)

	if catchChance < 5 {
		catchChance = 5
	}
	if catchChance > 95 {
		catchChance = 95
	}

	if rand.Intn(100) < catchChance {
    fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

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
