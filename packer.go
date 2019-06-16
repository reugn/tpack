package packer

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Packer - command boxing struct
// Wraps go application to act as a Unix pipeline process
// Packer anatomy:
// ------------------------------------------------
//            _________________________
//    stdin  |                         |  stdout
//   ------> | Filter -> Map -> Reduce | ------->
//           |_________________________| [stderr]
//
// ------------------------------------------------
type Packer struct {
	Filter func(string) bool
	Map    func(string) string
	Reduce func([]string) string
}

// Execute command
func (p Packer) Execute() {
	stdout := os.Stdout
	stderr := os.Stderr

	info, err := os.Stdin.Stat()
	if err != nil {
		stderr.Write([]byte(err.Error()))
		os.Exit(-1)
	}
	if info.Mode()&os.ModeNamedPipe == 0 {
		stderr.Write([]byte("named pipe (FIFO)"))
		os.Exit(-1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	var r []string
	for scanner.Scan() {
		line := scanner.Text()
		if p.Filter != nil && !p.Filter(line) {
			continue
		}
		if p.Map != nil {
			line = p.Map(line)
		}
		if p.Reduce != nil {
			r = append(r, line)
		} else {
			stdout.Write([]byte(line + "\n"))
		}
	}
	if p.Reduce != nil {
		outStr := []byte(p.Reduce(r))
		stdout.Write([]byte(outStr))
	}
}

// MkString with delimiter reducer
func MkString(delimiter string) func([]string) string {
	return func(s []string) string {
		return strings.Join(s, delimiter) + "\n"
	}
}

// Count reducer
var Count = func(s []string) string {
	return strconv.Itoa(len(s))
}
