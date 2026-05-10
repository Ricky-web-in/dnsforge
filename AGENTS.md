# Agent Instructions for dnsforge

You are building dnsforge, a safe DNS resolution and validation tool for  reconnaissance workflows.

Rules:

- Only resolve user-provided hostnames and optional wildcard test hostnames.
- Add timeout, retries, concurrency, and optional rate limit.
- Keep output structured and automation-friendly.
- Prefer JSONL for bulk output.
- Keep code modular.
- Add tests for normalization, parsing, output, resolution result handling, and wildcard logic.
