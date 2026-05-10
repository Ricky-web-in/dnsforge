package input

import "testing"

func TestNormalizeAndDedupeMarkdownInput(t *testing.T) {
	in := []string{
		"[api.example.com](https://api.example.com/path)",
		"*.api.example.com",
		"API.EXAMPLE.COM.",
	}
	got := NormalizeAndDedupe(in)
	if len(got) != 1 {
		t.Fatalf("expected 1 unique host, got %d", len(got))
	}
	if got[0] != "api.example.com" {
		t.Fatalf("unexpected hostname: %s", got[0])
	}
}
