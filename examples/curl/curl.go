package main

import (
	"strings"

	"github.com/reugn/packer"
)

// used in packer_test.go
func main() {
	packer.Packer{
		Filter: func(s string) bool {
			return strings.Contains(s, "Programming Language")
		},
		Map: func(s string) string {
			return strings.ToUpper(s)
		},
		Reduce: packer.Count,
	}.Execute()
}
