package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"modscan/internal/checker"

	"github.com/charmbracelet/lipgloss"
)

var (
	okStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)
	critStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	dimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	boldStyle = lipgloss.NewStyle().Bold(true)
)

func PrintTerm(results []checker.CheckResult) {
	summary := map[checker.ResultStatus]int{
		checker.StatusHealthy:  0,
		checker.StatusWarning:  0,
		checker.StatusCritical: 0,
	}

	fmt.Println(boldStyle.Render("\n📦 Scanning dependencies...\n"))

	for _, r := range results {
		summary[r.Status]++

		switch r.Status {
		case checker.StatusHealthy:
			fmt.Printf("%s  %s %s\n",
				okStyle.Render("✓"),
				dimStyle.Render(r.Dep.Path),
				dimStyle.Render(r.Dep.Version))
		case checker.StatusWarning:
			fmt.Println(warnStyle.Render("⚠"), r.Dep.Path, r.Dep.Version)
			fmt.Println("   ", dimStyle.Render(r.Message))
			for _, alt := range r.Alternatives {
				fmt.Println("   ", "→", boldStyle.Render(alt.Path), dimStyle.Render("-", alt.Reason))
			}
			fmt.Println()
		case checker.StatusCritical:
			fmt.Println(critStyle.Render("✗"), r.Dep.Path, r.Dep.Version)
			fmt.Println("   ", critStyle.Render(r.Message))
			for _, alt := range r.Alternatives {
				fmt.Println("   ", "→", boldStyle.Render(alt.Path), dimStyle.Render("-", alt.Reason))
			}
			fmt.Println()
		}
	}

	fmt.Printf("\n%s\n  %s %d healthy\n  %s %d warnings\n  %s %d critical\n",
		boldStyle.Render("📊 Summary:"),
		okStyle.Render("✓"), summary[checker.StatusHealthy],
		warnStyle.Render("⚠"), summary[checker.StatusWarning],
		critStyle.Render("✗"), summary[checker.StatusCritical])
}

func PrintJSON(results []checker.CheckResult) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(results)
}

func HasCritical(results []checker.CheckResult) bool {
	for _, r := range results {
		if r.Status == checker.StatusCritical {
			return true
		}
	}
	return false
}

func CountCritical(results []checker.CheckResult) int {
	count := 0
	for _, r := range results {
		if r.Status == checker.StatusCritical {
			count++
		}
	}
	return count
}

func CountWarning(results []checker.CheckResult) int {
	count := 0
	for _, r := range results {
		if r.Status == checker.StatusWarning {
			count++
		}
	}
	return count
}
