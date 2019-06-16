package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/reugn/packer"
)

func main() {
	var db map[string]string
	f, _ := ioutil.ReadFile("db.json")
	json.Unmarshal(f, &db)
	packer.Packer{
		Filter: func(s string) bool {
			return strings.HasPrefix(s, "+")
		},
		Map: func(s string) string {
			s = strings.Replace(s, "+", "", 1)
			return fmt.Sprintf("%s -> %s", s, db[s])
		},
		Reduce: packer.MkString(", "),
	}.Execute()
}
