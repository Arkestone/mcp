// Package github provides a client for the GitHub Contents API.
// It fetches file contents and directory listings from repositories.
package github

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Client accesses the GitHub Contents API.
type Client struct {
	BaseURL    string       // defaults to "https://api.github.com"
	Token      string       // optional bearer token
	HTTPClient *http.Client // optional; falls back to http.DefaultClient
}

// ContentEntry represents a file or directory in a GitHub repository listing.
type ContentEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"` // "file" or "dir"
}

// FetchFile downloads the raw content of a file from a GitHub repository.
func (c *Client) FetchFile(ctx context.Context, owner, repo, ref, path string) (string, error) {
	url := c.buildURL(owner, repo, ref, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.raw+json")
	c.setAuth(req)

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", c.httpError(resp.StatusCode, path, resp.Header, body)
	}

	return string(body), nil
}

// FetchDir lists the contents of a directory in a GitHub repository.
func (c *Client) FetchDir(ctx context.Context, owner, repo, ref, path string) ([]ContentEntry, error) {
	url := c.buildURL(owner, repo, ref, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	c.setAuth(req)

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, c.httpError(resp.StatusCode, path, resp.Header, body)
	}

	var entries []ContentEntry
	if err := json.NewDecoder(strings.NewReader(string(body))).Decode(&entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// FetchDirRecursive lists all files under root recursively using BFS.
// Directories are not included in the results.
// The initial call to root propagates errors; errors on individual subdirectories are silently skipped.
// Pass root="" to start from the repository root.
func (c *Client) FetchDirRecursive(ctx context.Context, owner, repo, ref, root string) ([]ContentEntry, error) {
	var results []ContentEntry
	queue := []string{root}
	initial := true
	for len(queue) > 0 {
		dir := queue[0]
		queue = queue[1:]
		entries, err := c.FetchDir(ctx, owner, repo, ref, dir)
		if err != nil {
			if initial {
				return nil, err
			}
			continue
		}
		initial = false
		for _, e := range entries {
			switch e.Type {
			case "file":
				results = append(results, e)
			case "dir":
				queue = append(queue, e.Path)
			}
		}
	}
	return results, nil
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *Client) buildURL(owner, repo, ref, path string) string {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = "https://api.github.com"
	}
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", base, owner, repo, path)
	if ref != "" {
		url += "?ref=" + ref
	}
	return url
}

func (c *Client) setAuth(req *http.Request) {
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
}

func (c *Client) httpError(code int, path string, headers http.Header, body []byte) error {
	// Try to extract GitHub's error message from the JSON response body.
	var ghErr struct {
		Message string `json:"message"`
	}
	_ = json.Unmarshal(body, &ghErr)
	msg := strings.ToLower(ghErr.Message)

	// Rate limit: GitHub returns 403 with "rate limit" in the message,
	// or sets X-RateLimit-Remaining: 0.
	if code == 403 && (strings.Contains(msg, "rate limit") || headers.Get("X-RateLimit-Remaining") == "0") {
		hint := "set GITHUB_TOKEN to increase the rate limit (5000 req/hr vs 60 req/hr)"
		return fmt.Errorf("HTTP 403 for %s: rate limit exceeded — %s", path, hint)
	}

	// Auth / private repo errors when no token is set.
	if (code == 401 || code == 403 || code == 404) && c.Token == "" {
		return fmt.Errorf("HTTP %d for %s (no GITHUB_TOKEN set — is this a private repo?)", code, path)
	}
	return fmt.Errorf("HTTP %d for %s", code, path)
}

// IsRateLimitError reports whether err is a GitHub API rate-limit error.
func IsRateLimitError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "rate limit exceeded")
}

// FetchZipAndExtract downloads the repository as a ZIP archive and extracts all
// files into targetDir, preserving the directory tree but stripping the
// top-level "{repo}-{sha}/" prefix that GitHub adds inside the archive.
//
// It uses the GitHub zipball API endpoint which supports both authenticated and
// unauthenticated (public repos) requests. For public repos this bypasses the
// 60 req/hr Contents API rate limit — the zipball is served from codeload.github.com.
//
// ref may be a branch, tag, or commit SHA; an empty ref uses the default branch.
func (c *Client) FetchZipAndExtract(ctx context.Context, owner, repo, ref, targetDir string) error {
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = "https://api.github.com"
	}
	url := fmt.Sprintf("%s/repos/%s/%s/zipball/%s", base, owner, repo, ref)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	c.setAuth(req)

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return c.httpError(resp.StatusCode, owner+"/"+repo, resp.Header, body)
	}

	return extractZip(body, targetDir)
}

// extractZip extracts a GitHub zipball (in-memory) into targetDir.
// GitHub zips have a single top-level directory "owner-repo-sha/" that is stripped.
func extractZip(data []byte, targetDir string) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("opening zip: %w", err)
	}

	// Determine the top-level prefix to strip (e.g. "awesome-copilot-abc1234/")
	prefix := ""
	for _, f := range zr.File {
		if strings.Contains(f.Name, "/") {
			prefix = f.Name[:strings.Index(f.Name, "/")+1]
			break
		}
	}

	for _, f := range zr.File {
		// Strip top-level prefix
		rel := strings.TrimPrefix(f.Name, prefix)
		if rel == "" {
			continue
		}

		destPath := filepath.Join(targetDir, filepath.FromSlash(rel))

		// Security: prevent zip slip
		if !strings.HasPrefix(destPath+string(os.PathSeparator), targetDir+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(destPath, 0o755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			continue
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}
		_ = os.WriteFile(destPath, content, 0o644)
	}
	return nil
}
