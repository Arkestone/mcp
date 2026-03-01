package optimizer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewNilWhenNotConfigured(t *testing.T) {
	tests := []struct {
		name string
		cfg  LLMConfig
	}{
		{"empty", LLMConfig{}},
		{"no endpoint", LLMConfig{APIKey: "sk-key"}},
		{"no api key", LLMConfig{Endpoint: "http://example.com"}},
		{"both empty", LLMConfig{Endpoint: "", APIKey: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := New(tt.cfg)
			if opt != nil {
				t.Error("New should return nil when not configured")
			}
		})
	}
}

func TestNewReturnsOptimizer(t *testing.T) {
	opt := New(LLMConfig{
		Endpoint: "http://example.com",
		APIKey:   "sk-test",
		Model:    "gpt-4o-mini",
	})
	if opt == nil {
		t.Fatal("New should return non-nil when configured")
	}
}

func TestEnabledNilSafe(t *testing.T) {
	var opt *Optimizer
	if opt.Enabled() {
		t.Error("nil Optimizer.Enabled() should return false")
	}
}

func TestEnabledTrue(t *testing.T) {
	opt := New(LLMConfig{
		Endpoint: "http://example.com",
		APIKey:   "sk-test",
	})
	if !opt.Enabled() {
		t.Error("configured Optimizer.Enabled() should return true")
	}
}

func TestConcatRawMultiple(t *testing.T) {
	inputs := []ContentInput{
		{Source: "repo1", Path: ".github/copilot-instructions.md", Content: "Use Go"},
		{Source: "repo2", Path: ".github/instructions/style.instructions.md", Content: "Use gofmt"},
	}
	result := ConcatRaw(inputs)

	if !strings.Contains(result, "repo1") {
		t.Error("result should contain source repo1")
	}
	if !strings.Contains(result, "repo2") {
		t.Error("result should contain source repo2")
	}
	if !strings.Contains(result, "Use Go") {
		t.Error("result should contain first instruction content")
	}
	if !strings.Contains(result, "Use gofmt") {
		t.Error("result should contain second instruction content")
	}
}

func TestConcatRawEmpty(t *testing.T) {
	result := ConcatRaw(nil)
	if result != "" {
		t.Errorf("ConcatRaw(nil) = %q, want empty", result)
	}
}

func TestConcatRawSingle(t *testing.T) {
	inputs := []ContentInput{
		{Source: "src", Path: "path", Content: "content"},
	}
	result := ConcatRaw(inputs)
	if !strings.Contains(result, "content") {
		t.Error("result should contain the content")
	}
	if !strings.Contains(result, "src") {
		t.Error("result should contain the source")
	}
}

func TestOptimizeSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate request
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Authorization") != "Bearer sk-test" {
			t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
		}

		// Validate body
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("invalid request body: %v", err)
		}
		if req.Model != "test-model" {
			t.Errorf("model = %q", req.Model)
		}
		if len(req.Messages) != 2 {
			t.Fatalf("messages count = %d, want 2", len(req.Messages))
		}
		if req.Messages[0].Role != "system" {
			t.Errorf("messages[0].role = %q", req.Messages[0].Role)
		}
		if req.Messages[1].Role != "user" {
			t.Errorf("messages[1].role = %q", req.Messages[1].Role)
		}

		resp := chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Role: "assistant", Content: "# Consolidated\nUse Go and gofmt"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "test-model",
	})

	inputs := []ContentInput{
		{Source: "repo1", Path: "p1", Content: "Use Go"},
		{Source: "repo2", Path: "p2", Content: "Use gofmt"},
	}

	result, err := opt.Optimize(context.Background(), inputs)
	if err != nil {
		t.Fatalf("Optimize failed: %v", err)
	}
	if result != "# Consolidated\nUse Go and gofmt" {
		t.Errorf("result = %q", result)
	}
}

func TestOptimizeHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	_, err := opt.Optimize(context.Background(), []ContentInput{
		{Source: "s", Path: "p", Content: "c"},
	})
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error = %v", err)
	}
}

func TestOptimizeEmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{Choices: []chatChoice{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	_, err := opt.Optimize(context.Background(), []ContentInput{
		{Source: "s", Path: "p", Content: "c"},
	})
	if err == nil {
		t.Fatal("expected error for empty choices")
	}
	if !strings.Contains(err.Error(), "no choices") {
		t.Errorf("error = %v", err)
	}
}

func TestOptimizeMalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	_, err := opt.Optimize(context.Background(), []ContentInput{
		{Source: "s", Path: "p", Content: "c"},
	})
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestOptimizeEndpointNormalization(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		wantPath string
	}{
		{"bare URL", "URL_PLACEHOLDER", "/chat/completions"},
		{"with trailing slash", "URL_PLACEHOLDER/", "/chat/completions"},
		{"already has path", "URL_PLACEHOLDER/chat/completions", "/chat/completions"},
		{"with v1 prefix", "URL_PLACEHOLDER/v1", "/v1/chat/completions"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				resp := chatResponse{
					Choices: []chatChoice{
						{Message: chatMessage{Content: "ok"}},
					},
				}
				json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			endpoint := strings.ReplaceAll(tt.endpoint, "URL_PLACEHOLDER", server.URL)
			opt := New(LLMConfig{
				Endpoint: endpoint,
				APIKey:   "sk-test",
				Model:    "model",
			})

			opt.Optimize(context.Background(), []ContentInput{
				{Source: "s", Path: "p", Content: "c"},
			})

			if gotPath != tt.wantPath {
				t.Errorf("request path = %q, want %q", gotPath, tt.wantPath)
			}
		})
	}
}

func TestOptimizeNotEnabledFallback(t *testing.T) {
	var opt *Optimizer

	inputs := []ContentInput{
		{Source: "s1", Path: "p1", Content: "c1"},
		{Source: "s2", Path: "p2", Content: "c2"},
	}

	result, err := opt.Optimize(context.Background(), inputs)
	if err != nil {
		t.Fatalf("Optimize on nil should not error: %v", err)
	}
	// Should fall back to ConcatRaw
	if !strings.Contains(result, "c1") || !strings.Contains(result, "c2") {
		t.Errorf("result = %q, should contain both contents", result)
	}
}

func TestOptimizeTemperature(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		json.Unmarshal(body, &req)
		if req.Temperature == nil {
			t.Error("temperature should be set")
		} else if *req.Temperature != 0.2 {
			t.Errorf("temperature = %f, want 0.2", *req.Temperature)
		}
		resp := chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Content: "ok"}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})
	opt.Optimize(context.Background(), []ContentInput{{Source: "s", Path: "p", Content: "c"}})
}

func TestOptimizeUserMessageContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		json.Unmarshal(body, &req)

		userMsg := req.Messages[1].Content
		if !strings.Contains(userMsg, "Source1") {
			t.Error("user message should contain source names")
		}
		if !strings.Contains(userMsg, "content A") {
			t.Error("user message should contain instruction content")
		}
		if !strings.Contains(userMsg, "---") {
			t.Error("user message should contain separator")
		}

		resp := chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Content: "merged"}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})
	opt.Optimize(context.Background(), []ContentInput{
		{Source: "Source1", Path: "path1", Content: "content A"},
		{Source: "Source2", Path: "path2", Content: "content B"},
	})
}

func TestOptimizeContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response; the context should cancel before we respond
		<-r.Context().Done()
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := opt.Optimize(ctx, []ContentInput{
		{Source: "s", Path: "p", Content: "c"},
	})
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

func TestConcatRawPreservesOrder(t *testing.T) {
	inputs := []ContentInput{
		{Source: "alpha", Path: "a.md", Content: "first"},
		{Source: "bravo", Path: "b.md", Content: "second"},
		{Source: "charlie", Path: "c.md", Content: "third"},
	}
	result := ConcatRaw(inputs)

	idxAlpha := strings.Index(result, "alpha")
	idxBravo := strings.Index(result, "bravo")
	idxCharlie := strings.Index(result, "charlie")

	if idxAlpha >= idxBravo || idxBravo >= idxCharlie {
		t.Errorf("order not preserved: alpha@%d bravo@%d charlie@%d", idxAlpha, idxBravo, idxCharlie)
	}
}

func TestConcatRawSpecialChars(t *testing.T) {
	inputs := []ContentInput{
		{Source: "md-src", Path: "special.md", Content: "# Heading\n\n```go\nfunc main() {}\n```\n\n> blockquote"},
		{Source: "html-src", Path: "html.md", Content: "<div class=\"test\">hello &amp; world</div>"},
		{Source: "backtick-src", Path: "bt.md", Content: "Use `backticks` and ``double backticks``"},
	}
	result := ConcatRaw(inputs)

	for _, want := range []string{
		"```go", "func main() {}", "```",
		"<div class=\"test\">", "&amp;",
		"`backticks`", "``double backticks``",
	} {
		if !strings.Contains(result, want) {
			t.Errorf("result missing %q", want)
		}
	}
}

func TestOptimizeLargeInput(t *testing.T) {
	var receivedReq chatRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedReq)
		resp := chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Content: "consolidated"}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	inputs := make([]ContentInput, 50)
	for i := range inputs {
		inputs[i] = ContentInput{
			Source:  fmt.Sprintf("source-%d", i),
			Path:    fmt.Sprintf("path-%d.md", i),
			Content: strings.Repeat(fmt.Sprintf("line %d content. ", i), 20),
		}
	}

	result, err := opt.Optimize(context.Background(), inputs)
	if err != nil {
		t.Fatalf("Optimize failed: %v", err)
	}
	if result != "consolidated" {
		t.Errorf("result = %q", result)
	}
	if len(receivedReq.Messages) != 2 {
		t.Fatalf("messages count = %d, want 2", len(receivedReq.Messages))
	}
	userMsg := receivedReq.Messages[1].Content
	for i := 0; i < 50; i++ {
		if !strings.Contains(userMsg, fmt.Sprintf("source-%d", i)) {
			t.Errorf("user message missing source-%d", i)
			break
		}
	}
}

func TestOptimizeConcurrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Content: "ok"}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	inputs := []ContentInput{{Source: "s", Path: "p", Content: "c"}}

	errs := make(chan error, 5)
	for i := 0; i < 5; i++ {
		go func() {
			_, err := opt.Optimize(context.Background(), inputs)
			errs <- err
		}()
	}
	for i := 0; i < 5; i++ {
		if err := <-errs; err != nil {
			t.Errorf("goroutine %d error: %v", i, err)
		}
	}
}

func TestOptimizeHTTP429RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		w.Write([]byte("rate limit exceeded"))
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	_, err := opt.Optimize(context.Background(), []ContentInput{
		{Source: "s", Path: "p", Content: "c"},
	})
	if err == nil {
		t.Fatal("expected error for HTTP 429")
	}
	if !strings.Contains(err.Error(), "429") {
		t.Errorf("error should mention 429, got: %v", err)
	}
}

func TestOptimizeHTTP401Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte("unauthorized"))
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-bad",
		Model:    "model",
	})

	_, err := opt.Optimize(context.Background(), []ContentInput{
		{Source: "s", Path: "p", Content: "c"},
	})
	if err == nil {
		t.Fatal("expected error for HTTP 401")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error should mention 401, got: %v", err)
	}
}

func TestOptimizeEmptyInstructions(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		resp := chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Content: "empty"}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	result, err := opt.Optimize(context.Background(), []ContentInput{})
	if err != nil {
		t.Fatalf("Optimize failed: %v", err)
	}
	if !called {
		t.Error("expected HTTP request to be made even with empty instructions")
	}
	if result != "empty" {
		t.Errorf("result = %q, want %q", result, "empty")
	}
}

func TestOptimizeMultipleChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Content: "first choice"}},
				{Message: chatMessage{Content: "second choice"}},
				{Message: chatMessage{Content: "third choice"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})

	result, err := opt.Optimize(context.Background(), []ContentInput{
		{Source: "s", Path: "p", Content: "c"},
	})
	if err != nil {
		t.Fatalf("Optimize failed: %v", err)
	}
	if result != "first choice" {
		t.Errorf("result = %q, want %q", result, "first choice")
	}
}

func TestNewWithAllFields(t *testing.T) {
	opt := New(LLMConfig{
		Endpoint: "http://example.com/v1",
		APIKey:   "sk-full-test",
		Model:    "gpt-4o-mini",
		Enabled:  true,
	})
	if opt == nil {
		t.Fatal("New should return non-nil when all fields are populated")
	}
	if !opt.Enabled() {
		t.Error("Enabled() should return true")
	}
}

func TestOptimizeSystemPromptContent(t *testing.T) {
	var systemContent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		json.Unmarshal(body, &req)
		for _, m := range req.Messages {
			if m.Role == "system" {
				systemContent = m.Content
			}
		}
		resp := chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Content: "ok"}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	opt := New(LLMConfig{
		Endpoint: server.URL,
		APIKey:   "sk-test",
		Model:    "model",
	})
	opt.Optimize(context.Background(), []ContentInput{
		{Source: "s", Path: "p", Content: "c"},
	})

	for _, keyword := range []string{"merge", "duplicate", "Markdown"} {
		if !strings.Contains(strings.ToLower(systemContent), strings.ToLower(keyword)) {
			t.Errorf("system prompt missing keyword %q", keyword)
		}
	}
}

func TestNewWithCustomClient(t *testing.T) {
	custom := &http.Client{Timeout: 30 * time.Second}
	opt := New(LLMConfig{Endpoint: "http://example.com", APIKey: "sk-test"}, custom)
	if opt == nil {
		t.Fatal("expected non-nil optimizer")
	}
	if opt.client != custom {
		t.Error("expected custom client to be used")
	}
	if opt.client.Timeout != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", opt.client.Timeout)
	}
}

func TestNewWithNilClient(t *testing.T) {
	opt := New(LLMConfig{Endpoint: "http://example.com", APIKey: "sk-test"}, nil)
	if opt == nil {
		t.Fatal("expected non-nil optimizer")
	}
	if opt.client.Timeout != 60*time.Second {
		t.Errorf("timeout = %v, want 60s (default)", opt.client.Timeout)
	}
}

func TestNewDefaultClientTimeout(t *testing.T) {
	opt := New(LLMConfig{Endpoint: "http://example.com", APIKey: "sk-test"})
	if opt == nil {
		t.Fatal("expected non-nil optimizer")
	}
	if opt.client.Timeout != 60*time.Second {
		t.Errorf("timeout = %v, want 60s", opt.client.Timeout)
	}
}

func TestConcatRawFormat(t *testing.T) {
	inputs := []ContentInput{
		{Source: "mysrc", Path: "mypath.md", Content: "my content"},
	}
	result := ConcatRaw(inputs)

	wantHeader := "# Instructions from mysrc (mypath.md)"
	if !strings.Contains(result, wantHeader) {
		t.Errorf("missing header %q in:\n%s", wantHeader, result)
	}
	if !strings.Contains(result, "---") {
		t.Error("missing separator ---")
	}
	// Verify structure: header, blank line, content, blank line, separator
	lines := strings.Split(result, "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines, got %d", len(lines))
	}
	if lines[0] != wantHeader {
		t.Errorf("first line = %q, want %q", lines[0], wantHeader)
	}
	if lines[1] != "" {
		t.Errorf("second line should be blank, got %q", lines[1])
	}
	if lines[2] != "my content" {
		t.Errorf("third line = %q, want %q", lines[2], "my content")
	}
}

// ---------------------------------------------------------------------------
// Additional nominal / error / limit tests
// ---------------------------------------------------------------------------

func TestConsolidate_SingleInput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		json.Unmarshal(body, &req)

		// Verify single input is still sent to LLM
		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages (system + user), got %d", len(req.Messages))
		}

		json.NewEncoder(w).Encode(chatResponse{
			Choices: []chatChoice{{Message: chatMessage{Content: "single result"}}},
		})
	}))
	defer srv.Close()

	opt := New(LLMConfig{
		Endpoint: srv.URL,
		APIKey:   "test-key",
		Model:    "test-model",
	})

	got, err := opt.Optimize(context.Background(), []ContentInput{{Source: "s1", Content: "single"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "single result" {
		t.Errorf("got %q, want %q", got, "single result")
	}
}

func TestConsolidate_InvalidEndpoint(t *testing.T) {
	opt := New(LLMConfig{
		Endpoint: "http://127.0.0.1:1", // unreachable port
		APIKey:   "test-key",
		Model:    "test-model",
	})

	_, err := opt.Optimize(context.Background(), []ContentInput{{Source: "s1", Content: "data"}})
	if err == nil {
		t.Fatal("expected error for unreachable endpoint")
	}
}

func TestConsolidate_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // slow server
	}))
	defer srv.Close()

	opt := New(LLMConfig{
		Endpoint: srv.URL,
		APIKey:   "test-key",
		Model:    "test-model",
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // immediately cancel

	_, err := opt.Optimize(ctx, []ContentInput{{Source: "s1", Content: "data"}})
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestConsolidate_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is not json"))
	}))
	defer srv.Close()

	opt := New(LLMConfig{
		Endpoint: srv.URL,
		APIKey:   "test-key",
		Model:    "test-model",
	})

	_, err := opt.Optimize(context.Background(), []ContentInput{{Source: "s1", Content: "data"}})
	if err == nil {
		t.Fatal("expected error for malformed JSON response")
	}
}

func TestConsolidate_EmptyChoices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(chatResponse{Choices: []chatChoice{}})
	}))
	defer srv.Close()

	opt := New(LLMConfig{
		Endpoint: srv.URL,
		APIKey:   "test-key",
		Model:    "test-model",
	})

	_, err := opt.Optimize(context.Background(), []ContentInput{{Source: "s1", Content: "data"}})
	if err == nil {
		t.Fatal("expected error for empty choices")
	}
}

func TestNew_MultipleClients(t *testing.T) {
	c1 := &http.Client{Timeout: 1 * time.Second}
	c2 := &http.Client{Timeout: 2 * time.Second}

	// Only first non-nil client should be used
	opt := New(LLMConfig{Endpoint: "http://example.com", APIKey: "key"}, c1, c2)
	if opt.client != c1 {
		t.Error("should use first client, not second")
	}
}

func TestEnabled(t *testing.T) {
	opt := New(LLMConfig{Endpoint: "http://example.com", APIKey: "key"})
	if !opt.Enabled() {
		t.Error("configured optimizer should be enabled")
	}
}
