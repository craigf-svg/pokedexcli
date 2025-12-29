package main

import (
	"fmt" 
	"strings"
	"bufio"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
				name:				 "help",
				description: "Displays a help message",
				callback:		 commandHelp,
		},
	}

	for {
		fmt.Printf("Pokedex > ")

		if !scanner.Scan() {
			return // or break, if input is closed
    }

    input := scanner.Text()
    clean := cleanInput(input)
		command := clean[0]
		
		if cli, ok := cliCommands[command]; ok {
			cli.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}

	return
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")	
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")

	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	/* for _, message := range cliCommands {
		fmt.Println(message.description)
	} */
	return nil
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	return strings.Fields(lower)
}
