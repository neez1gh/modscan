package cmd

import (
	"fmt"
	"os"
	"sync"

	"modscan/internal/checker"
	"modscan/internal/fetcher"
	"modscan/internal/parser"
	"modscan/internal/reporter"

	"github.com/spf13/cobra"
)

var (
	ciMode     bool
	outputJSON bool
	rulesPath  string
	gomodPath  string
)

var rootCmd = &cobra.Command{
	Use:   "modscan",
	Short: "Audit your Go dependencies",
	Long:  `modscan scans your go.mod file and checks dependencies against known issues, vulnerabilities, and outdated packages.`,
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan go.mod for dependency issues",
	Run: func(cmd *cobra.Command, args []string) {
		deps, err := parser.ParseGoMod(gomodPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing go.mod: %v\n", err)
			os.Exit(1)
		}

		rules, err := checker.LoadRules(rulesPath)
		if err != nil {
			if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Warning: could not load rules: %v\n", err)
			}
			rules = nil
		}

		var wg sync.WaitGroup
		mu := sync.Mutex{}
		results := make([]checker.CheckResult, 0, len(deps))

		sem := make(chan struct{}, 5)

		for _, dep := range deps {
			wg.Add(1)
			sem <- struct{}{}
			go func(d parser.Dep) {
				defer wg.Done()
				defer func() { <-sem }()

				info, err := fetcher.FetchModuleInfo(d.Path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not fetch info for %s: %v\n", d.Path, err)
				}

				result := checker.CheckDep(d, rules, info)

				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}(dep)
		}

		wg.Wait()

		if outputJSON {
			reporter.PrintJSON(results)
		} else {
			reporter.PrintTerm(results)
		}

		if ciMode && reporter.HasCritical(results) {
			os.Exit(1)
		}
	},
}

func init() {
	scanCmd.Flags().BoolVar(&ciMode, "ci", false, "Exit with code 1 if critical issues found")
	scanCmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	scanCmd.Flags().StringVar(&rulesPath, "rules", "rules/alternatives.yaml", "Path to rules file")
	scanCmd.Flags().StringVar(&gomodPath, "go-mod", "go.mod", "Path to go.mod file")

	rootCmd.AddCommand(scanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
