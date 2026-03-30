package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func ListLocations(url string) (Shallow, error) {
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
	return data, nil
}
