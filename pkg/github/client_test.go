package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
		err := c.httpError(code, "test/path")
		if !strings.Contains(err.Error(), "no GITHUB_TOKEN") {
			t.Errorf("httpError(%d) = %q, want token hint", code, err.Error())
		}
	}
}

func TestHttpError_WithTokenNoHint(t *testing.T) {
	c := &Client{Token: "ghp_test"}
	for _, code := range []int{401, 403, 404} {
		err := c.httpError(code, "test/path")
		if strings.Contains(err.Error(), "no GITHUB_TOKEN") {
			t.Errorf("httpError(%d) = %q, should NOT contain token hint when token is set", code, err.Error())
		}
	}
}

func TestHttpError_OtherCodes(t *testing.T) {
	c := &Client{}
	err := c.httpError(500, "test/path")
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
