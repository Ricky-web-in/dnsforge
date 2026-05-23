package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
)

const Version = "v0.1.0-milestone1"

type Config struct {
	Domain            string
	InputFile         string
	Format            string
	OutputFile        string
	Records           []string
	Timeout           time.Duration
	Retries           int
	Concurrency       int
	RateLimit         int
	Resolver          string
	IncludeUnresolved bool
	Silent            bool
	Verbose           bool
	DetectWildcard    bool
	Root              string
	ShowHelp          bool
	ShowVersion       bool
}

func Parse(args []string, stderr io.Writer) (*Config, error) {
	cfg := &Config{}
	fs := flag.NewFlagSet("dnsforge", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}

	fs.StringVar(&cfg.Domain, "d", "", "Single domain to resolve")
	fs.StringVar(&cfg.InputFile, "i", "", "Input file containing hostnames")
	fs.StringVar(&cfg.Format, "format", "jsonl", "Output format: jsonl or txt")
	fs.StringVar(&cfg.OutputFile, "o", "", "Output file path")
	fs.StringVar(&cfg.OutputFile, "output", "", "Output file path (optional)")
	recordsStr := fs.String("records", "A,AAAA,CNAME", "Comma-separated record types (A,AAAA,CNAME,MX,NS,TXT)")
	fs.DurationVar(&cfg.Timeout, "timeout", 5*time.Second, "DNS resolution timeout")
	fs.IntVar(&cfg.Retries, "retries", 1, "Number of retries per record type")
	fs.IntVar(&cfg.Concurrency, "c", 50, "Number of concurrent workers")
	fs.IntVar(&cfg.Concurrency, "concurrency", 50, "Number of concurrent workers")
	fs.IntVar(&cfg.RateLimit, "rate-limit", 0, "Rate limit in requests/second (0 for unlimited)")
	fs.StringVar(&cfg.Resolver, "resolver", "", "Custom resolver IP:PORT (e.g. 8.8.8.8:53)")
	fs.BoolVar(&cfg.IncludeUnresolved, "include-unresolved", true, "Include unresolved hosts in JSONL output")
	fs.BoolVar(&cfg.Silent, "silent", false, "Suppress non-essential stderr logs")
	fs.BoolVar(&cfg.Verbose, "v", false, "Enable verbose logs")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logs")
	fs.BoolVar(&cfg.DetectWildcard, "detect-wildcard", false, "Enable wildcard detection")
	fs.StringVar(&cfg.Root, "root", "", "Root domain for wildcard detection")
	fs.BoolVar(&cfg.ShowVersion, "version", false, "Show version")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			cfg.ShowHelp = true
			return cfg, nil
		}
		return nil, err
	}

	records, err := ParseRecordTypes(*recordsStr)
	if err != nil {
		return nil, err
	}
	cfg.Records = records

	if cfg.Format != "jsonl" && cfg.Format != "txt" {
		return nil, fmt.Errorf("unsupported format %q", cfg.Format)
	}
	if cfg.Concurrency <= 0 {
		return nil, fmt.Errorf("concurrency must be > 0")
	}
	if cfg.Retries < 0 {
		return nil, fmt.Errorf("retries must be >= 0")
	}
	if cfg.Timeout <= 0 {
		return nil, fmt.Errorf("timeout must be > 0")
	}
	if cfg.RateLimit < 0 {
		return nil, fmt.Errorf("rate-limit must be >= 0")
	}
	if cfg.DetectWildcard && strings.TrimSpace(cfg.Root) == "" {
		return nil, fmt.Errorf("--root is required with --detect-wildcard")
	}
	return cfg, nil
}

func ParseRecordTypes(input string) ([]string, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("records cannot be empty")
	}
	allowed := map[string]bool{
		"A": true, "AAAA": true, "CNAME": true, "MX": true, "NS": true, "TXT": true,
	}
	var out []string
	seen := make(map[string]bool)
	for _, part := range strings.Split(input, ",") {
		record := strings.ToUpper(strings.TrimSpace(part))
		if record == "" {
			continue
		}
		if !allowed[record] {
			return nil, fmt.Errorf("unsupported record type %q", record)
		}
		if !seen[record] {
			seen[record] = true
			out = append(out, record)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("records cannot be empty")
	}
	return out, nil
}

func PrintHelp(stdout io.Writer) {
	fmt.Fprintln(stdout, "dnsforge - safe DNS resolution and validation tool")
	fmt.Fprintln(stdout, "Usage: dnsforge [options]")
	fmt.Fprintln(stdout, "Options:")
	fmt.Fprintln(stdout, "  -d string              Single domain to resolve")
	fmt.Fprintln(stdout, "  -i string              Input file containing hostnames")
	fmt.Fprintln(stdout, "  -format string         Output format: jsonl or txt (default \"jsonl\")")
	fmt.Fprintln(stdout, "  -o, --output string    Output file path")
	fmt.Fprintln(stdout, "  -records string        Comma-separated record types (default \"A,AAAA,CNAME\")")
	fmt.Fprintln(stdout, "  -timeout duration      DNS resolution timeout (default 5s)")
	fmt.Fprintln(stdout, "  -retries int           Number of retries per record type (default 1)")
	fmt.Fprintln(stdout, "  -c, --concurrency int  Number of concurrent workers (default 50)")
	fmt.Fprintln(stdout, "  -rate-limit int        Requests per second, 0 means unlimited")
	fmt.Fprintln(stdout, "  -resolver string       Custom resolver IP:PORT (default \"system\")")
	fmt.Fprintln(stdout, "  --include-unresolved   Include unresolved hosts in JSONL output (default true)")
	fmt.Fprintln(stdout, "  --silent               Suppress non-essential stderr logs")
	fmt.Fprintln(stdout, "  -v, --verbose          Enable verbose logs")
	fmt.Fprintln(stdout, "  -detect-wildcard       Enable wildcard detection")
	fmt.Fprintln(stdout, "  -root string           Root domain for wildcard detection")
	fmt.Fprintln(stdout, "  -version               Show version")
}

func PrintVersion(stdout io.Writer) {
	fmt.Fprintf(stdout, "dnsforge %s\n", Version)
}
