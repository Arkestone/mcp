package github

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// reqCapture wraps a handler that captures the last received request.
type reqCapture struct {
	Req *http.Request
}

func newServer(t *testing.T, status int, body string) (*httptest.Server, *reqCapture) {
	t.Helper()
	rc := &reqCapture{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc.Req = r.Clone(r.Context())
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv, rc
}

// ---------------------------------------------------------------------------
// FetchFile tests
// ---------------------------------------------------------------------------

func TestFetchFile_Basic(t *testing.T) {
	srv, rc := newServer(t, http.StatusOK, "file-content")
	c := &Client{BaseURL: srv.URL}

	got, err := c.FetchFile(context.Background(), "owner", "repo", "", "README.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "file-content" {
		t.Errorf("body = %q, want %q", got, "file-content")
	}
	if accept := rc.Req.Header.Get("Accept"); accept != "application/vnd.github.raw+json" {
		t.Errorf("Accept = %q, want %q", accept, "application/vnd.github.raw+json")
	}
	if !strings.HasSuffix(rc.Req.URL.Path, "/repos/owner/repo/contents/README.md") {
		t.Errorf("unexpected path: %s", rc.Req.URL.Path)
	}
}

func TestFetchFile_WithAuthToken(t *testing.T) {
	srv, rc := newServer(t, http.StatusOK, "ok")
	c := &Client{BaseURL: srv.URL, Token: "ghp_secret"}

	_, err := c.FetchFile(context.Background(), "o", "r", "", "f")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth := rc.Req.Header.Get("Authorization"); auth != "Bearer ghp_secret" {
		t.Errorf("Authorization = %q, want %q", auth, "Bearer ghp_secret")
	}
}

func TestFetchFile_WithoutAuthToken(t *testing.T) {
	srv, rc := newServer(t, http.StatusOK, "ok")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchFile(context.Background(), "o", "r", "", "f")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth := rc.Req.Header.Get("Authorization"); auth != "" {
		t.Errorf("Authorization = %q, want empty", auth)
	}
}

func TestFetchFile_WithRef(t *testing.T) {
	srv, rc := newServer(t, http.StatusOK, "ok")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchFile(context.Background(), "o", "r", "v1.0", "f")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref := rc.Req.URL.Query().Get("ref"); ref != "v1.0" {
		t.Errorf("ref = %q, want %q", ref, "v1.0")
	}
}

func TestFetchFile_ServerError500(t *testing.T) {
	srv, _ := newServer(t, http.StatusInternalServerError, "boom")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchFile(context.Background(), "o", "r", "", "f")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error = %q, want it to contain '500'", err.Error())
	}
}

func TestFetchFile_NotFound404(t *testing.T) {
	srv, _ := newServer(t, http.StatusNotFound, "not found")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchFile(context.Background(), "o", "r", "", "missing.txt")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error = %q, want it to contain '404'", err.Error())
	}
}

func TestFetchFile_EmptyBaseURLDefault(t *testing.T) {
	c := &Client{}
	url := c.buildURL("owner", "repo", "main", "path/to/file")
	want := "https://api.github.com/repos/owner/repo/contents/path/to/file?ref=main"
	if url != want {
		t.Errorf("buildURL = %q, want %q", url, want)
	}
}

// ---------------------------------------------------------------------------
// FetchDir tests
// ---------------------------------------------------------------------------

func TestFetchDir_Basic(t *testing.T) {
	entries := []ContentEntry{
		{Name: "README.md", Path: "README.md", Type: "file"},
		{Name: "src", Path: "src", Type: "dir"},
	}
	body, _ := json.Marshal(entries)
	srv, rc := newServer(t, http.StatusOK, string(body))
	c := &Client{BaseURL: srv.URL}

	got, err := c.FetchDir(context.Background(), "owner", "repo", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Name != "README.md" || got[0].Type != "file" {
		t.Errorf("entry[0] = %+v, want README.md file", got[0])
	}
	if got[1].Name != "src" || got[1].Type != "dir" {
		t.Errorf("entry[1] = %+v, want src dir", got[1])
	}
	if accept := rc.Req.Header.Get("Accept"); accept != "application/vnd.github+json" {
		t.Errorf("Accept = %q, want %q", accept, "application/vnd.github+json")
	}
}

func TestFetchDir_WithAuthToken(t *testing.T) {
	body, _ := json.Marshal([]ContentEntry{})
	srv, rc := newServer(t, http.StatusOK, string(body))
	c := &Client{BaseURL: srv.URL, Token: "tok123"}

	_, err := c.FetchDir(context.Background(), "o", "r", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth := rc.Req.Header.Get("Authorization"); auth != "Bearer tok123" {
		t.Errorf("Authorization = %q, want %q", auth, "Bearer tok123")
	}
}

func TestFetchDir_NotFound404(t *testing.T) {
	srv, _ := newServer(t, http.StatusNotFound, "not found")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchDir(context.Background(), "o", "r", "", "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error = %q, want it to contain '404'", err.Error())
	}
}

func TestFetchDir_InvalidJSON(t *testing.T) {
	srv, _ := newServer(t, http.StatusOK, "not-json{{{")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchDir(context.Background(), "o", "r", "", "dir")
	if err == nil {
		t.Fatal("expected JSON decode error, got nil")
	}
}

func TestFetchDir_EmptyDirectory(t *testing.T) {
	body, _ := json.Marshal([]ContentEntry{})
	srv, _ := newServer(t, http.StatusOK, string(body))
	c := &Client{BaseURL: srv.URL}

	got, err := c.FetchDir(context.Background(), "o", "r", "", "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

// ---------------------------------------------------------------------------
// Edge-case tests
// ---------------------------------------------------------------------------

func TestBuildURL_TrailingSlash(t *testing.T) {
	c := &Client{BaseURL: "https://example.com/"}
	got := c.buildURL("owner", "repo", "", "file.txt")
	// Should not produce double slash between base and "repos"
	if strings.Contains(got, "//repos") {
		t.Errorf("buildURL produced double slash: %s", got)
	}
}

func TestFetchFile_SpecialCharactersInPath(t *testing.T) {
	srv, rc := newServer(t, http.StatusOK, "content")
	c := &Client{BaseURL: srv.URL}

	path := "dir/my file (1).md"
	_, err := c.FetchFile(context.Background(), "o", "r", "", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The path segment should be present in the request URL
	if !strings.Contains(rc.Req.URL.Path, "my file (1).md") &&
		!strings.Contains(rc.Req.URL.RawPath, "my%20file%20%281%29.md") &&
		!strings.Contains(rc.Req.RequestURI, "my") {
		t.Errorf("path with special chars not in request URL: %s", rc.Req.URL.String())
	}
}

func TestFetchFile_NetworkError(t *testing.T) {
	c := &Client{BaseURL: "http://127.0.0.1:1"} // port 1 — connection refused
	_, err := c.FetchFile(context.Background(), "o", "r", "", "f.md")
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

func TestFetchDir_NetworkError(t *testing.T) {
	c := &Client{BaseURL: "http://127.0.0.1:1"}
	_, err := c.FetchDir(context.Background(), "o", "r", "", "dir")
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
}

func TestHttpClient_Default(t *testing.T) {
	c := &Client{}
	got := c.httpClient()
	if got != http.DefaultClient {
		t.Error("nil HTTPClient should return http.DefaultClient")
	}
}

func TestHttpClient_Custom(t *testing.T) {
	custom := &http.Client{}
	c := &Client{HTTPClient: custom}
	got := c.httpClient()
	if got != custom {
		t.Error("should return the custom HTTPClient")
	}
}

func TestFetchFile_CancelledContext(t *testing.T) {
	srv, _ := newServer(t, http.StatusOK, "content")
	c := &Client{BaseURL: srv.URL}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	_, err := c.FetchFile(ctx, "o", "r", "", "f.md")
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestFetchDir_CancelledContext(t *testing.T) {
	srv, _ := newServer(t, http.StatusOK, "[]")
	c := &Client{BaseURL: srv.URL}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.FetchDir(ctx, "o", "r", "", "dir")
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestHttpError_NoTokenHint(t *testing.T) {
	c := &Client{}
	for _, code := range []int{401, 403, 404} {
		err := c.httpError(code, "test/path", http.Header{}, nil)
		if !strings.Contains(err.Error(), "no GITHUB_TOKEN") {
			t.Errorf("httpError(%d) = %q, want token hint", code, err.Error())
		}
	}
}

func TestHttpError_WithTokenNoHint(t *testing.T) {
	c := &Client{Token: "ghp_test"}
	for _, code := range []int{401, 403, 404} {
		err := c.httpError(code, "test/path", http.Header{}, nil)
		if strings.Contains(err.Error(), "no GITHUB_TOKEN") {
			t.Errorf("httpError(%d) = %q, should NOT contain token hint when token is set", code, err.Error())
		}
	}
}

func TestHttpError_OtherCodes(t *testing.T) {
	c := &Client{}
	err := c.httpError(500, "test/path", http.Header{}, nil)
	if strings.Contains(err.Error(), "no GITHUB_TOKEN") {
		t.Errorf("httpError(500) = %q, should NOT suggest token for server error", err.Error())
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("httpError(500) = %q, should contain status code", err.Error())
	}
}

func TestFetchFile_NoToken404HintsAboutAuth(t *testing.T) {
	srv, _ := newServer(t, http.StatusNotFound, "not found")
	c := &Client{BaseURL: srv.URL}
	_, err := c.FetchFile(context.Background(), "o", "r", "", "private.md")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "no GITHUB_TOKEN") {
		t.Errorf("error = %q, want hint about missing token", err.Error())
	}
}

func TestFetchFile_WithToken404NoHint(t *testing.T) {
	srv, _ := newServer(t, http.StatusNotFound, "not found")
	c := &Client{BaseURL: srv.URL, Token: "ghp_test"}
	_, err := c.FetchFile(context.Background(), "o", "r", "", "missing.md")
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "no GITHUB_TOKEN") {
		t.Errorf("error = %q, should NOT hint about token when one is set", err.Error())
	}
}

// ---------------------------------------------------------------------------
// Additional nominal / error / limit tests
// ---------------------------------------------------------------------------

func TestFetchDir_WithRef(t *testing.T) {
	body, _ := json.Marshal([]ContentEntry{{Name: "a.go", Type: "file"}})
	srv, rc := newServer(t, http.StatusOK, string(body))
	c := &Client{BaseURL: srv.URL}

	got, err := c.FetchDir(context.Background(), "o", "r", "v2.0", "pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if !strings.Contains(rc.Req.URL.RawQuery, "ref=v2.0") {
		t.Errorf("query = %q, want ref=v2.0", rc.Req.URL.RawQuery)
	}
}

func TestFetchFile_EmptyBody(t *testing.T) {
	srv, _ := newServer(t, http.StatusOK, "")
	c := &Client{BaseURL: srv.URL}

	got, err := c.FetchFile(context.Background(), "o", "r", "", "empty.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestFetchFile_LargeBody(t *testing.T) {
	large := strings.Repeat("A", 1024*1024) // 1MB
	srv, _ := newServer(t, http.StatusOK, large)
	c := &Client{BaseURL: srv.URL}

	got, err := c.FetchFile(context.Background(), "o", "r", "", "big.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1024*1024 {
		t.Errorf("len = %d, want %d", len(got), 1024*1024)
	}
}

func TestFetchDir_ManyEntries(t *testing.T) {
	entries := make([]ContentEntry, 150)
	for i := range entries {
		entries[i] = ContentEntry{Name: fmt.Sprintf("file%d.go", i), Type: "file"}
	}
	body, _ := json.Marshal(entries)
	srv, _ := newServer(t, http.StatusOK, string(body))
	c := &Client{BaseURL: srv.URL}

	got, err := c.FetchDir(context.Background(), "o", "r", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 150 {
		t.Errorf("len = %d, want 150", len(got))
	}
}

func TestFetchFile_HTTP401NoToken(t *testing.T) {
	srv, _ := newServer(t, http.StatusUnauthorized, "unauthorized")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchFile(context.Background(), "o", "r", "", "secret.md")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should mention 401: %v", err)
	}
	if !strings.Contains(err.Error(), "no GITHUB_TOKEN") {
		t.Errorf("error should hint about missing token: %v", err)
	}
}

func TestFetchFile_HTTP403NoToken(t *testing.T) {
	srv, _ := newServer(t, http.StatusForbidden, "forbidden")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchFile(context.Background(), "o", "r", "", "private.md")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("error should mention 403: %v", err)
	}
	if !strings.Contains(err.Error(), "no GITHUB_TOKEN") {
		t.Errorf("error should hint about missing token: %v", err)
	}
}

func TestFetchDir_HTTP401NoToken(t *testing.T) {
	srv, _ := newServer(t, http.StatusUnauthorized, "unauthorized")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchDir(context.Background(), "o", "r", "", "dir")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should mention 401: %v", err)
	}
}

func TestFetchDir_HTTP403RateLimitViaBody(t *testing.T) {
	// GitHub returns a JSON body with "rate limit" in the message when rate-limited.
	body := `{"message":"API rate limit exceeded for 1.2.3.4.","documentation_url":"https://docs.github.com/..."}`
	srv, _ := newServer(t, http.StatusForbidden, body)
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchDir(context.Background(), "o", "r", "", "dir")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "rate limit") {
		t.Errorf("rate-limit 403 error should mention rate limit, got: %v", err)
	}
	if strings.Contains(err.Error(), "private repo") {
		t.Errorf("rate-limit error should NOT say 'private repo', got: %v", err)
	}
	if !strings.Contains(err.Error(), "GITHUB_TOKEN") {
		t.Errorf("rate-limit error should suggest setting GITHUB_TOKEN, got: %v", err)
	}
}

func TestFetchDir_HTTP403RateLimitViaHeader(t *testing.T) {
	// Rate limit can also be detected via X-RateLimit-Remaining: 0 header.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"some error"}`))
	}))
	t.Cleanup(srv.Close)

	c := &Client{BaseURL: srv.URL}
	_, err := c.FetchDir(context.Background(), "o", "r", "", "dir")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "rate limit") {
		t.Errorf("rate-limit header 403 error should mention rate limit, got: %v", err)
	}
}

func TestFetchFile_HTTP403RateLimitViaBody(t *testing.T) {
	body := `{"message":"API rate limit exceeded for 5.5.5.5."}`
	srv, _ := newServer(t, http.StatusForbidden, body)
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchFile(context.Background(), "o", "r", "", "file.md")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "rate limit") {
		t.Errorf("rate-limit 403 FetchFile error should mention rate limit, got: %v", err)
	}
}

func TestFetchDir_SpecialCharsInPath(t *testing.T) {
	body, _ := json.Marshal([]ContentEntry{})
	srv, rc := newServer(t, http.StatusOK, string(body))
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchDir(context.Background(), "o", "r", "", "path/with spaces/dir")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(rc.Req.URL.Path, "path/with spaces/dir") {
		t.Errorf("path = %q, should contain special chars", rc.Req.URL.Path)
	}
}

func TestBuildURL_NoRef(t *testing.T) {
	c := &Client{BaseURL: "https://api.github.com"}
	url := c.buildURL("owner", "repo", "", "path")
	if strings.Contains(url, "ref=") {
		t.Errorf("URL should not contain ref= when empty: %s", url)
	}
}

func TestFetchFile_HTTP200WithToken(t *testing.T) {
	srv, _ := newServer(t, http.StatusOK, "content")
	c := &Client{BaseURL: srv.URL, Token: "ghp_test"}

	got, err := c.FetchFile(context.Background(), "o", "r", "", "file.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "content" {
		t.Errorf("got %q, want %q", got, "content")
	}
}

// TestFetchFile_RequestCreationError covers the http.NewRequestWithContext error
// branch by supplying a BaseURL that produces an invalid URL (space in host).
func TestFetchFile_RequestCreationError(t *testing.T) {
	c := &Client{BaseURL: "http://invalid host"}
	_, err := c.FetchFile(context.Background(), "o", "r", "", "f.md")
	if err == nil {
		t.Fatal("expected error for malformed URL")
	}
}

// errReader returns an error on every Read, simulating a broken response body.
type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, fmt.Errorf("read error") }
func (errReader) Close() error             { return nil }

// brokenBodyTransport returns HTTP 200 with a body that always errors on Read.
type brokenBodyTransport struct{}

func (brokenBodyTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       errReader{},
		Header:     make(http.Header),
	}, nil
}

func TestFetchFile_BodyReadError(t *testing.T) {
	c := &Client{HTTPClient: &http.Client{Transport: brokenBodyTransport{}}, BaseURL: "http://example.com"}
	_, err := c.FetchFile(context.Background(), "o", "r", "", "f.md")
	if err == nil {
		t.Fatal("expected error from broken body, got nil")
	}
}

func TestFetchDir_DecodeError(t *testing.T) {
	srv, _ := newServer(t, http.StatusOK, "not-valid-json")
	c := &Client{BaseURL: srv.URL}

	_, err := c.FetchDir(context.Background(), "o", "r", "", "dir")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}

type errTransport struct{}

func (errTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("transport error")
}

func TestFetchDir_DoError(t *testing.T) {
	c := &Client{HTTPClient: &http.Client{Transport: errTransport{}}, BaseURL: "http://example.com"}
	_, err := c.FetchDir(context.Background(), "o", "r", "", "dir")
	if err == nil {
		t.Fatal("expected transport error, got nil")
	}
}

// ---------------------------------------------------------------------------
// FetchDirRecursive tests
// ---------------------------------------------------------------------------

// recursiveServer serves multiple directory listings keyed by URL path suffix.
func recursiveServer(t *testing.T, routes map[string]interface{}) *httptest.Server {
t.Helper()
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Strip query string for matching.
path := r.URL.Path
for key, val := range routes {
if strings.HasSuffix(path, key) {
w.Header().Set("Content-Type", "application/json")
if err := json.NewEncoder(w).Encode(val); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
}
return
}
}
w.WriteHeader(http.StatusNotFound)
}))
t.Cleanup(srv.Close)
return srv
}

func TestFetchDirRecursive_Basic(t *testing.T) {
// Root contains one file and one subdir; subdir contains two files.
rootEntries := []ContentEntry{
{Name: "README.md", Path: "README.md", Type: "file"},
{Name: "src", Path: "src", Type: "dir"},
}
srcEntries := []ContentEntry{
{Name: "main.go", Path: "src/main.go", Type: "file"},
{Name: "util.go", Path: "src/util.go", Type: "file"},
}
srv := recursiveServer(t, map[string]interface{}{
"/repos/o/r/contents/": rootEntries,
"/repos/o/r/contents/src": srcEntries,
})
c := &Client{BaseURL: srv.URL}

got, err := c.FetchDirRecursive(context.Background(), "o", "r", "", "")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 3 {
t.Fatalf("expected 3 files, got %d: %v", len(got), got)
}
paths := map[string]bool{}
for _, e := range got {
paths[e.Path] = true
}
for _, want := range []string{"README.md", "src/main.go", "src/util.go"} {
if !paths[want] {
t.Errorf("expected path %q in results", want)
}
}
}

func TestFetchDirRecursive_SubdirErrorSkipped(t *testing.T) {
// Root contains a file and a subdir; the subdir returns 404 → silently skipped.
rootEntries := []ContentEntry{
{Name: "README.md", Path: "README.md", Type: "file"},
{Name: "missing", Path: "missing", Type: "dir"},
}
srv := recursiveServer(t, map[string]interface{}{
"/repos/o/r/contents/": rootEntries,
// "missing" subdir is not registered → 404
})
c := &Client{BaseURL: srv.URL}

got, err := c.FetchDirRecursive(context.Background(), "o", "r", "", "")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 1 || got[0].Path != "README.md" {
t.Errorf("expected only README.md, got %v", got)
}
}

func TestFetchDirRecursive_RootError(t *testing.T) {
// If the root directory itself returns an error, it propagates.
srv := recursiveServer(t, map[string]interface{}{}) // no routes → all 404
c := &Client{BaseURL: srv.URL}

_, err := c.FetchDirRecursive(context.Background(), "o", "r", "", "")
if err == nil {
t.Fatal("expected error from root 404, got nil")
}
}

func TestFetchDirRecursive_EmptyRoot(t *testing.T) {
srv := recursiveServer(t, map[string]interface{}{
"/repos/o/r/contents/": []ContentEntry{},
})
c := &Client{BaseURL: srv.URL}

got, err := c.FetchDirRecursive(context.Background(), "o", "r", "", "")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if len(got) != 0 {
t.Errorf("expected empty result, got %v", got)
}
}

func TestFetchDirRecursive_WithRef(t *testing.T) {
rootEntries := []ContentEntry{
{Name: "file.go", Path: "file.go", Type: "file"},
}
var capturedRef string
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
capturedRef = r.URL.Query().Get("ref")
w.Header().Set("Content-Type", "application/json")
_ = json.NewEncoder(w).Encode(rootEntries)
}))
t.Cleanup(srv.Close)
c := &Client{BaseURL: srv.URL}

_, err := c.FetchDirRecursive(context.Background(), "o", "r", "v1.2.3", "")
if err != nil {
t.Fatalf("unexpected error: %v", err)
}
if capturedRef != "v1.2.3" {
t.Errorf("ref = %q, want %q", capturedRef, "v1.2.3")
}
}

// ---------------------------------------------------------------------------
// FetchZipAndExtract
// ---------------------------------------------------------------------------

func TestFetchZipAndExtract_Success(t *testing.T) {
// Build an in-memory ZIP with a top-level prefix directory (as GitHub does).
zipData := buildTestZip(t, map[string]string{
"myrepo-abc1234/README.md":        "# readme",
"myrepo-abc1234/src/main.go":      "package main",
"myrepo-abc1234/docs/guide.md":    "# guide",
})

srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/zip")
_, _ = w.Write(zipData)
}))
t.Cleanup(srv.Close)

destDir := t.TempDir()
c := &Client{BaseURL: srv.URL}
if err := c.FetchZipAndExtract(context.Background(), "owner", "myrepo", "main", destDir); err != nil {
t.Fatalf("unexpected error: %v", err)
}

for _, want := range []struct{ path, content string }{
{"README.md", "# readme"},
{"src/main.go", "package main"},
{"docs/guide.md", "# guide"},
} {
got, err := os.ReadFile(destDir + "/" + want.path)
if err != nil {
t.Errorf("missing %s: %v", want.path, err)
continue
}
if string(got) != want.content {
t.Errorf("%s = %q, want %q", want.path, got, want.content)
}
}
}

func TestFetchZipAndExtract_HTTPError(t *testing.T) {
srv, _ := newServer(t, 403, `{"message":"API rate limit exceeded"}`)
c := &Client{BaseURL: srv.URL}
err := c.FetchZipAndExtract(context.Background(), "o", "r", "", t.TempDir())
if err == nil {
t.Fatal("expected error, got nil")
}
if !strings.Contains(err.Error(), "rate limit") {
t.Errorf("expected rate limit error, got: %v", err)
}
}

func TestIsRateLimitError(t *testing.T) {
tests := []struct {
err  error
want bool
}{
{nil, false},
{fmt.Errorf("connection refused"), false},
{fmt.Errorf("HTTP 404 for owner/repo"), false},
{fmt.Errorf("HTTP 403 for /path: rate limit exceeded — set GITHUB_TOKEN"), true},
{fmt.Errorf("rate limit exceeded"), true},
}
for _, tt := range tests {
if got := IsRateLimitError(tt.err); got != tt.want {
t.Errorf("IsRateLimitError(%v) = %v, want %v", tt.err, got, tt.want)
}
}
}

// buildTestZip creates an in-memory zip from a map of path→content.
func buildTestZip(t *testing.T, files map[string]string) []byte {
t.Helper()
var buf bytes.Buffer
zw := zip.NewWriter(&buf)
for name, content := range files {
w, err := zw.Create(name)
if err != nil {
t.Fatalf("creating zip entry %s: %v", name, err)
}
_, _ = w.Write([]byte(content))
}
if err := zw.Close(); err != nil {
t.Fatalf("closing zip: %v", err)
}
return buf.Bytes()
}
