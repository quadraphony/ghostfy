package openvpn

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestRunnerMissingBinary(t *testing.T) {
	t.Parallel()

	runner := NewRunner(&bytes.Buffer{}, &bytes.Buffer{})
	runner.lookPath = func(string) (string, error) {
		return "", exec.ErrNotFound
	}

	cfg := Config{OpenVPNConfig: "/tmp/ghostify.conf"}
	err := runner.Run(context.Background(), cfg)
	if err == nil || err.Error() != `find openvpn binary: executable file not found in $PATH` {
		t.Fatalf("expected missing binary error, got %v", err)
	}
}

func TestRunnerWithShortCommand(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	runner := NewRunner(&stdout, &stderr)
	runner.command = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "/bin/sh", "-c", "printf 'started\n'; sleep 1")
	}
	runner.lookPath = func(string) (string, error) {
		return "/bin/sh", nil
	}

	cfg := Config{OpenVPNConfig: "/tmp/config"}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := runner.Run(ctx, cfg)
	if err == nil {
		t.Fatal("expected context error")
	}
	if !bytes.Contains(stderr.Bytes(), []byte("started")) && !bytes.Contains(stdout.Bytes(), []byte("started")) {
		t.Fatalf("expected process output, got stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}
