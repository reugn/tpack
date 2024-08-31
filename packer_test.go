package tpack_test

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/reugn/tpack"
)

func TestPacker(t *testing.T) {
	in := bytes.NewBufferString("a\nb\n1\nc\n2")
	var out bytes.Buffer
	var err bytes.Buffer

	tpack.NewPacker(in, &out, &err, tpack.NewProcessor(
		func(in []byte) ([][]byte, error) {
			var result [][]byte
			str := string(in)
			_, err := strconv.Atoi(str)
			if err != nil {
				return nil, errors.New(str)
			}
			return append(result, []byte(str)), nil
		},
	)).Execute()

	assertEqual(t, out.String(), "1\n2\n")
	assertEqual(t, err.String(), "a\nb\nc\n")
}

func TestPacker_parallel(t *testing.T) {
	in := bytes.NewBufferString("a\nb\n1\nc\n2")
	var out bytes.Buffer
	var err bytes.Buffer

	tpack.NewPacker(in, &out, &err, tpack.NewProcessor(
		func(str string) ([]string, error) {
			var result []string
			_, err := strconv.Atoi(str)
			if err != nil {
				return nil, errors.New(str)
			}
			return append(result, str), nil
		},
		tpack.Parallel(3),
	)).Execute()

	assertEqual(t, len(strings.ReplaceAll(out.String(), "\n", "")), 2)
	assertEqual(t, len(strings.ReplaceAll(err.String(), "\n", "")), 3)
}

func assertEqual(t *testing.T, a, b any) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
