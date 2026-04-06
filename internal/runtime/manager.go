package runtime

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
)

const (
	defaultStartupWait  = 2 * time.Second
	defaultShutdownWait = 5 * time.Second
)

type State string

const (
	StateInit     State = "INIT"
	StateStarting State = "STARTING"
	StateRunning  State = "RUNNING"
	StateStopping State = "STOPPING"
	StateStopped  State = "STOPPED"
	StateError    State = "ERROR"
)

type execCommand func(ctx context.Context, name string, args ...string) *exec.Cmd

type logger interface {
	Log(level, msg string, fields map[string]any)
}

type Manager struct {
	mu           sync.Mutex
	state        State
	stdout       io.Writer
	stderr       io.Writer
	command      execCommand
	lookPath     func(string) (string, error)
	startupWait  time.Duration
	shutdownWait time.Duration
	log          logger
	cmd          *exec.Cmd
	waitCh       chan error
	streamWG     sync.WaitGroup
}

func NewManager(stdout, stderr io.Writer, log logger) *Manager {
	return &Manager{
		state:        StateInit,
		stdout:       stdout,
		stderr:       stderr,
		command:      exec.CommandContext,
		lookPath:     exec.LookPath,
		startupWait:  defaultStartupWait,
		shutdownWait: defaultShutdownWait,
		log:          log,
	}
}

func (m *Manager) State() State {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.state
}

func (m *Manager) Start(ctx context.Context, configPath string) error {
	m.mu.Lock()
	if m.state == StateStarting || m.state == StateRunning || m.state == StateStopping {
		m.mu.Unlock()
		return fmt.Errorf("cannot start from state %s", m.state)
	}
	m.state = StateStarting
	m.mu.Unlock()

	binaryPath, err := m.lookPath("sing-box")
	if err != nil {
		m.setState(StateError)
		return fmt.Errorf("find sing-box binary: %w", err)
	}

	m.log.Log("info", "resolved sing-box binary", map[string]any{"path": binaryPath})
	m.log.Log("info", "starting sing-box runtime", map[string]any{"config": configPath})

	cmd := m.command(ctx, binaryPath, "run", "-c", configPath)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		m.setState(StateError)
		return fmt.Errorf("open stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		m.setState(StateError)
		return fmt.Errorf("open stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		m.setState(StateError)
		return fmt.Errorf("start sing-box: %w", err)
	}

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	m.mu.Lock()
	m.cmd = cmd
	m.waitCh = waitCh
	m.streamWG.Add(2)
	go m.streamPipe(m.stdout, "stdout", stdoutPipe)
	go m.streamPipe(m.stderr, "stderr", stderrPipe)
	m.mu.Unlock()

	select {
	case err := <-waitCh:
		m.streamWG.Wait()
		if err != nil {
			m.setState(StateError)
			return fmt.Errorf("sing-box exited during startup: %w", err)
		}
		m.setState(StateStopped)
		return nil
	case <-time.After(m.startupWait):
		m.setState(StateRunning)
		m.log.Log("info", "sing-box entered running state", map[string]any{"pid": cmd.Process.Pid})
		return nil
	case <-ctx.Done():
		_ = m.Stop()
		return ctx.Err()
	}
}

func (m *Manager) Wait() error {
	m.mu.Lock()
	waitCh := m.waitCh
	m.mu.Unlock()

	if waitCh == nil {
		return errors.New("runtime not started")
	}

	err := <-waitCh
	m.streamWG.Wait()
	if err != nil {
		m.setState(StateError)
		return fmt.Errorf("wait for sing-box: %w", err)
	}

	m.setState(StateStopped)
	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	cmd := m.cmd
	if cmd == nil || cmd.Process == nil {
		m.mu.Unlock()
		return errors.New("process not started")
	}
	m.state = StateStopping
	waitCh := m.waitCh
	m.mu.Unlock()

	process := cmd.Process
	m.log.Log("info", "sending SIGTERM to sing-box", map[string]any{"pid": process.Pid})
	if err := process.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
		m.setState(StateError)
		return fmt.Errorf("signal sing-box: %w", err)
	}

	select {
	case err := <-waitCh:
		m.streamWG.Wait()
		if err != nil && !isSignalExit(err) {
			m.setState(StateError)
			return fmt.Errorf("wait for sing-box shutdown: %w", err)
		}
		m.setState(StateStopped)
		m.log.Log("info", "sing-box stopped cleanly", nil)
		return nil
	case <-time.After(m.shutdownWait):
		m.log.Log("warn", "graceful shutdown timed out; sending SIGKILL", map[string]any{"pid": process.Pid})
		if err := process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			m.setState(StateError)
			return fmt.Errorf("kill sing-box: %w", err)
		}
		if err := <-waitCh; err != nil && !isSignalExit(err) {
			m.setState(StateError)
			return fmt.Errorf("wait after kill: %w", err)
		}
		m.streamWG.Wait()
		m.setState(StateStopped)
		return nil
	}
}

func (m *Manager) Restart(ctx context.Context, configPath string) error {
	if m.State() == StateRunning {
		if err := m.Stop(); err != nil {
			return err
		}
	}

	return m.Start(ctx, configPath)
}

func (m *Manager) setState(state State) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.state = state
}

func (m *Manager) streamPipe(out io.Writer, streamName string, pipe io.Reader) {
	defer m.streamWG.Done()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Fprintf(out, "[sing-box/%s] %s\n", streamName, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(out, "[sing-box/%s] stream error: %v\n", streamName, err)
	}
}

func isSignalExit(err error) bool {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}

	return exitErr.ProcessState != nil
}
