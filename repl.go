package main

import (
	"strings"
)

func cleanInput(text string) []string {
	lowercase := strings.ToLower(text)
	var wordSlice []string
	wordSlice = strings.Fields(lowercase)
	return wordSlice
}
