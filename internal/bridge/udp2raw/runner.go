package udp2raw

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

type Runner struct {
	stdout io.Writer
	stderr io.Writer
	cmd    func(ctx context.Context, name string, args ...string) *exec.Cmd
	look   func(string) (string, error)
	log    *logging.Logger
}

func NewRunner(stdout, stderr io.Writer) *Runner {
	return &Runner{
		stdout: stdout,
		stderr: stderr,
		cmd:    exec.CommandContext,
		look:   exec.LookPath,
		log:    logging.New(stdout),
	}
}

func (r *Runner) Run(ctx context.Context, cfg Config) error {
	bin := cfg.Executable
	if bin == "" {
		var err error
		bin, err = r.look("udp2raw")
		if err != nil {
			return fmt.Errorf("find udp2raw binary: %w", err)
		}
	}

	cmd := r.cmd(ctx, bin, cfg.Args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	r.log.Log("info", "udp2raw started", map[string]any{"pid": cmd.Process.Pid})

	var wg sync.WaitGroup
	wg.Add(2)
	go r.stream(&wg, stdout, r.stdout, "stdout")
	go r.stream(&wg, stderr, r.stderr, "stderr")

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-ctx.Done():
		cmd.Process.Signal(os.Interrupt)
		wg.Wait()
		return ctx.Err()
	case err := <-done:
		wg.Wait()
		return err
	}
}

func (r *Runner) stream(wg *sync.WaitGroup, in io.Reader, out io.Writer, label string) {
	defer wg.Done()
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		fmt.Fprintf(out, "[udp2raw/%s] %s\n", label, scanner.Text())
	}
	if err := scanner.Err(); err != nil && !strings.Contains(err.Error(), "file already closed") {
		fmt.Fprintf(out, "[udp2raw/%s] stream error: %v\n", label, err)
	}
}
