package main

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config file (.yaml) structure
type Config struct {
	HTTPS struct {
		Certificate string `yaml:"certificate"`
		Key         string `yaml:"key"`
	} `yaml:"https"`
	Port      int    `yaml:"port"`
	Database  string `yaml:"database"`
	Endpoints []struct {
		URL   string `yaml:"url"`
		Query string `yaml:"query"`
	} `yaml:"endpoints"`
}

func loadConfig() Config {
	var config Config

	// Read the config file
	file, err := ioutil.ReadFile("randa.yaml")
	if err != nil {
		log.Fatal("config:", err)
	}

	// Parse it
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("config:", err)
	}

	// Use the default port number '80' for HTTP and '443' for HTTPS if port number is not set
	if config.Port == 0 {
		if config.HTTPS.Certificate != "" && config.HTTPS.Key != "" {
			config.Port = 443
		} else {
			config.Port = 80
		}
	}

	// Check if the database is set
	if config.Database == "" {
		log.Fatal("Database not found!")
	}
	log.Println("Database:", config.Database)

	// Try to open it
	checkDatabase(config.Database)

	// Check if there is at least one endpoint
	if len(config.Endpoints) == 0 {
		log.Fatal("endpoints: 0 endpoints found")
	}
	log.Println("Endpoints:", len(config.Endpoints))

	// Check the endpoints
	for i, endpoint := range config.Endpoints {
		if endpoint.URL == "" {
			log.Fatal("Endpoints: 'url' cannot be empty")
		}
		log.Printf("\t#%d => %s\n", i+1, endpoint.URL)

		// Split the 'url' to get the HTTP method
		url := strings.Fields(endpoint.URL)
		if len(url) != 2 {
			log.Fatal("endpoints: 'url' is in the wrong format")
		}

		// Only support GET method for now
		if url[0] != "GET" {
			log.Fatal("endpoints: only 'GET' is supported")
		}

		// Make sure 'query' is not empty
		if endpoint.Query == "" {
			log.Fatal("endpoints: 'query' cannot be empty")
		}

		regex, unpack := regexRoute(url[1])

		// All good, add the new endpoint
		addEndpoint(url[0], url[1], regex, unpack, endpoint.Query)
	}

	return config
}
