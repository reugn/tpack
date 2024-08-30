package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type store map[string]string

func (s store) get(key string) (string, error) {
	time.Sleep(200 * time.Millisecond) // simulate latency
	value, exists := s[key]
	if !exists {
		return "", fmt.Errorf("%s not found", key)
	}
	return value, nil
}

// emulates a database store
var db store

// initialize the global db from the file
func init() {
	data, err := os.ReadFile("db.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(data, &db); err != nil {
		log.Fatal(err)
	}
}

func doETL(in string) ([]string, error) {
	var result []string
	key, found := strings.CutPrefix(in, "+")
	if found {
		value, err := db.get(key)
		if err != nil {
			return nil, err
		}
		return append(result, value), nil
	}
	return result, nil
}
