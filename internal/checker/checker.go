package checker

import (
	"os"

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
	StatusHealthy   ResultStatus = "healthy"
	StatusWarning   ResultStatus = "warning"
	StatusCritical  ResultStatus = "critical"
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

	for _, rule := range rules {
		if rule.Package == dep.Path {
			status := StatusWarning
			if rule.Severity == "critical" {
				status = StatusCritical
			}
			return CheckResult{
				Dep:          dep,
				Status:       status,
				Message:      rule.Reason,
				Alternatives: rule.Alternatives,
				ModuleInfo:   info,
			}
		}
	}

	if info == nil {
		return CheckResult{Dep: dep, Status: StatusHealthy}
	}

	if info.Retracted {
		return CheckResult{
			Dep:    dep,
			Status: StatusCritical,
			Message: "Retracted by author",
			ModuleInfo: info,
		}
	}

	if info.Deprecated {
		return CheckResult{
			Dep:    dep,
			Status: StatusCritical,
			Message: "Deprecated",
			ModuleInfo: info,
		}
	}

	behind := info.LatestVersion != "" && info.LatestVersion != "v0.0.0" && info.LatestVersion != dep.Version
	if behind {
		return CheckResult{
			Dep:    dep,
			Status: StatusWarning,
			Message: "Newer version available: " + info.LatestVersion,
			ModuleInfo: info,
		}
	}

	return CheckResult{Dep: dep, Status: StatusHealthy, ModuleInfo: info}
}
