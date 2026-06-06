package checker

import (
	"os"
	"time"

	"modscan/internal/fetcher"
	"modscan/internal/parser"

	"gopkg.in/yaml.v3"
)

type Alternative struct {
	Path   string `yaml:"path"`
	Reason string `yaml:"reason"`
}

type Rule struct {
	Package      string        `yaml:"package"`
	Reason       string        `yaml:"reason"`
	Severity     string        `yaml:"severity"`
	Alternatives []Alternative `yaml:"alternatives"`
}

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type ResultStatus string

const (
	StatusHealthy  ResultStatus = "healthy"
	StatusWarning  ResultStatus = "warning"
	StatusCritical ResultStatus = "critical"
)

type CheckResult struct {
	Dep          parser.Dep
	Status       ResultStatus
	Message      string
	Alternatives []Alternative
	ModuleInfo   *fetcher.ModuleInfo
}

func LoadRules(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg.Rules, nil
}

func CheckDep(dep parser.Dep, rules []Rule, info *fetcher.ModuleInfo) CheckResult {
	if dep.Indirect {
		return CheckResult{Dep: dep, Status: StatusHealthy}
	}

	if result := checkAgainstRules(dep, rules, info); result != nil {
		return *result
	}

	if info == nil {
		return CheckResult{Dep: dep, Status: StatusHealthy}
	}

	if info.Retracted {
		return checkResult(dep, StatusCritical, "Retracted by author", info)
	}

	if info.Deprecated {
		return checkResult(dep, StatusCritical, "Deprecated", info)
	}

	latestVersion := resolveLatestVersion(dep, info)
	if latestVersion != "" && latestVersion != "v0.0.0" && latestVersion != dep.Version {
		return checkResult(dep, StatusWarning, "Newer version available: "+latestVersion, info)
	}

	if isStale(dep, info) {
		return checkResult(dep, StatusWarning, "No updates in over a year", info)
	}

	return CheckResult{Dep: dep, Status: StatusHealthy, ModuleInfo: info}
}

func checkAgainstRules(dep parser.Dep, rules []Rule, info *fetcher.ModuleInfo) *CheckResult {
	for _, rule := range rules {
		if rule.Package == dep.Path {
			status := StatusWarning
			if rule.Severity == "critical" {
				status = StatusCritical
			}
			return &CheckResult{
				Dep:          dep,
				Status:       status,
				Message:      rule.Reason,
				Alternatives: rule.Alternatives,
				ModuleInfo:   info,
			}
		}
	}
	return nil
}

func checkResult(dep parser.Dep, status ResultStatus, message string, info *fetcher.ModuleInfo) CheckResult {
	return CheckResult{
		Dep:        dep,
		Status:     status,
		Message:    message,
		ModuleInfo: info,
	}
}

func resolveLatestVersion(dep parser.Dep, info *fetcher.ModuleInfo) string {
	v := info.LatestVersion
	if (v == "" || v == "v0.0.0") && fetcher.IsGitHubModule(dep.Path) {
		if ghVer, err := fetcher.FetchLatestReleaseVersion(dep.Path); err == nil && ghVer != "" {
			return ghVer
		}
	}
	return v
}

func isStale(dep parser.Dep, info *fetcher.ModuleInfo) bool {
	if info.CommitTime == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, info.CommitTime)
	if err != nil {
		return false
	}
	if time.Since(t) <= 365*24*time.Hour {
		return false
	}
	if !fetcher.IsGitHubModule(dep.Path) {
		return true
	}
	ghTime, err := fetcher.FetchLatestCommitTime(dep.Path)
	if err != nil {
		return true
	}
	return time.Since(ghTime) > 365*24*time.Hour
}
