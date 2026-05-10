package normalize

import "testing"

func TestHostname(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"api.example.com", "api.example.com"},
		{"https://api.example.com/path", "api.example.com"},
		{"[api.example.com](http://api.example.com)", "api.example.com"},
		{"*.api.example.com", "api.example.com"},
		{"API.EXAMPLE.COM.", "api.example.com"},
		{"  Sub.Domain.com  ", "sub.domain.com"},
	}

	for _, c := range cases {
		res := Hostname(c.in)
		if res != c.out {
			t.Errorf("expected %s but got %s for %s", c.out, res, c.in)
		}
	}
}
