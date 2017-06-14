package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func startServer(port int) {
	http.HandleFunc("/", handler)

	log.Println()
	strPort := strconv.Itoa(port)
	log.Printf("Listening on port: %s\n", strPort)
	log.Println("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(":"+strPort, nil))
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
		r = []byte("{}")
	} else {
		r = runQuery(request.URL.String(), endpoints[id])
	}

	response.Write(r)

	finish := time.Since(start)
	log.Printf("[%s] %q %v \n", request.Method, request.URL, finish)
}
