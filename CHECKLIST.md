# dnsforge Checklist

## Project

- [ ] SPEC.md exists
- [ ] AGENTS.md exists
- [ ] README.md exists
- [ ] .gitignore exists
- [ ] Tests exist

## CLI

- [ ] dnsforge --help works
- [ ] dnsforge --version works
- [ ] dnsforge -d example.com works
- [ ] dnsforge -i hosts.txt works
- [ ] cat hosts.txt | dnsforge works
- [ ] --format jsonl works
- [ ] --format txt works
- [ ] --output works
- [ ] --records works
- [ ] --timeout works
- [ ] --retries works
- [ ] --concurrency works
- [ ] --rate-limit works
- [ ] --resolver works
- [ ] --detect-wildcard --root example.com works

## Normalization

- [ ] api.example.com normalizes correctly
- [ ] https://api.example.com/path normalizes correctly
- [ ] [api.example.com](http://api.example.com) normalizes correctly
- [ ] *.api.example.com normalizes correctly
- [ ] API.EXAMPLE.COM. normalizes correctly

## DNS

- [ ] A resolution works
- [ ] AAAA resolution works
- [ ] CNAME resolution works
- [ ] MX/NS/TXT supported when requested
- [ ] unresolved hosts handled cleanly
- [ ] retries are controlled
- [ ] timeout is enforced

## Output

- [ ] JSONL records are valid JSON
- [ ] TXT output contains only resolved hostnames
- [ ] output file writing works
- [ ] errors appear in JSONL when appropriate

## Wildcard

- [ ] random labels generated under root
- [ ] wildcard records collected
- [ ] wildcard comparison works
- [ ] results are marked wildcard when appropriate


## Quality

- [ ] go test ./... passes
- [ ] go vet ./... passes
- [ ] README has examples
- [ ] README has safety note
