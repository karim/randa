package main

import (
	"log"
	"regexp"
)

// routeMatch takes a url and returns the endpoint id that matches it
func routeMatch(url string) int {
	for id, endpoint := range endpoints {
		matched, err := regexp.Match(endpoint.regex, []byte(url))
		if err != nil {
			log.Fatal(err)
		}

		if matched {
			return id
		}
	}

	return -1
}
