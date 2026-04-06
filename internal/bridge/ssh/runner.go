package ssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/quadraphony/ghostfy/internal/logging"
)

type commandLauncher interface {
	LookPath(string) (string, error)
	Exec(ctx context.Context, bin string, args ...string) *exec.Cmd
}

type Runner struct {
	stdout   io.Writer
	stderr   io.Writer
	launcher commandLauncher
	logger   *logging.Logger
}

func NewRunner(stdout, stderr io.Writer) *Runner {
	return &Runner{
		stdout:   stdout,
		stderr:   stderr,
		launcher: &defaultLauncher{},
		logger:   logging.New(stdout),
	}
}

func (r *Runner) Run(ctx context.Context, cfg Config) error {
	bin := cfg.Binary
	if bin == "" {
		var err error
		bin, err = r.launcher.LookPath("ssh")
		if err != nil {
			return fmt.Errorf("find ssh binary: %w", err)
		}
	}
	args := append([]string{cfg.Connection}, cfg.Args...)
	cmd := r.launcher.Exec(ctx, bin, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	r.logger.Log("info", "ssh bridge running", map[string]any{"target": cfg.Connection, "pid": cmd.Process.Pid})

	var wg sync.WaitGroup
	wg.Add(2)
	go r.copyStream(&wg, stdoutPipe, r.stdout, "stdout")
	go r.copyStream(&wg, stderrPipe, r.stderr, "stderr")

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		_ = cmd.Process.Signal(os.Interrupt)
		wg.Wait()
		return ctx.Err()
	case err := <-done:
		wg.Wait()
		return err
	}
}

func (r *Runner) copyStream(wg *sync.WaitGroup, in io.Reader, out io.Writer, stream string) {
	defer wg.Done()
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		fmt.Fprintf(out, "[ssh/%s] %s\n", stream, scanner.Text())
	}
	if err := scanner.Err(); err != nil && !strings.Contains(err.Error(), "file already closed") {
		fmt.Fprintf(out, "[ssh/%s] stream error: %v\n", stream, err)
	}
}

type defaultLauncher struct{}

func (defaultLauncher) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func (defaultLauncher) Exec(ctx context.Context, bin string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, bin, args...)
}
