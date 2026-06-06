package parser

import (
	"os"

	"golang.org/x/mod/modfile"
)

type Dep struct {
	Path     string
	Version  string
	Indirect bool
}

func ParseGoMod(path string) ([]Dep, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := modfile.Parse(path, data, nil)
	if err != nil {
		return nil, err
	}

	var deps []Dep
	for _, req := range f.Require {
		deps = append(deps, Dep{
			Path:     req.Mod.Path,
			Version:  req.Mod.Version,
			Indirect: req.Indirect,
		})
	}
	return deps, nil
}
