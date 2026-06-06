# modscan

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
```

## Install

```bash
go install github.com/neez1gh/modscan@latest
```

Or build from source:

```bash
git clone https://github.com/neez1gh/modscan.git
cd modscan
go build -o modscan .
```

## Usage

```bash
# Scan current project
modscan scan

# Scan with custom paths
modscan scan --go-mod /path/to/go.mod --rules /path/to/rules.yaml

# JSON output for CI/CD
modscan scan --json

# Exit with code 1 on critical issues
modscan scan --ci

# Look up why a package is flagged
modscan why github.com/sirupsen/logrus
```

## How it works

1. Parses `go.mod` using `golang.org/x/mod/modfile`
2. Fetches module metadata from [pkg.go.dev API](https://pkg.go.dev)
3. Checks against YAML rules for known problematic packages
4. Falls back to GitHub API for accurate staleness data
5. Reports results with color-coded terminal output or JSON

### Data sources

| Source | What it provides |
|---|---|
| [pkg.go.dev](https://pkg.go.dev) | Latest version, commit time, deprecation and retraction flags |
| [GitHub API](https://docs.github.com/en/rest) | Actual latest commit time, latest release tag (fallback) |
| YAML rules | Curated list of known deprecated/malicious packages + alternatives |

## Rules

Add custom rules in `rules/alternatives.yaml`:

```yaml
rules:
  - package: "github.com/sirupsen/logrus"
    reason: "Low activity, no major updates since 2021"
    severity: "warning"
    alternatives:
      - path: "go.uber.org/zap"
        reason: "10x faster, battle-tested at scale"
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--ci` | `false` | Exit code 1 if critical issues found |
| `--json` | `false` | Output as JSON |
| `--rules` | `rules/alternatives.yaml` | Path to rules file |
| `--go-mod` | `go.mod` | Path to go.mod file |

## License

MIT
