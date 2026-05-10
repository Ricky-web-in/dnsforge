You are implementing a professional cybersecurity reconnaissance tool .

Read these files first:

- SPEC.md
- AGENTS.md
- CHECKLIST.md

Task:

Implement Milestone 1 of dnsforge.

Important:

- Implement only what is listed in SPEC.md.
- Use Go.
- Prefer the Go standard library where practical.
- Keep code modular.
- Add tests for hostname normalization, Markdown input normalization, record type parsing, deduplication, JSONL output, TXT output, unresolved handling, and wildcard comparison logic.
- Add README.md with clear explanation, install instructions, usage examples, output examples, and safety note.
- Ensure CLI works as dnsforge.
- Run go test ./... and fix failures.
- Run go vet ./... and fix issues.

Expected structure:

- cmd/dnsforge/main.go
- internal/config/
- internal/input/
- internal/normalize/
- internal/resolve/
- internal/output/
- internal/wildcard/
- internal/runner/
- examples/
- README.md

At the end, summarize:

1. Files created/changed
2. How to build
3. How to run
4. How to test
