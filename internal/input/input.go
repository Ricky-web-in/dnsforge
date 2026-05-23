package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Ricky-web-in/dnsforge/internal/config"
	"github.com/Ricky-web-in/dnsforge/internal/normalize"
)

func Collect(cfg *config.Config, stdin io.Reader) ([]string, error) {
	var raw []string
	if strings.TrimSpace(cfg.Domain) != "" {
		raw = append(raw, cfg.Domain)
	}

	if strings.TrimSpace(cfg.InputFile) != "" {
		lines, err := readFileLines(cfg.InputFile)
		if err != nil {
			return nil, err
		}
		raw = append(raw, lines...)
	}

	if hasPipedInput() {
		lines, err := ReadLines(stdin)
		if err != nil {
			return nil, err
		}
		raw = append(raw, lines...)
	}

	hosts := NormalizeAndDedupe(raw)
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no valid hostnames provided")
	}
	return hosts, nil
}

func ReadLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	out := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out = append(out, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func NormalizeAndDedupe(lines []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		host := normalize.Hostname(line)
		if host == "" {
			continue
		}
		if !seen[host] {
			seen[host] = true
			out = append(out, host)
		}
	}
	return out
}

func readFileLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadLines(f)
}

func hasPipedInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice == 0
}
