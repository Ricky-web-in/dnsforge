package config

import "testing"

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
