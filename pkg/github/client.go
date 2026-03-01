// Package github provides a client for the GitHub Contents API.
// It fetches file contents and directory listings from repositories.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	if resp.StatusCode != 200 {
		return "", c.httpError(resp.StatusCode, path)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
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

	if resp.StatusCode != 200 {
		return nil, c.httpError(resp.StatusCode, path)
	}

	var entries []ContentEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}
	return entries, nil
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

func (c *Client) httpError(code int, path string) error {
	if (code == 401 || code == 403 || code == 404) && c.Token == "" {
		return fmt.Errorf("HTTP %d for %s (no GITHUB_TOKEN set — is this a private repo?)", code, path)
	}
	return fmt.Errorf("HTTP %d for %s", code, path)
}
