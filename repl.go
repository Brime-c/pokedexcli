package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
	var comms map[string]cliCommand

	comms = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback: func() error {
				return commandHelp(comms)
			},
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		clean := cleanInput(scanner.Text())
		firstWord := clean[0]
		comm, ok := comms[firstWord]
		if ok {
			err := comm.callback()
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
