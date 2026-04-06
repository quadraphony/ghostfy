package udp2raw

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
)

func TestUDP2RAWRunnerMissingBinary(t *testing.T) {
	t.Parallel()

	runner := NewRunner(&bytes.Buffer{}, &bytes.Buffer{})
	runner.look = func(string) (string, error) {
		return "", exec.ErrNotFound
	}

	cfg := Config{Args: []string{"-l", "12345"}}
	if err := runner.Run(context.Background(), cfg); err == nil {
		t.Fatal("expected lookup error")
	}
}
