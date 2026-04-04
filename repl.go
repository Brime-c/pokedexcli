package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Brime/pokedexcli/internal/pokeapi"
	"github.com/Brime/pokedexcli/internal/pokecache"
)

func startRepl() {
	var comms map[string]cliCommand

	comms = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback: func(_ *config, _ []string) error {
				return commandExit()
			},
		},

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback: func(_ *config, _ []string) error {
				return commandHelp(comms)
			},
		},

		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback: func(cfg *config, args []string) error {
				return commandMap(cfg)
			},
		},

		"mapb": {
			name:        "mapb",
			description: "Displays the last 20 locations",
			callback: func(cfg *config, args []string) error {
				return commandMapb(cfg)
			},
		},
		"explore": {
			name:        "explore",
			description: "Displays the pokemons present in the requested area",
			callback: func(cfg *config, args []string) error {
				return commandExplore(cfg, args)
			},
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		Cache: pokecache.NewCache(5 * time.Second),
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		clean := cleanInput(scanner.Text())
		firstWord := clean[0]
		args := clean[1:]
		comm, ok := comms[firstWord]
		if ok {
			err := comm.callback(cfg, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}
func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	var wordSlice []string
	wordSlice = strings.Fields(lowercase)
	return wordSlice
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(comms map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, comm := range comms {
		fmt.Printf("%s: %s\n", comm.name, comm.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if cfg.Next != nil {
		url = *cfg.Next
	}
	data, err := pokeapi.ListLocations(url, cfg.Cache)
	if err != nil {
		return err
	}
	cfg.Next = data.Next
	cfg.Previous = data.Previous
	for _, location := range data.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	url := *cfg.Previous

	data, err := pokeapi.ListLocations(url, cfg.Cache)
	if err != nil {
		return err
	}
	cfg.Next = data.Next
	cfg.Previous = data.Previous
	for _, location := range data.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("No area provided")
		return nil
	}
	data, err := pokeapi.ListPokemon(args[0], cfg.Cache)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, encounter := range data {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
