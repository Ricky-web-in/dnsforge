# dnsforge SPEC

## Goal

Build dnsforge, a safe DNS resolution and validation tool for  reconnaissance workflows.

The tool reads domains/subdomains, normalizes them, resolves DNS records, detects wildcard DNS when requested, and outputs structured results.


---

## Milestone 1: MVP

Implement only Milestone 1 first.

---

## Input Modes

### Single target

```bash
dnsforge -d api.example.com
