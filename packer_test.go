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

	tpack.NewPacker(in, &out, &err, tpack.NewFunctionProcessor(
		func(in []byte) ([][]byte, error) {
			var res [][]byte
			s := string(in)
			_, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return nil, errors.New(s)
			}
			return append(res, []byte(s)), nil
		},
	)).Execute()

	assertEqual(t, out.String(), "1\n2\n")
	assertEqual(t, err.String(), "a\nb\nc\n")
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
