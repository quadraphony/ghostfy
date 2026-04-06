package openvpn

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/quadraphony/ghostfy/internal/logging"
)

type execCommand func(ctx context.Context, name string, args ...string) *exec.Cmd

type logger interface {
	Log(level, msg string, fields map[string]any)
}

type Runner struct {
	stdout   io.Writer
	stderr   io.Writer
	command  execCommand
	lookPath func(string) (string, error)
	log      logger
}

func NewRunner(stdout, stderr io.Writer) *Runner {
	return &Runner{
		stdout:   stdout,
		stderr:   stderr,
		command:  exec.CommandContext,
		lookPath: exec.LookPath,
		log:      logging.New(stdout),
	}
}

func (r *Runner) Run(ctx context.Context, cfg Config) error {
	execPath := cfg.Executor
	if execPath == "" {
		var err error
		execPath, err = r.lookPath("openvpn")
		if err != nil {
			return fmt.Errorf("find openvpn binary: %w", err)
		}
	}

	args := append([]string{"--config", cfg.OpenVPNConfig}, cfg.Args...)
	cmd := r.command(ctx, execPath, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("open stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start openvpn: %w", err)
	}

	r.log.Log("info", "openvpn bridge running", map[string]any{"pid": cmd.Process.Pid, "config": cfg.OpenVPNConfig})

	var wg sync.WaitGroup
	wg.Add(2)
	go r.streamPipe(&wg, r.stdout, "stdout", stdoutPipe)
	go r.streamPipe(&wg, r.stderr, "stderr", stderrPipe)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		if err := r.shutdown(cmd.Process); err != nil {
			return err
		}
		wg.Wait()
		return ctx.Err()
	case err := <-done:
		wg.Wait()
		if err != nil {
			return fmt.Errorf("openvpn exited: %w", err)
		}
		return nil
	}
}

func (r *Runner) shutdown(process *os.Process) error {
	r.log.Log("info", "stopping openvpn bridge", nil)
	if err := process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return fmt.Errorf("signal openvpn: %w", err)
	}
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			if err := process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
				return fmt.Errorf("kill openvpn: %w", err)
			}
			return nil
		case <-ticker.C:
			if err := process.Signal(syscall.Signal(0)); err != nil {
				return nil
			}
		}
	}
}

func (r *Runner) streamPipe(wg *sync.WaitGroup, out io.Writer, name string, pipe io.Reader) {
	defer wg.Done()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Fprintf(out, "[openvpn/%s] %s\n", name, scanner.Text())
	}
}
