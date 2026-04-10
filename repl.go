package main

import (
	"bufio"
	"fmt"
	"math/rand"
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
		"catch": {
			name:        "catch",
			description: "Attempts to catch the specified pokemon and add it to the pokedex",
			callback: func(cfg *config, args []string) error {
				return commandCatch(cfg, args)
			},
		},
		"inspect": {
			name:        "inspect",
			description: "provides statistics of captured pokemon",
			callback: func(cfg *config, args []string) error {
				return commandInspect(cfg, args)
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows captured pokemons",
			callback: func(cfg *config, args []string) error {
				return commandPokedex(cfg)
			},
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		Cache:         pokecache.NewCache(5 * time.Second),
		caughtPokemon: map[string]pokeapi.Pokemon{},
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

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("No pokemon specified")
		return nil
	}
	data, err := pokeapi.GetPokemon(args[0], cfg.Cache)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	threshold := data.BaseExperience / 2
	result := rand.Intn(data.BaseExperience)
	if result < threshold {
		fmt.Printf("%s was caught!\n", args[0])
		fmt.Println("You may now inspect it with the inspect command.")
		cfg.caughtPokemon[data.Name] = data
	} else {
		fmt.Printf("%s escaped!\n", args[0])
	}
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("No pokemon specified")
		return nil
	}
	pokemon, ok := cfg.caughtPokemon[args[0]]
	if !ok {
		fmt.Println("Pokemon not captured")
		return nil
	}
	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height", pokemon.Height)
	fmt.Println("Weight", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeinfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeinfo.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *config) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf("  - %s\n", pokemon.Name)
	}
	return nil
}
