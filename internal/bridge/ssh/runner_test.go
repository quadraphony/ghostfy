package ssh

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
)

func TestRunnerMissingSSH(t *testing.T) {
	t.Parallel()

	runner := NewRunner(&bytes.Buffer{}, &bytes.Buffer{})
	runner.launcher = &mockLauncher{
		path: "",
		err:  exec.ErrNotFound,
	}

	_, err := runner.launcher.LookPath("ssh")
	if err == nil {
		t.Fatal("expected error")
	}
}

type mockLauncher struct {
	path string
	err  error
}

func (m *mockLauncher) LookPath(name string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.path, nil
}

func (m *mockLauncher) Exec(ctx context.Context, bin string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "/bin/sh", "-c", "sleep 1")
}
