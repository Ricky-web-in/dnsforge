package resolve

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

type Record struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Result struct {
	Hostname string   `json:"hostname"`
	Records  []Record `json:"records"`
	Resolved bool     `json:"resolved"`
	Wildcard bool     `json:"wildcard,omitempty"`
	Error    string   `json:"error,omitempty"`
}

type Resolver struct {
	client  DNSClient
	timeout time.Duration
	retries int
}

type DNSClient interface {
	LookupIP(ctx context.Context, network, host string) ([]net.IPAddr, error)
	LookupCNAME(ctx context.Context, host string) (string, error)
	LookupMX(ctx context.Context, host string) ([]*net.MX, error)
	LookupNS(ctx context.Context, host string) ([]*net.NS, error)
	LookupTXT(ctx context.Context, host string) ([]string, error)
}

type netClient struct {
	resolver *net.Resolver
}

func (n *netClient) LookupIP(ctx context.Context, network, host string) ([]net.IPAddr, error) {
	ips, err := n.resolver.LookupIP(ctx, network, host)
	if err != nil {
		return nil, err
	}
	out := make([]net.IPAddr, 0, len(ips))
	for _, ip := range ips {
		out = append(out, net.IPAddr{IP: ip})
	}
	return out, nil
}
func (n *netClient) LookupCNAME(ctx context.Context, host string) (string, error) {
	return n.resolver.LookupCNAME(ctx, host)
}
func (n *netClient) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	return n.resolver.LookupMX(ctx, host)
}
func (n *netClient) LookupNS(ctx context.Context, host string) ([]*net.NS, error) {
	return n.resolver.LookupNS(ctx, host)
}
func (n *netClient) LookupTXT(ctx context.Context, host string) ([]string, error) {
	return n.resolver.LookupTXT(ctx, host)
}

func NewResolver(customResolver string, timeout time.Duration, retries int) *Resolver {
	r := net.DefaultResolver
	if customResolver != "" {
		if !strings.Contains(customResolver, ":") {
			customResolver += ":53"
		}
		r = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: timeout}
				return d.DialContext(ctx, network, customResolver)
			},
		}
	}

	return &Resolver{
		client:  &netClient{resolver: r},
		timeout: timeout,
		retries: retries,
	}
}

func NewResolverWithClient(client DNSClient, timeout time.Duration, retries int) *Resolver {
	return &Resolver{
		client:  client,
		timeout: timeout,
		retries: retries,
	}
}

func (r *Resolver) Resolve(hostname string, recordTypes []string) Result {
	res := Result{Hostname: hostname}
	resolvedAny := false
	var errs []string

	for _, recType := range recordTypes {
		var err error
		for i := 0; i <= r.retries; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
			err = r.queryRecord(ctx, hostname, recType, &res)
			cancel()
			if err == nil {
				resolvedAny = true
				break
			}
			if i < r.retries {
				time.Sleep(100 * time.Millisecond)
			}
		}
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	res.Resolved = resolvedAny
	if !resolvedAny && len(errs) > 0 {
		res.Error = "resolution failed: " + strings.Join(errs, "; ")
	}

	res.Records = DedupeRecords(res.Records)

	return res
}

func DedupeRecords(records []Record) []Record {
	seen := make(map[string]bool)
	deduped := make([]Record, 0, len(records))
	for _, rec := range records {
		key := rec.Type + ":" + rec.Value
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, rec)
	}
	sort.Slice(deduped, func(i, j int) bool {
		if deduped[i].Type == deduped[j].Type {
			return deduped[i].Value < deduped[j].Value
		}
		return deduped[i].Type < deduped[j].Type
	})
	return deduped
}

func (r *Resolver) queryRecord(ctx context.Context, hostname string, recType string, res *Result) error {
	switch recType {
	case "A":
		ips, err := r.client.LookupIP(ctx, "ip4", hostname)
		if err != nil {
			return err
		}
		for _, ip := range ips {
			res.Records = append(res.Records, Record{Type: "A", Value: ip.IP.String()})
		}
	case "AAAA":
		ips, err := r.client.LookupIP(ctx, "ip6", hostname)
		if err != nil {
			return err
		}
		for _, ip := range ips {
			res.Records = append(res.Records, Record{Type: "AAAA", Value: ip.IP.String()})
		}
	case "CNAME":
		cname, err := r.client.LookupCNAME(ctx, hostname)
		if err != nil {
			return err
		}
		if cname != "" && cname != hostname+"." {
			res.Records = append(res.Records, Record{Type: "CNAME", Value: strings.TrimSuffix(cname, ".")})
		}
	case "MX":
		mxs, err := r.client.LookupMX(ctx, hostname)
		if err != nil {
			return err
		}
		for _, mx := range mxs {
			res.Records = append(res.Records, Record{Type: "MX", Value: strings.TrimSuffix(mx.Host, ".")})
		}
	case "NS":
		nss, err := r.client.LookupNS(ctx, hostname)
		if err != nil {
			return err
		}
		for _, ns := range nss {
			res.Records = append(res.Records, Record{Type: "NS", Value: strings.TrimSuffix(ns.Host, ".")})
		}
	case "TXT":
		txts, err := r.client.LookupTXT(ctx, hostname)
		if err != nil {
			return err
		}
		for _, txt := range txts {
			res.Records = append(res.Records, Record{Type: "TXT", Value: txt})
		}
	default:
		return fmt.Errorf("unsupported record type %s", recType)
	}
	return nil
}
