package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var db map[string]string

func init() {
	f, err := ioutil.ReadFile("db.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(f, &db); err != nil {
		log.Fatal(err)
	}
}

func doETL(in []byte) ([][]byte, error) {
	var res [][]byte
	s := string(in)
	if strings.HasPrefix(s, "+") {
		key := strings.Replace(s, "+", "", 1)
		value, err := getByKey(key)
		if err != nil {
			return nil, err
		}
		return append(res, []byte(value)), nil
	}
	return nil, nil
}

func getByKey(key string) (string, error) {
	value, exists := db[key]
	if !exists {
		return "", fmt.Errorf("%s not found", key)
	}
	return value, nil
}
