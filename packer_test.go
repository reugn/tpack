package packer

import (
	"os/exec"
	"testing"
)

func TestPacker(t *testing.T) {
	cmd := "curl https://golang.org/ | go run ./examples/curl/curl.go"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		t.Fatalf("Failed to execute command: %s", cmd)
	}
	assertEqual(t, "2", string(out[:]))
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
