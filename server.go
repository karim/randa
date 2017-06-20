package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func startServer(config Config) {
	http.HandleFunc("/", handler)

	strPort := strconv.Itoa(config.Port)

	// Start either a HTTP or HTTPS server
	var err error
	if config.HTTPS.Certificate != "" && config.HTTPS.Key != "" {
		log.Println("HTTPS server listening on port:", strPort)
		err = http.ListenAndServeTLS(":"+strPort, config.HTTPS.Certificate, config.HTTPS.Key, nil)
	} else {
		log.Println("HTTP server listening on port:", strPort)
		err = http.ListenAndServe(":"+strPort, nil)
	}
	log.Fatal(err)
}

func handler(response http.ResponseWriter, request *http.Request) {
	// Set this for every request
	response.Header().Set("Content-Type", "application/json")

	start := time.Now()

	if request.Method != "GET" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var r []byte

	// Check if the url matches any endpoint
	id := routeMatch(request.URL.String())
	if id == -1 {
		// no match
		response.WriteHeader(http.StatusNotFound)
		r = []byte(`{"message":"Not Found"}`)
	} else {
		r = runQuery(request.URL.String(), endpoints[id])
	}

	response.Write(r)

	finish := time.Since(start)
	log.Printf("[%s] %q %v \n", request.Method, request.URL, finish)
}
