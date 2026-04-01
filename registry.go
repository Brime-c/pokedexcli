package main

import "github.com/Brime/pokedexcli/internal/pokecache"

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
}
