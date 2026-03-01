package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/httputil"
	"github.com/Arkestone/mcp/pkg/optimizer"
	"github.com/Arkestone/mcp/pkg/testutil"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestShouldOptimize(t *testing.T) {
	enabled := optimizer.New(testutil.LLMConfig())
	if enabled == nil {
		t.Fatal("expected non-nil optimizer from testutil.LLMConfig")
	}

	tests := []struct {
		name          string
		opt           *optimizer.Optimizer
		globalDefault bool
		perRequest    string
		want          bool
	}{
		// nil optimizer → always false
		{"nil opt, default false, empty override", nil, false, "", false},
		{"nil opt, default true, empty override", nil, true, "", false},
		{"nil opt, default true, override true", nil, true, "true", false},

		// per-request true variants
		{"enabled, override true", enabled, false, "true", true},
		{"enabled, override TRUE", enabled, false, "TRUE", true},
		{"enabled, override True", enabled, false, "True", true},
		{"enabled, override 1", enabled, false, "1", true},
		{"enabled, override yes", enabled, false, "yes", true},
		{"enabled, override YES", enabled, false, "YES", true},

		// per-request false variants
		{"enabled, override false", enabled, true, "false", false},
		{"enabled, override FALSE", enabled, true, "FALSE", false},
		{"enabled, override False", enabled, true, "False", false},
		{"enabled, override 0", enabled, true, "0", false},
		{"enabled, override no", enabled, true, "no", false},
		{"enabled, override NO", enabled, true, "NO", false},

		// empty/unknown perRequest → globalDefault
		{"enabled, default true, empty override", enabled, true, "", true},
		{"enabled, default false, empty override", enabled, false, "", false},
		{"enabled, default true, unknown override", enabled, true, "maybe", true},
		{"enabled, default false, unknown override", enabled, false, "maybe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldOptimize(tt.opt, tt.globalDefault, tt.perRequest)
			if got != tt.want {
				t.Errorf("ShouldOptimize(%v, %v, %q) = %v, want %v",
					tt.opt != nil, tt.globalDefault, tt.perRequest, got, tt.want)
			}
		})
	}
}

func TestWrapHandler(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if headers were captured into context
		hdrs := httputil.HeadersFromContext(r.Context())
		if hdrs != nil && hdrs.Get("X-Test") != "" {
			w.Header().Set("X-Captured", hdrs.Get("X-Test"))
		}
		w.WriteHeader(200)
	})

	t.Run("no passthrough returns unwrapped handler", func(t *testing.T) {
		wrapped := WrapHandler(inner, nil)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Test", "hello")
		wrapped.ServeHTTP(rec, req)
		// No middleware → headers not captured into context
		if rec.Header().Get("X-Captured") != "" {
			t.Error("expected no captured header without middleware")
		}
	})

	t.Run("with passthrough captures headers", func(t *testing.T) {
		wrapped := WrapHandler(inner, []string{"X-Test"})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Test", "hello")
		wrapped.ServeHTTP(rec, req)
		if got := rec.Header().Get("X-Captured"); got != "hello" {
			t.Errorf("X-Captured = %q, want %q", got, "hello")
		}
	})

	t.Run("empty slice returns unwrapped handler", func(t *testing.T) {
		wrapped := WrapHandler(inner, []string{})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Test", "hello")
		wrapped.ServeHTTP(rec, req)
		if rec.Header().Get("X-Captured") != "" {
			t.Error("expected no captured header with empty slice")
		}
	})
}

func TestRunHTTP(t *testing.T) {
	t.Run("starts and stops cleanly", func(t *testing.T) {
		srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
		ctx, cancel := context.WithCancel(context.Background())

		errCh := make(chan error, 1)
		go func() {
			errCh <- RunHTTP(ctx, srv, "127.0.0.1:19847", nil)
		}()

		time.Sleep(100 * time.Millisecond)
		cancel()

		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("RunHTTP returned error: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("RunHTTP didn't stop after context cancel")
		}
	})

	t.Run("with header passthrough", func(t *testing.T) {
		srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
		ctx, cancel := context.WithCancel(context.Background())

		errCh := make(chan error, 1)
		go func() {
			errCh <- RunHTTP(ctx, srv, "127.0.0.1:19848", []string{"X-Request-Id", "X-Trace"})
		}()

		time.Sleep(100 * time.Millisecond)
		cancel()

		select {
		case err := <-errCh:
			if err != nil {
				t.Fatalf("RunHTTP returned error: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("RunHTTP didn't stop after context cancel")
		}
	})

	t.Run("returns error on invalid addr", func(t *testing.T) {
		srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
		ctx := context.Background()
		err := RunHTTP(ctx, srv, "127.0.0.1:-1", nil)
		if err == nil {
			t.Fatal("expected error for invalid address")
		}
	})
}

// ---------------------------------------------------------------------------
// Additional nominal / error / limit tests
// ---------------------------------------------------------------------------

func TestWrapHandler_NilProxy(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	wrapped := WrapHandler(handler, nil)
	if wrapped == nil {
		t.Fatal("WrapHandler should return non-nil handler even with nil headerPassthrough")
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	wrapped.ServeHTTP(rec, req)
}

func TestWrapHandler_WithHeaders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	wrapped := WrapHandler(handler, []string{"X-Custom"})
	if wrapped == nil {
		t.Fatal("WrapHandler should return non-nil handler")
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Custom", "value")
	wrapped.ServeHTTP(rec, req)
}

func TestShouldOptimize_AllCombinations(t *testing.T) {
	enabled := optimizer.New(testutil.LLMConfig())
	tests := []struct {
		name      string
		opt       *optimizer.Optimizer
		globalDef bool
		perReq    string
		want      bool
	}{
		{"nil_false_empty", nil, false, "", false},
		{"nil_true_empty", nil, true, "", false},
		{"nil_false_true", nil, false, "true", false},
		{"enabled_false_empty", enabled, false, "", false},
		{"enabled_true_empty", enabled, true, "", true},
		{"enabled_false_true", enabled, false, "true", true},
		{"enabled_true_false", enabled, true, "false", false},
		{"enabled_false_false", enabled, false, "false", false},
		{"enabled_false_1", enabled, false, "1", true},
		{"enabled_true_0", enabled, true, "0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ShouldOptimize(tt.opt, tt.globalDef, tt.perReq)
			if got != tt.want {
				t.Errorf("ShouldOptimize(%v, %v, %q) = %v, want %v", tt.opt != nil, tt.globalDef, tt.perReq, got, tt.want)
			}
		})
	}
}

func TestRunHTTP_AddressAlreadyInUse(t *testing.T) {
	// First bind a port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	addr := l.Addr().String()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	ctx := context.Background()
	err = RunHTTP(ctx, srv, addr, nil)
	if err == nil {
		t.Fatal("expected error for address already in use")
	}
}

func TestRunHTTP_Healthz(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	addr := l.Addr().String()
	l.Close()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go RunHTTP(ctx, srv, addr, nil) //nolint:errcheck

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf(`body["status"] = %q, want "ok"`, body["status"])
	}
}

func TestRunHTTP_GracefulShutdown(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	addr := l.Addr().String()
	l.Close()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- RunHTTP(ctx, srv, addr, nil)
	}()

	time.Sleep(100 * time.Millisecond)

	// Verify the server is up.
	resp, err := http.Get("http://" + addr + "/healthz")
	if err != nil {
		t.Fatalf("server not up before shutdown: %v", err)
	}
	resp.Body.Close()

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("RunHTTP returned error after graceful shutdown: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("RunHTTP did not stop after context cancel")
	}
}

func TestRunHTTP_ServesHTTPRequests(t *testing.T) {
	// Find a free port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	addr := l.Addr().String()
	l.Close()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- RunHTTP(ctx, srv, addr, nil)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Make a real HTTP request
	resp, err := http.Post("http://"+addr+"/", "application/json", nil)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Any HTTP response means the server is actually serving
	if resp.StatusCode == 0 {
		t.Error("expected a valid HTTP status code")
	}

	cancel()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("RunHTTP returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("RunHTTP did not stop after context cancel")
	}
}
