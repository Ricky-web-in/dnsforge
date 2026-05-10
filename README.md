# dnsforge

`dnsforge` is a safe DNS resolution and validation CLI for reconnaissance workflows.

It reads hostnames from a single target flag, file input, or stdin, normalizes them, resolves selected DNS records, optionally detects wildcard DNS behavior, and outputs either JSONL or TXT.

## Safety

- `dnsforge` only resolves hostnames you provide.
- Wildcard checks only query randomized labels under the root domain you explicitly set.

## Install

```bash
go build -o dnsforge ./cmd/dnsforge
```

## Usage

```bash
dnsforge --help
dnsforge --version
dnsforge -d api.example.com
dnsforge -i examples/hosts.txt --records A,AAAA,CNAME --format jsonl
cat examples/hosts.txt | dnsforge --format txt
dnsforge -i examples/hosts.txt --output results.jsonl
dnsforge -d api.example.com --resolver 8.8.8.8:53 --timeout 3s --retries 1
dnsforge -i examples/hosts.txt --detect-wildcard --root example.com
```

## Output examples

JSONL:

```json
{"hostname":"api.example.com","records":[{"type":"A","value":"93.184.216.34"}],"resolved":true}
{"hostname":"missing.example.com","records":[],"resolved":false,"error":"resolution failed: lookup missing.example.com: no such host"}
```

TXT:

```txt
api.example.com
www.example.com
```

Only resolved hostnames appear in TXT output.
