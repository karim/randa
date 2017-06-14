package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func checkDatabase(database string) {
	var err error
	db, err = sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func runQuery(url string, e endpoint) []byte {
	query := e.query

	if e.unpack != "" {
		regex := regexp.MustCompile(e.unpack)
		match := regex.FindStringSubmatch(url)
		param := map[string]string{}
		for i, name := range regex.SubexpNames() {
			if i > 0 {
				param[name] = match[i]
			}
		}
		for key, value := range param {
			query = strings.Replace(query, "{"+key+"}", value, -1)
		}
	}

	return runSQLQuery(query)
}

func runSQLQuery(query string) []byte {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, err := rows.ColumnTypes()
	if err != nil {
		log.Fatal(err)
	}

	length := len(columns)

	pointers := make([]interface{}, length)
	values := make([]interface{}, length)

	for i := range columns {
		pointers[i] = &values[i]
	}

	table := make([]map[string]interface{}, 0)
	var entry map[string]interface{}

	rowNum := 0

	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			log.Fatal(err)
		}

		entry = make(map[string]interface{})

		for i, value := range values {
			var v interface{}
			b, ok := value.([]byte)
			if ok {
				v = string(b)
			} else {
				v = value
			}

			// Don't show 'NULL' columns
			if value != nil {
				entry[columns[i].Name()] = v
			}
		}
		table = append(table, entry)

		rowNum++
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	var jsonResult []byte
	if rowNum == 0 {
		// No results, return an empty object
		return []byte(`{}`)
	} else if rowNum == 1 {
		// Only one result, return the JSON object
		jsonResult, err = json.Marshal(entry)
	} else {
		// More than one result, return the JSON array
		jsonResult, err = json.Marshal(table)
	}
	if err != nil {
		log.Fatal(err)
	}
	return jsonResult
}
