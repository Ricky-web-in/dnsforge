package wildcard

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"dnsforge/internal/resolve"
)

func RandomLabel(root string) (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", hex.EncodeToString(b), strings.TrimSpace(root)), nil
}

func Signature(records []resolve.Record) map[string]bool {
	sig := make(map[string]bool)
	for _, rec := range records {
		sig[rec.Type+":"+rec.Value] = true
	}
	return sig
}

func IsWildcardMatch(result resolve.Result, wildcardSig map[string]bool) bool {
	if !result.Resolved {
		return false
	}
	if len(wildcardSig) == 0 {
		return false
	}
	for _, rec := range result.Records {
		if wildcardSig[rec.Type+":"+rec.Value] {
			return true
		}
	}
	return false
}
