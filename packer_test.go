package tpack_test

import (
	"bytes"
	"errors"
	"strconv"
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

func assertEqual(t *testing.T, a, b any) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
