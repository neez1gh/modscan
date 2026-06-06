# modscan

[![Go](https://github.com/neez1gh/modscan/actions/workflows/go.yml/badge.svg)](https://github.com/neez1gh/modscan/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/neez1gh/modscan)](https://goreportcard.com/report/github.com/neez1gh/modscan)
[![License](https://img.shields.io/github/license/neez1gh/modscan)](LICENSE)
[![Release](https://img.shields.io/github/v/release/neez1gh/modscan)](https://github.com/neez1gh/modscan/releases/latest)

Scan your Go dependencies for known vulnerabilities, deprecated packages, and outdated versions.

```
⚠ github.com/sirupsen/logrus v1.9.3
   Low activity, no major updates since 2021
   → go.uber.org/zap — 10x faster, battle-tested at scale
   → github.com/rs/zerolog — Zero allocation, simple API

✗ github.com/dgrijalva/jwt-go v3.2.0
   Officially deprecated by author
   → github.com/golang-jwt/jwt/v5 — Official community successor

✓ github.com/gin-gonic/gin v1.9.1
   Up to date
```

## Features

- **YAML rules** — define known deprecated, vulnerable, or malicious packages with suggested replacements
- **pkg.go.dev integration** — checks latest version, deprecation, and retraction status for every dependency
- **GitHub fallback** — fetches actual commit activity and releases when pkg.go.dev data is stale
- **Parallel scanning** — checks all dependencies concurrently (max 5 requests at a time)
- **CI/CD ready** — `--ci` exits with code 1 on critical issues, `--json` for machine-readable output
- **Color-coded output** — healthy (green), warning (yellow), critical (red)

## Installation

### Quick install

```bash
go install github.com/neez1gh/modscan@latest
```

### Build from source

```bash
git clone https://github.com/neez1gh/modscan.git
cd modscan
go build -o modscan .
```

## Quick start

```bash
# Scan current project
modscan scan

# Exit with code 1 if critical issues found
modscan scan --ci

# JSON output for CI/CD pipelines
modscan scan --json

# Look up why a package is flagged
modscan why github.com/sirupsen/logrus
```

## Usage

### Scan

```bash
# Default scan (uses go.mod in current directory)
modscan scan

# Specify go.mod and rules file paths
modscan scan --go-mod /path/to/go.mod --rules /path/to/rules.yaml

# CI mode — exit code 1 on critical issues
modscan scan --ci

# JSON output
modscan scan --json

# Combine flags
modscan scan --ci --json --go-mod ./submodule/go.mod
```

### Why

```bash
# Show why a package is flagged and its alternatives
modscan why github.com/sirupsen/logrus
```

```
Package: github.com/sirupsen/logrus
Reason:  Low activity, no major updates since 2021
Severity: warning

Alternatives:
  → go.uber.org/zap — 10x faster, battle-tested at scale
  → github.com/rs/zerolog — Zero allocation, simple API
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--ci` | `false` | Exit with code 1 if critical issues found |
| `--json` | `false` | Output as JSON |
| `--rules` | `rules/alternatives.yaml` | Path to YAML rules file |
| `--go-mod` | `go.mod` | Path to go.mod file |

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Scan completed, no critical issues |
| `1` | Critical issues found (with `--ci`) or error |

## Rules

Define custom rules in a YAML file. By default `rules/alternatives.yaml`:

```yaml
rules:
  - package: "github.com/sirupsen/logrus"
    reason: "Low activity, no major updates since 2021"
    severity: "warning"
    alternatives:
      - path: "go.uber.org/zap"
        reason: "10x faster, battle-tested at scale"
      - path: "github.com/rs/zerolog"
        reason: "Zero allocation, simple API"

  - package: "github.com/dgrijalva/jwt-go"
    reason: "Officially deprecated by author"
    severity: "critical"
    alternatives:
      - path: "github.com/golang-jwt/jwt/v5"
        reason: "Official community successor"
```

Severity levels:
- `warning` — displayed but does not fail CI
- `critical` — fails CI when `--ci` is set

## How it works

```
go.mod ──► Parser ──► pkg.go.dev API ──► Checker ──► Reporter
                              │
                              ▼
                        GitHub API
                    (fallback for stale data)
```

1. **Parser** reads `go.mod` using `golang.org/x/mod/modfile`
2. **Fetcher** queries [pkg.go.dev](https://pkg.go.dev) for each dependency's metadata (latest version, commit time, deprecation, retraction) and falls back to GitHub API for stale data
3. **Checker** evaluates each dependency against YAML rules and API data using priority: YAML rule > retracted > deprecated > newer version > stale > healthy
4. **Reporter** displays results with color-coded terminal output or JSON

## CI Integration

### GitHub Actions

```yaml
name: Dependency Audit
on: [push, pull_request]

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Run modscan
        run: |
          go install github.com/neez1gh/modscan@latest
          modscan scan --ci
```

## Contributing

Contributions are welcome! Feel free to:

- Add new rules to [`rules/alternatives.yaml`](rules/alternatives.yaml)
- Report issues or suggest features via GitHub Issues
- Submit pull requests

## License

[MIT](LICENSE)
