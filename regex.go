package main

import (
	"bytes"
	"fmt"
	"log"
)

func regexRoute(route string) (string, string) {
	indices, err := findParamsIndices(route)
	if err != nil {
		log.Fatal(err)
	}

	var buffer bytes.Buffer
	buffer.WriteString(`^`)

	// Static path
	if indices == nil {
		buffer.WriteString(route)
		buffer.WriteString(`$`)
		return buffer.String(), ""
	}

	var unpack bytes.Buffer
	unpack.WriteString(``)

	var end int
	for i := 0; i < len(indices); i += 2 {
		raw := route[end:indices[i]]
		buffer.WriteString(raw)
		end = indices[i+1]

		name := route[indices[i]+1 : end-1]
		if name == "" {
			log.Fatal("regex: param name not found")
		}

		unpack.WriteString(raw)
		unpack.WriteString(`(?P<`)
		unpack.WriteString(name)
		unpack.WriteString(`>[^/]+)`)

		buffer.WriteString(`[^/]+`)
	}
	buffer.WriteString(route[end:])
	buffer.WriteString(`$`)
	unpack.WriteString(`$`)

	return buffer.String(), unpack.String()
}

func findParamsIndices(url string) ([]int, error) {
	var level, index int
	var indices []int

	for i := 0; i < len(url); i++ {
		switch url[i] {
		case '{':
			if level++; level == 1 {
				index = i
			} else if level > 1 {
				return nil, fmt.Errorf("regex: braces in %v", url)
			}
		case '}':
			if level--; level == 0 {
				indices = append(indices, index, i+1)
			} else if level < 0 {
				return nil, fmt.Errorf("regex: braces in %v", url)
			}
		}
	}

	return indices, nil
}
