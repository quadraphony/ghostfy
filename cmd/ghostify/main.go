package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/quadraphony/ghostfy/internal/app"
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
		return errors.New("expected a command; available: run-singbox-test, run, render, import-uri, protocols")
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
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printProtocols() error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(app.New().ListProtocols())
}
