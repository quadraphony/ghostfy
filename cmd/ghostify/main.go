package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/quadraphony/ghostfy/internal/app"
	"github.com/quadraphony/ghostfy/internal/bridge"
	"github.com/quadraphony/ghostfy/internal/bridge/openvpn"
	"github.com/quadraphony/ghostfy/internal/bridge/ssh"
	"github.com/quadraphony/ghostfy/internal/bridge/udp2raw"
	"github.com/quadraphony/ghostfy/internal/observability"
	"github.com/quadraphony/ghostfy/internal/runtime"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ghostify: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("expected a command; available: run-singbox-test, run, render, import-uri, health, status, status-server, openvpn-bridge, ssh-bridge, udp2raw-bridge, bridges, protocols")
	}

	switch args[0] {
	case "run-singbox-test":
		if len(args) > 1 {
			return errors.New("run-singbox-test does not accept extra arguments")
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		runner := runtime.NewSingboxTestRunner(os.Stdout, os.Stderr)
		return runner.Run(ctx)
	case "run":
		if len(args) != 3 || args[1] != "-c" {
			return errors.New(`usage: ghostify run -c <ghostify.json>`)
		}
		return app.New().Run(args[2], os.Stdout, os.Stderr)
	case "render":
		if len(args) != 3 || args[1] != "-c" {
			return errors.New(`usage: ghostify render -c <ghostify.json>`)
		}
		return app.New().Render(args[2], os.Stdout)
	case "import-uri":
		if len(args) != 2 {
			return errors.New(`usage: ghostify import-uri '<uri>'`)
		}
		return app.New().ImportURI(args[1], os.Stdout)
	case "protocols":
		if len(args) != 1 {
			return errors.New("usage: ghostify protocols")
		}
		return printProtocols()
	case "health":
		if len(args) != 3 || args[1] != "-c" {
			return errors.New(`usage: ghostify health -c <ghostify.json>`)
		}
		return app.New().HealthCheck(args[2], os.Stdout)
	case "status":
		if len(args) != 1 {
			return errors.New("usage: ghostify status")
		}
		status := map[string]any{
			"metrics":   observability.Default.Snapshot(),
			"protocols": app.New().ListProtocols(),
			"bridges":   bridge.Default().List(),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(status)
	case "openvpn-bridge":
		if len(args) != 3 || args[1] != "-c" {
			return errors.New(`usage: ghostify openvpn-bridge -c <config.json>`)
		}
		cfg, err := openvpn.Load(args[2])
		if err != nil {
			return err
		}

		runner := openvpn.NewRunner(os.Stdout, os.Stderr)
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		return runner.Run(ctx, cfg)
	case "ssh-bridge":
		if len(args) != 3 || args[1] != "-c" {
			return errors.New(`usage: ghostify ssh-bridge -c <config.json>`)
		}
		cfg, err := ssh.Load(args[2])
		if err != nil {
			return err
		}
		runner := ssh.NewRunner(os.Stdout, os.Stderr)
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		return runner.Run(ctx, cfg)
	case "udp2raw-bridge":
		if len(args) != 3 || args[1] != "-c" {
			return errors.New(`usage: ghostify udp2raw-bridge -c <config.json>`)
		}
		cfg, err := udp2raw.Load(args[2])
		if err != nil {
			return err
		}
		runner := udp2raw.NewRunner(os.Stdout, os.Stderr)
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		return runner.Run(ctx, cfg)
	case "bridges":
		if len(args) != 1 {
			return errors.New("usage: ghostify bridges")
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(bridge.Default().List())
	case "status-server":
		fs := flag.NewFlagSet("status-server", flag.ExitOnError)
		addr := fs.String("addr", "127.0.0.1:9111", "listen address")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		server := &observability.Server{Addr: *addr}
		return server.ListenAndServe()
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printProtocols() error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(app.New().ListProtocols())
}
