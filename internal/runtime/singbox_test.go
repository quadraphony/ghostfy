package runtime

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMinimalConfig(t *testing.T) {
	cfg := minimalConfig()
	outbounds, ok := cfg["outbounds"].([]map[string]any)
	if !ok {
		t.Fatalf("outbounds has wrong type: %T", cfg["outbounds"])
	}

	if len(outbounds) != 1 {
		t.Fatalf("expected one outbound, got %d", len(outbounds))
	}

	if got := outbounds[0]["type"]; got != "direct" {
		t.Fatalf("expected direct outbound, got %v", got)
	}
}

func TestWriteConfigCreatesJSONFile(t *testing.T) {
	t.Parallel()

	runner := NewSingboxTestRunner(&bytes.Buffer{}, &bytes.Buffer{})
	runner.tempDir = t.TempDir()

	path, err := runner.writeConfig()
	if err != nil {
		t.Fatalf("writeConfig returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, `"outbounds"`) {
		t.Fatalf("config missing outbounds: %s", content)
	}
	if !strings.Contains(content, `"direct"`) {
		t.Fatalf("config missing direct outbound: %s", content)
	}
}

func TestRunFailsWhenSingboxMissing(t *testing.T) {
	t.Parallel()

	runner := NewSingboxTestRunner(&bytes.Buffer{}, &bytes.Buffer{})
	runner.manager.lookPath = func(string) (string, error) {
		return "", exec.ErrNotFound
	}

	err := runner.Run(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "find sing-box binary") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWithShortLivedCommand(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	runner := NewSingboxTestRunner(&stdout, &stderr)
	runner.manager.lookPath = func(string) (string, error) {
		return "/bin/sh", nil
	}
	runner.manager.command = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		return exec.CommandContext(ctx, "/bin/sh", "-c", "printf 'started\\n'; exit 0")
	}
	runner.tempDir = t.TempDir()
	runner.startupWait = 50 * time.Millisecond

	if err := runner.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if !strings.Contains(stdout.String(), "started") {
		t.Fatalf("stdout did not include process output: %s", stdout.String())
	}
}

func TestStopProcessKillsStuckCommand(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("/bin/sh", "-c", "trap '' TERM; sleep 10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	defer cmd.Process.Kill()

	runner := NewSingboxTestRunner(&bytes.Buffer{}, &bytes.Buffer{})
	runner.manager.shutdownWait = 200 * time.Millisecond

	logger := bytes.NewBuffer(nil)
	err := runner.manager.Stop()
	if err == nil {
		t.Fatal("expected stop error without started manager")
	}

	testManager := NewManager(&bytes.Buffer{}, &bytes.Buffer{}, NewSingboxLogger(logger))
	testManager.shutdownWait = 200 * time.Millisecond
	testManager.cmd = cmd
	testManager.waitCh = make(chan error, 1)
	go func() {
		testManager.waitCh <- cmd.Wait()
	}()
	err = testManager.Stop()
	if err != nil {
		t.Fatalf("stopProcess returned error: %v", err)
	}
}

func TestConfigFileName(t *testing.T) {
	t.Parallel()

	got := configFileName(filepath.Join("/tmp", "ghostify-singbox.json"))
	if got != "ghostify-singbox.json" {
		t.Fatalf("unexpected base name: %s", got)
	}
}

func NewSingboxLogger(buf *bytes.Buffer) *loggingAdapter {
	return &loggingAdapter{buf: buf}
}

type loggingAdapter struct {
	buf *bytes.Buffer
}

func (l *loggingAdapter) Log(level, msg string, fields map[string]any) {
	l.buf.WriteString(level)
	l.buf.WriteString(msg)
	if fields != nil {
		l.buf.WriteString(" fields")
	}
}

func TestStopProcessNilProcess(t *testing.T) {
	t.Parallel()

	runner := NewManager(&bytes.Buffer{}, &bytes.Buffer{}, NewSingboxLogger(&bytes.Buffer{}))
	err := runner.Stop()
	if err == nil {
		t.Fatal("expected process error, got nil")
	}
	if !strings.Contains(err.Error(), "process not started") {
		t.Fatalf("unexpected error: %v", err)
	}
}
