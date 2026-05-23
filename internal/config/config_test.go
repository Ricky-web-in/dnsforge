package config

import (
	"bytes"
	"testing"
	"time"
)

func TestParseRecordTypes(t *testing.T) {
	got, err := ParseRecordTypes("a,aaaa,cname, A ,TXT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"A", "AAAA", "CNAME", "TXT"}
	if len(got) != len(want) {
		t.Fatalf("unexpected length: got=%d want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected value at %d: got=%s want=%s", i, got[i], want[i])
		}
	}
}

func TestParseRecordTypesRejectsInvalid(t *testing.T) {
	if _, err := ParseRecordTypes("A,BOGUS"); err == nil {
		t.Fatal("expected invalid record type error")
	}
}

func TestParseDefaultsMatchSpec(t *testing.T) {
	cfg, err := Parse([]string{"-d", "example.com"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if cfg.Concurrency != 50 {
		t.Fatalf("unexpected default concurrency: %d", cfg.Concurrency)
	}
	if cfg.Timeout != 5*time.Second {
		t.Fatalf("unexpected default timeout: %s", cfg.Timeout)
	}
	if cfg.Retries != 1 {
		t.Fatalf("unexpected default retries: %d", cfg.Retries)
	}
	want := []string{"A", "AAAA", "CNAME"}
	for i := range want {
		if cfg.Records[i] != want[i] {
			t.Fatalf("unexpected default records: got=%v", cfg.Records)
		}
	}
}

func TestParseAliases(t *testing.T) {
	cfg, err := Parse([]string{"-d", "example.com", "-o", "out.jsonl", "-c", "77"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if cfg.OutputFile != "out.jsonl" {
		t.Fatalf("expected -o alias to set output file, got %q", cfg.OutputFile)
	}
	if cfg.Concurrency != 77 {
		t.Fatalf("expected -c alias to set concurrency, got %d", cfg.Concurrency)
	}
}
