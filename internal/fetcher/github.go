package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func parseModulePath(modulePath string) (owner, repo string, err error) {
	parts := strings.SplitN(modulePath, "/", 4)
	if len(parts) < 3 {
		return "", "", fmt.Errorf("not a GitHub module path: %s", modulePath)
	}
	return parts[1], parts[2], nil
}

func ghGet(url string, result any) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API: status %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

type GitHubCommit struct {
	Commit struct {
		Committer struct {
			Date string `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
}

func FetchLatestCommitTime(modulePath string) (time.Time, error) {
	owner, repo, err := parseModulePath(modulePath)
	if err != nil {
		return time.Time{}, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?per_page=1", owner, repo)
	var commits []GitHubCommit
	if err := ghGet(url, &commits); err != nil {
		return time.Time{}, err
	}

	if len(commits) == 0 {
		return time.Time{}, fmt.Errorf("no commits found")
	}

	t, err := time.Parse(time.RFC3339, commits[0].Commit.Committer.Date)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func FetchLatestReleaseVersion(modulePath string) (string, error) {
	owner, repo, err := parseModulePath(modulePath)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	var release GitHubRelease
	if err := ghGet(url, &release); err != nil {
		return "", err
	}
	return release.TagName, nil
}

func IsGitHubModule(path string) bool {
	return strings.HasPrefix(path, "github.com/")
}
