package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
		fmt.Println("Scanning dependencies...")
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
	}
}
