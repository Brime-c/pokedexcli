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
