# Pokédex CLI

> Explore different regions and catch Pokémon from your terminal.

A command-line Pokédex built with [Go](https://github.com/golang/go) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

<p align="center">
  <img width="100%" src="https://github.com/user-attachments/assets/c2404b49-1971-47e8-85f4-fecfb35e1105" alt="User catching a gastly from the terminal" />
</p>

## Motivation
I wanted to interact with the PokéAPI and get some hands-on practice building a CLI tool in Go. 
By connecting the terminal to the PokéAPI you can browse locations, run into new Pokémon, and catch them with mechanics that feel similar to the real thing.

## Features
* Scroll through locations
* Explore new areas
* Catch Pokémon
* Check out your collection and the stats of your Pokémon
* Caching for faster performance

## Quick Start

### Run
```go run main.go```

### Usage
```
Pokedex > map                         # Show the next 20 locations
Pokedex > explore canalave-city-area  # See what pops up there
Pokedex > catch pikachu               # Give catching Pikachu a shot
Pokedex > pokedex                     # List caught Pokémon
Pokedex > inspect pikachu             # Check Pikachu's stats, types, and moves
```

## Commands
| Command | What it does |
|---------|--------------|
| `map` | Jump to the next page of locations |
| `mapb` | Go back to the previous page |
| `explore <location>` | Find out which Pokémon live there |
| `catch <pokemon>` | Attempt catching a Pokémon |
| `inspect <pokemon>` | Get the details on a caught Pokémon |
| `pokedex` | List your caught Pokémon |
| `help` | List all commands |
| `exit` | Quit the app |

## How It Works
Catching odds depend on the Pokémon's base experience with catch rates anywhere from 5% to 95%. The app caches API responses for 5 seconds.
