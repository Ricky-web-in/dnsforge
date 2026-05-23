package main

import (
	"fmt"
	"os"

	"github.com/Ricky-web-in/dnsforge/internal/config"
	"github.com/Ricky-web-in/dnsforge/internal/input"
	"github.com/Ricky-web-in/dnsforge/internal/output"
	"github.com/Ricky-web-in/dnsforge/internal/runner"
)

func main() {
	cfg, err := config.Parse(os.Args[1:], os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if cfg.ShowHelp {
		config.PrintHelp(os.Stdout)
		return
	}
	if cfg.ShowVersion {
		config.PrintVersion(os.Stdout)
		return
	}

	hosts, err := input.Collect(cfg, os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	results, err := runner.Run(cfg, hosts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	w, err := output.WriterForPath(cfg.OutputFile, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer w.Close()

	if err := output.WriteResults(results, cfg.Format, cfg.IncludeUnresolved, cfg.Records, w); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
