package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ModuleInfo struct {
	ModulePath    string `json:"modulePath"`
	LatestVersion string `json:"latestVersion"`
	CommitTime    string `json:"commitTime"`
	Deprecated    bool   `json:"deprecated"`
	Retracted     bool   `json:"retracted"`
}

func FetchModuleInfo(path string) (*ModuleInfo, error) {
	url := fmt.Sprintf("https://pkg.go.dev/v1beta/module/%s", path)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fetch %s: status %d: %s", path, resp.StatusCode, string(body))
	}

	var info ModuleInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	return &info, nil
}
