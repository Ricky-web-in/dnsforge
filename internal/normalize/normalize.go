package normalize

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	mdRegex = regexp.MustCompile(`\[.*?\]\((.*?)\)`)
)

func Hostname(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	if matches := mdRegex.FindStringSubmatch(input); len(matches) > 1 {
		input = matches[1]
	}

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		u, err := url.Parse(input)
		if err == nil && u.Hostname() != "" {
			input = u.Hostname()
		}
	}

	input = strings.TrimPrefix(input, "*.")
	input = strings.TrimSuffix(input, ".")

	return strings.ToLower(input)
}
