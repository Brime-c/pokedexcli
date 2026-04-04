package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Brime/pokedexcli/internal/pokecache"
)

type Shallow struct {
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type LocationArea struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

func ListLocations(url string, cache *pokecache.Cache) (Shallow, error) {
	val, ok := cache.Get(url)
	if ok {
		data := Shallow{}

		err := json.Unmarshal(val, &data)

		if err != nil {
			return Shallow{}, err
		}
		return data, nil
	}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	data := Shallow{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return Shallow{}, err
	}
	cache.Add(url, body)
	return data, nil
}

func ListPokemon(area string, cache *pokecache.Cache) ([]PokemonEncounter, error) {
	const baseURL = "https://pokeapi.co/api/v2/location-area/"
	fullUrl := baseURL + area
	val, ok := cache.Get(fullUrl)
	if ok {
		data := LocationArea{}

		err := json.Unmarshal(val, &data)

		if err != nil {
			return []PokemonEncounter{}, err
		}
		return data.PokemonEncounters, nil
	}
	res, err := http.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	data := LocationArea{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return []PokemonEncounter{}, err
	}
	cache.Add(fullUrl, body)
	return data.PokemonEncounters, nil

}
