package runner

import (
	"sync"
	"time"

	"dnsforge/internal/config"
	"dnsforge/internal/resolve"
	"dnsforge/internal/wildcard"
)

func Run(cfg *config.Config, hosts []string) ([]resolve.Result, error) {
	resolver := resolve.NewResolver(cfg.Resolver, cfg.Timeout, cfg.Retries)
	return RunWithResolver(cfg, hosts, resolver)
}

func RunWithResolver(cfg *config.Config, hosts []string, resolver *resolve.Resolver) ([]resolve.Result, error) {
	var wildcardSig map[string]bool
	if cfg.DetectWildcard {
		testHost, err := wildcard.RandomLabel(cfg.Root)
		if err != nil {
			return nil, err
		}
		wildRes := resolver.Resolve(testHost, cfg.Records)
		wildcardSig = wildcard.Signature(wildRes.Records)
	}

	jobs := make(chan string)
	results := make(chan resolve.Result)
	var wg sync.WaitGroup
	rateTicker := newRateTicker(cfg.RateLimit)
	if rateTicker != nil {
		defer rateTicker.Stop()
	}

	worker := func() {
		defer wg.Done()
		for host := range jobs {
			if rateTicker != nil {
				<-rateTicker.C
			}
			result := resolver.Resolve(host, cfg.Records)
			if cfg.DetectWildcard {
				result.Wildcard = wildcard.IsWildcardMatch(result, wildcardSig)
			}
			results <- result
		}
	}

	wg.Add(cfg.Concurrency)
	for i := 0; i < cfg.Concurrency; i++ {
		go worker()
	}

	go func() {
		for _, host := range hosts {
			jobs <- host
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	out := make([]resolve.Result, 0, len(hosts))
	for result := range results {
		out = append(out, result)
	}
	return out, nil
}

func newRateTicker(rate int) *time.Ticker {
	if rate <= 0 {
		return nil
	}
	return time.NewTicker(time.Second / time.Duration(rate))
}
