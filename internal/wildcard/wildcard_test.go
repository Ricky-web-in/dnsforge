package wildcard

import (
	"strings"
	"testing"

	"github.com/Ricky-web-in/dnsforge/internal/resolve"
)

func TestRandomLabel(t *testing.T) {
	label, err := RandomLabel("example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(label, ".example.com") {
		t.Fatalf("unexpected wildcard label: %s", label)
	}
}

func TestWildcardComparison(t *testing.T) {
	wildSig := Signature([]resolve.Record{{Type: "A", Value: "1.1.1.1"}})
	result := resolve.Result{
		Hostname: "api.example.com",
		Resolved: true,
		Records:  []resolve.Record{{Type: "A", Value: "1.1.1.1"}},
	}
	if !IsWildcardMatch(result, wildSig) {
		t.Fatal("expected wildcard match")
	}
}
