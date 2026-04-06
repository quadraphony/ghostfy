package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/quadraphony/ghostfy/internal/logging"
)

type SingboxTestRunner struct {
	stdout      io.Writer
	stderr      io.Writer
	tempDir     string
	startupWait time.Duration
	manager     *Manager
}

func NewSingboxTestRunner(stdout, stderr io.Writer) *SingboxTestRunner {
	log := logging.New(stdout)
	return &SingboxTestRunner{
		stdout:      stdout,
		stderr:      stderr,
		tempDir:     os.TempDir(),
		startupWait: defaultStartupWait,
		manager:     NewManager(stdout, stderr, log),
	}
}

func (r *SingboxTestRunner) Run(ctx context.Context) error {
	configPath, err := r.writeConfig()
	if err != nil {
		return fmt.Errorf("write temporary config: %w", err)
	}
	defer os.Remove(configPath)

	r.manager.log.Log("info", "generated sing-box test config", map[string]any{"path": configPath})
	r.manager.startupWait = r.startupWait

	if err := r.manager.Start(ctx, configPath); err != nil {
		return err
	}

	if r.manager.State() == StateStopped {
		return nil
	}

	r.manager.log.Log("info", "startup window passed; stopping sing-box cleanly", map[string]any{"wait": r.startupWait.String()})
	return r.manager.Stop()
}

func (r *SingboxTestRunner) writeConfig() (string, error) {
	if err := os.MkdirAll(r.tempDir, 0o755); err != nil {
		return "", err
	}

	file, err := os.CreateTemp(r.tempDir, "ghostify-singbox-*.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(minimalConfig()); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func minimalConfig() map[string]any {
	return map[string]any{
		"log": map[string]any{
			"level":     "info",
			"timestamp": true,
		},
		"outbounds": []map[string]any{
			{
				"type": "direct",
			},
		},
	}
}

func configFileName(path string) string {
	return filepath.Base(path)
}
