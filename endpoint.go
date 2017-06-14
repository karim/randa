package main

type endpoint struct {
	method string
	url    string
	regex  string
	unpack string
	query  string
}

var endpoints []endpoint

func addEndpoint(method, url, regex, unpack, query string) {
	endpoints = append(endpoints, endpoint{method, url, regex, unpack, query})
}
