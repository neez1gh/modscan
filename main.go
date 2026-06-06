package main

import (
	_ "embed"

	"modscan/cmd"
)

//go:embed rules/alternatives.yaml
var defaultRules []byte

func main() {
	cmd.DefaultRules = defaultRules
	cmd.Execute()
}
