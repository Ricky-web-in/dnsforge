package resolve

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"
)

type fakeDNSClient struct {
	ip4   []net.IPAddr
	ip6   []net.IPAddr
	cname string
	txt   []string
	err   error
}

func (f *fakeDNSClient) LookupIP(_ context.Context, network, _ string) ([]net.IPAddr, error) {
	if f.err != nil {
		return nil, f.err
	}
	if network == "ip4" {
		return f.ip4, nil
	}
	return f.ip6, nil
}
func (f *fakeDNSClient) LookupCNAME(_ context.Context, _ string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.cname, nil
}
func (f *fakeDNSClient) LookupMX(_ context.Context, _ string) ([]*net.MX, error) {
	return nil, f.err
}
func (f *fakeDNSClient) LookupNS(_ context.Context, _ string) ([]*net.NS, error) {
	return nil, f.err
}
func (f *fakeDNSClient) LookupTXT(_ context.Context, _ string) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.txt, nil
}

func TestResolveSuccess(t *testing.T) {
	client := &fakeDNSClient{
		ip4: []net.IPAddr{{IP: net.ParseIP("1.1.1.1")}, {IP: net.ParseIP("1.1.1.1")}},
		txt: []string{"v=spf1 include:_spf.example.com ~all"},
	}
	resolver := NewResolverWithClient(client, 2*time.Second, 0)
	res := resolver.Resolve("example.com", []string{"A", "TXT"})
	if !res.Resolved {
		t.Fatal("expected resolved result")
	}
	if len(res.Records) != 2 {
		t.Fatalf("expected 2 deduped records, got %d", len(res.Records))
	}
}

func TestResolveUnresolvedHandling(t *testing.T) {
	resolver := NewResolverWithClient(&fakeDNSClient{err: errors.New("nxdomain")}, 2*time.Second, 0)
	res := resolver.Resolve("missing.example.com", []string{"A"})
	if res.Resolved {
		t.Fatal("expected unresolved result")
	}
	if !strings.Contains(res.Error, "resolution failed") {
		t.Fatalf("unexpected unresolved error: %s", res.Error)
	}
}

func TestDedupeRecords(t *testing.T) {
	in := []Record{
		{Type: "A", Value: "1.1.1.1"},
		{Type: "A", Value: "1.1.1.1"},
		{Type: "TXT", Value: "hello"},
		{Type: "TXT", Value: "hello"},
	}
	out := DedupeRecords(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 deduped records, got %d", len(out))
	}
}
