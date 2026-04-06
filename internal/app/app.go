package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quadraphony/ghostfy/internal/adapters/singbox"
	"github.com/quadraphony/ghostfy/internal/config"
	"github.com/quadraphony/ghostfy/internal/importers"
	"github.com/quadraphony/ghostfy/internal/logging"
	"github.com/quadraphony/ghostfy/internal/observability"
	"github.com/quadraphony/ghostfy/internal/runtime"
)

type App struct {
	mapper *singbox.Mapper
}

func New() *App {
	return &App{
		mapper: singbox.NewMapper(),
	}
}

func (a *App) Run(configPath string, stdout, stderr *os.File) error {
	log := logging.New(stdout)
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	log.Log("info", "loaded ghostify config", map[string]any{
		"profile": cfg.Profile,
		"path":    configPath,
	})

	singboxConfig, err := a.mapper.Build(cfg)
	if err != nil {
		return err
	}

	runtimeConfigPath, err := writeRuntimeConfig(singboxConfig)
	if err != nil {
		return err
	}
	defer os.Remove(runtimeConfigPath)

	log.Log("info", "generated sing-box runtime config", map[string]any{"path": runtimeConfigPath})

	manager := runtime.NewManager(stdout, stderr, log)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := manager.Start(ctx, runtimeConfigPath); err != nil {
		observability.Default.Record("run-error")
		return err
	}

	log.Log("info", "ghostify runtime running", map[string]any{
		"state":          manager.State(),
		"inbound_type":   cfg.Inbound.Type,
		"inbound_listen": fmt.Sprintf("%s:%d", cfg.Inbound.Listen, cfg.Inbound.Port),
		"outbound_type":  cfg.Outbound.Type,
	})

	<-ctx.Done()
	log.Log("info", "shutdown signal received", map[string]any{"signal": ctx.Err().Error()})
	observability.Default.Record("run-success")
	return manager.Stop()
}

func (a *App) Render(configPath string, out io.Writer) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	singboxConfig, err := a.mapper.Build(cfg)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(singboxConfig)
}

func (a *App) ImportURI(raw string, out io.Writer) error {
	cfg, err := importers.Import(raw)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

func (a *App) ListProtocols() []map[string]string {
	entries := a.mapper.Registry().List()
	out := make([]map[string]string, 0, len(entries))
	for _, entry := range entries {
		out = append(out, map[string]string{
			"name":          entry.Name,
			"category":      entry.Category,
			"status":        entry.Status,
			"support_class": entry.SupportClass,
		})
	}
	return out
}

func (a *App) HealthCheck(configPath string, out io.Writer) error {
	log := logging.New(out)
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	singboxConfig, err := a.mapper.Build(cfg)
	if err != nil {
		return err
	}

	runtimeConfigPath, err := writeRuntimeConfig(singboxConfig)
	if err != nil {
		return err
	}
	defer os.Remove(runtimeConfigPath)

	manager := runtime.NewManager(out, out, log)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := manager.Start(ctx, runtimeConfigPath); err != nil {
		observability.Default.Record("health-error")
		return err
	}

	observability.Default.Record("health-success")
	return manager.Stop()
}

func writeRuntimeConfig(data map[string]any) (string, error) {
	file, err := os.CreateTemp("", "ghostify-runtime-*.json")
	if err != nil {
		return "", fmt.Errorf("create runtime config: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return "", fmt.Errorf("encode runtime config: %w", err)
	}

	return file.Name(), nil
}
