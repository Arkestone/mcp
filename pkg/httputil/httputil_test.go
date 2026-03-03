package httputil

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// --------------- NewTransport tests ---------------

func TestNewTransport_Default(t *testing.T) {
	t.Parallel()
	tr, err := NewTransport(ProxyConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil transport")
	}
	if tr.Proxy == nil {
		t.Error("expected default proxy function to be set")
	}
}

func TestNewTransport_WithProxyURL(t *testing.T) {
	t.Parallel()

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	defer backend.Close()

	tr, err := NewTransport(ProxyConfig{ProxyURL: backend.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Proxy == nil {
		t.Fatal("expected proxy function to be set")
	}

	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	resp, err := client.Get(backend.URL + "/ping")
	if err != nil {
		t.Fatalf("request through proxy transport failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestNewTransport_InvalidProxyURL(t *testing.T) {
	t.Parallel()
	_, err := NewTransport(ProxyConfig{ProxyURL: "://bad"})
	if err == nil {
		t.Fatal("expected error for invalid proxy URL")
	}
}

func TestNewTransport_WithValidCACertFile(t *testing.T) {
	t.Parallel()

	pemData := generateSelfSignedCACertPEM(t)
	certFile := filepath.Join(t.TempDir(), "ca.pem")
	if err := os.WriteFile(certFile, pemData, 0o600); err != nil {
		t.Fatal(err)
	}

	tr, err := NewTransport(ProxyConfig{CACertFile: certFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.TLSClientConfig == nil {
		t.Fatal("expected TLSClientConfig to be set")
	}
	if tr.TLSClientConfig.RootCAs == nil {
		t.Fatal("expected RootCAs to be set")
	}
}

func TestNewTransport_NonexistentCACertFile(t *testing.T) {
	t.Parallel()
	_, err := NewTransport(ProxyConfig{CACertFile: "/no/such/file.pem"})
	if err == nil {
		t.Fatal("expected error for nonexistent CA cert file")
	}
}

func TestNewTransport_InvalidPEMData(t *testing.T) {
	t.Parallel()
	certFile := filepath.Join(t.TempDir(), "bad.pem")
	if err := os.WriteFile(certFile, []byte("not a cert"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := NewTransport(ProxyConfig{CACertFile: certFile})
	if err == nil {
		t.Fatal("expected error for invalid PEM data")
	}
}

func TestNewTransport_TLSInsecureSkipVerify(t *testing.T) {
	t.Parallel()
	tr, err := NewTransport(ProxyConfig{TLSInsecureSkipVerify: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.TLSClientConfig == nil {
		t.Fatal("expected TLSClientConfig to be set")
	}
	if !tr.TLSClientConfig.InsecureSkipVerify {
		t.Error("expected InsecureSkipVerify to be true")
	}
}

// --------------- NewClient tests ---------------

func TestNewClient_Timeout(t *testing.T) {
	t.Parallel()
	c, err := NewClient(ProxyConfig{}, 42*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Timeout != 42*time.Second {
		t.Errorf("expected timeout 42s, got %v", c.Timeout)
	}
}

func TestNewClient_TransportIsHeaderPropagator(t *testing.T) {
	t.Parallel()
	c, err := NewClient(ProxyConfig{}, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := c.Transport.(*headerPropagator); !ok {
		t.Errorf("expected transport to be *headerPropagator, got %T", c.Transport)
	}
}

// --------------- Header context helpers ---------------

func TestWithHeaders_RoundTrip(t *testing.T) {
	t.Parallel()
	h := http.Header{"X-Foo": {"bar"}}
	ctx := WithHeaders(context.Background(), h)
	got := HeadersFromContext(ctx)
	if got.Get("X-Foo") != "bar" {
		t.Errorf("expected X-Foo=bar, got %q", got.Get("X-Foo"))
	}
}

func TestHeadersFromContext_NoHeaders(t *testing.T) {
	t.Parallel()
	got := HeadersFromContext(context.Background())
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestWithHeaders_LastOneWins(t *testing.T) {
	t.Parallel()
	h1 := http.Header{"X-Foo": {"first"}}
	h2 := http.Header{"X-Foo": {"second"}}
	ctx := WithHeaders(context.Background(), h1)
	ctx = WithHeaders(ctx, h2)
	got := HeadersFromContext(ctx)
	if got.Get("X-Foo") != "second" {
		t.Errorf("expected last headers to win, got %q", got.Get("X-Foo"))
	}
}

// --------------- headerPropagator tests ---------------

func TestHeaderPropagator_AllHeaders_EmptyAllowList(t *testing.T) {
	t.Parallel()

	var received http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Clone()
	}))
	defer srv.Close()

	prop := &headerPropagator{
		base:      &http.Transport{},
		allowList: nil,
	}
	ctxHeaders := http.Header{
		"X-Custom-One": {"v1"},
		"X-Custom-Two": {"v2"},
	}
	ctx := WithHeaders(context.Background(), ctxHeaders)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)

	resp, err := prop.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	resp.Body.Close()

	if received.Get("X-Custom-One") != "v1" {
		t.Errorf("expected X-Custom-One=v1, got %q", received.Get("X-Custom-One"))
	}
	if received.Get("X-Custom-Two") != "v2" {
		t.Errorf("expected X-Custom-Two=v2, got %q", received.Get("X-Custom-Two"))
	}
}

func TestHeaderPropagator_OnlyAllowListed(t *testing.T) {
	t.Parallel()

	var received http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Clone()
	}))
	defer srv.Close()

	prop := &headerPropagator{
		base:      &http.Transport{},
		allowList: []string{"X-Allowed"},
	}
	ctxHeaders := http.Header{
		"X-Allowed":    {"yes"},
		"X-Disallowed": {"no"},
	}
	ctx := WithHeaders(context.Background(), ctxHeaders)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)

	resp, err := prop.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	resp.Body.Close()

	if received.Get("X-Allowed") != "yes" {
		t.Errorf("expected X-Allowed=yes, got %q", received.Get("X-Allowed"))
	}
	if received.Get("X-Disallowed") != "" {
		t.Errorf("expected X-Disallowed to be absent, got %q", received.Get("X-Disallowed"))
	}
}

func TestHeaderPropagator_DoesNotOverrideExplicit(t *testing.T) {
	t.Parallel()

	var received http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = r.Header.Clone()
	}))
	defer srv.Close()

	prop := &headerPropagator{
		base:      &http.Transport{},
		allowList: nil,
	}
	ctxHeaders := http.Header{"X-Auth": {"from-context"}}
	ctx := WithHeaders(context.Background(), ctxHeaders)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	req.Header.Set("X-Auth", "explicit")

	resp, err := prop.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	resp.Body.Close()

	if received.Get("X-Auth") != "explicit" {
		t.Errorf("expected explicitly-set header to win, got %q", received.Get("X-Auth"))
	}
}

func TestHeaderPropagator_NoContextHeaders(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	defer srv.Close()

	prop := &headerPropagator{
		base:      &http.Transport{},
		allowList: nil,
	}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL, nil)
	resp, err := prop.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok" {
		t.Errorf("expected ok, got %q", body)
	}
}

func TestHeaderPropagator_PropagatedHeadersReachServer(t *testing.T) {
	t.Parallel()

	gotCh := make(chan http.Header, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotCh <- r.Header.Clone()
	}))
	defer srv.Close()

	prop := &headerPropagator{
		base:      &http.Transport{},
		allowList: nil,
	}
	ctxHeaders := http.Header{
		"X-Trace-Id":   {"abc123"},
		"X-Session-Id": {"sess-1"},
	}
	ctx := WithHeaders(context.Background(), ctxHeaders)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)

	resp, err := prop.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip error: %v", err)
	}
	resp.Body.Close()

	select {
	case got := <-gotCh:
		if got.Get("X-Trace-Id") != "abc123" {
			t.Errorf("expected X-Trace-Id=abc123, got %q", got.Get("X-Trace-Id"))
		}
		if got.Get("X-Session-Id") != "sess-1" {
			t.Errorf("expected X-Session-Id=sess-1, got %q", got.Get("X-Session-Id"))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for headers")
	}
}

// --------------- HeaderCaptureMiddleware tests ---------------

func TestHeaderCaptureMiddleware_CapturesAll_EmptyAllowList(t *testing.T) {
	t.Parallel()

	var captured http.Header
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = HeadersFromContext(r.Context())
	})
	mw := HeaderCaptureMiddleware(nil, inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-One", "1")
	req.Header.Set("X-Two", "2")
	mw.ServeHTTP(httptest.NewRecorder(), req)

	if captured.Get("X-One") != "1" {
		t.Errorf("expected X-One=1, got %q", captured.Get("X-One"))
	}
	if captured.Get("X-Two") != "2" {
		t.Errorf("expected X-Two=2, got %q", captured.Get("X-Two"))
	}
}

func TestHeaderCaptureMiddleware_CapturesOnlyAllowListed(t *testing.T) {
	t.Parallel()

	var captured http.Header
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = HeadersFromContext(r.Context())
	})
	mw := HeaderCaptureMiddleware([]string{"X-Keep"}, inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Keep", "yes")
	req.Header.Set("X-Drop", "no")
	mw.ServeHTTP(httptest.NewRecorder(), req)

	if captured.Get("X-Keep") != "yes" {
		t.Errorf("expected X-Keep=yes, got %q", captured.Get("X-Keep"))
	}
	if captured.Get("X-Drop") != "" {
		t.Errorf("expected X-Drop to be absent, got %q", captured.Get("X-Drop"))
	}
}

func TestHeaderCaptureMiddleware_AvailableDownstream(t *testing.T) {
	t.Parallel()

	done := make(chan struct{})
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := HeadersFromContext(r.Context())
		if h == nil {
			t.Error("expected headers in context, got nil")
		} else if h.Get("X-Request-Id") != "req-42" {
			t.Errorf("expected X-Request-Id=req-42, got %q", h.Get("X-Request-Id"))
		}
		close(done)
	})
	mw := HeaderCaptureMiddleware(nil, inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "req-42")
	mw.ServeHTTP(httptest.NewRecorder(), req)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not run")
	}
}

func TestHeaderCaptureMiddleware_ExcludesNonAllowListed(t *testing.T) {
	t.Parallel()

	var captured http.Header
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = HeadersFromContext(r.Context())
	})
	mw := HeaderCaptureMiddleware([]string{"X-A"}, inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-A", "a")
	req.Header.Set("X-B", "b")
	req.Header.Set("X-C", "c")
	mw.ServeHTTP(httptest.NewRecorder(), req)

	if captured.Get("X-A") != "a" {
		t.Errorf("expected X-A=a, got %q", captured.Get("X-A"))
	}
	for _, k := range []string{"X-B", "X-C"} {
		if captured.Get(k) != "" {
			t.Errorf("expected %s to be excluded, got %q", k, captured.Get(k))
		}
	}
}

// --------------- Integration: end-to-end header propagation ---------------

func TestIntegration_HeaderPropagation_EndToEnd(t *testing.T) {
	t.Parallel()

	// 1. "GitHub" mock — records headers received on outbound calls.
	outboundHeaders := make(chan http.Header, 1)
	githubMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outboundHeaders <- r.Header.Clone()
		fmt.Fprint(w, `{"status":"ok"}`)
	}))
	defer githubMock.Close()

	// 2. Build an outbound client that propagates X-Request-Id.
	outboundClient, err := NewClient(ProxyConfig{
		HeaderPassthrough: []string{"X-Request-Id"},
	}, 5*time.Second)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// 3. "MCP tool handler" — reads context headers, makes outbound call.
	toolHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		outReq, _ := http.NewRequestWithContext(r.Context(), http.MethodGet, githubMock.URL+"/api", nil)
		resp, reqErr := outboundClient.Do(outReq)
		if reqErr != nil {
			http.Error(w, reqErr.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	})

	// 4. Wrap with capture middleware (allow only X-Request-Id).
	mcpServer := httptest.NewServer(HeaderCaptureMiddleware(
		[]string{"X-Request-Id"}, toolHandler,
	))
	defer mcpServer.Close()

	// 5. Simulate an inbound MCP request with X-Request-Id.
	req, _ := http.NewRequest(http.MethodGet, mcpServer.URL+"/tool/call", nil)
	req.Header.Set("X-Request-Id", "trace-999")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("inbound request failed: %v", err)
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	// 6. Verify the outbound call to "GitHub" carried the header.
	select {
	case got := <-outboundHeaders:
		if got.Get("X-Request-Id") != "trace-999" {
			t.Errorf("expected outbound X-Request-Id=trace-999, got %q", got.Get("X-Request-Id"))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for outbound request")
	}
}

// --------------- helpers ---------------

// generateSelfSignedCACertPEM creates a minimal self-signed CA certificate PEM.
func generateSelfSignedCACertPEM(t *testing.T) []byte {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
}

func TestDefaultClient(t *testing.T) {
	c := DefaultClient(5 * time.Second)
	if c == nil {
		t.Fatal("DefaultClient returned nil")
	}
	if c.Timeout != 5*time.Second {
		t.Errorf("timeout = %v, want 5s", c.Timeout)
	}
}

func TestNewClient_InvalidProxy(t *testing.T) {
	_, err := NewClient(ProxyConfig{ProxyURL: "://bad"}, time.Second)
	if err == nil {
		t.Fatal("expected error for invalid proxy URL")
	}
}

// ---------------------------------------------------------------------------
// Additional nominal / error / limit tests
// ---------------------------------------------------------------------------

func TestHeaderCaptureMiddleware_NoMatchingHeaders(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdrs := HeadersFromContext(r.Context())
		// Only headers in allow list should be captured
		if hdrs.Get("Accept") != "" {
			t.Error("Accept should not be captured with specific allowList")
		}
		w.WriteHeader(200)
	})

	handler := HeaderCaptureMiddleware([]string{"X-Custom", "Authorization"}, inner)
	req := httptest.NewRequest("GET", "/", nil)
	// No X-Custom or Authorization headers, only Accept
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
}

func TestHeaderCaptureMiddleware_MultiValueHeaders(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdrs := HeadersFromContext(r.Context())
		vals := hdrs.Values("X-Multi")
		if len(vals) != 2 {
			t.Errorf("expected 2 values for X-Multi, got %d: %v", len(vals), vals)
		}
		w.WriteHeader(200)
	})

	handler := HeaderCaptureMiddleware([]string{"X-Multi"}, inner)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Multi", "first")
	req.Header.Add("X-Multi", "second")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
}

func TestHeaderCaptureMiddleware_EmptyAllowList(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdrs := HeadersFromContext(r.Context())
		// With empty allowList, ALL headers should be captured
		if hdrs.Get("Accept") == "" {
			t.Error("Accept should be captured with empty allowList")
		}
		w.WriteHeader(200)
	})

	handler := HeaderCaptureMiddleware(nil, inner)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
}

func TestHeaderPropagator_PreservesExistingHeaders(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer existing" {
			t.Errorf("Authorization = %q, want Bearer existing", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Custom") != "injected" {
			t.Errorf("X-Custom = %q, want injected", r.Header.Get("X-Custom"))
		}
		w.WriteHeader(200)
	}))
	defer backend.Close()

	captured := http.Header{}
	captured.Set("Authorization", "Bearer override-attempt")
	captured.Set("X-Custom", "injected")

	ctx := WithHeaders(context.Background(), captured)

	tr := &headerPropagator{
		base:      &http.Transport{},
		allowList: []string{"Authorization", "X-Custom"},
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", backend.URL, nil)
	req.Header.Set("Authorization", "Bearer existing") // should NOT be overridden

	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip: %v", err)
	}
	resp.Body.Close()
}

func TestNewTransport_CACertFromFile(t *testing.T) {
	// Create a valid PEM cert file
	certFile := filepath.Join(t.TempDir(), "ca.pem")
	cert, key := generateSelfSignedCert(t)
	_ = key
	os.WriteFile(certFile, cert, 0o644)

	tr, err := NewTransport(ProxyConfig{CACertFile: certFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.TLSClientConfig == nil {
		t.Fatal("expected TLS config with custom CA")
	}
}

func TestNewTransport_InvalidCACertFile(t *testing.T) {
	_, err := NewTransport(ProxyConfig{CACertFile: "/nonexistent/ca.pem"})
	if err == nil {
		t.Fatal("expected error for nonexistent CA cert file")
	}
}

func TestNewTransport_InvalidCACertContent(t *testing.T) {
	certFile := filepath.Join(t.TempDir(), "bad.pem")
	os.WriteFile(certFile, []byte("not a valid cert"), 0o644)

	_, err := NewTransport(ProxyConfig{CACertFile: certFile})
	if err == nil {
		t.Fatal("expected error for invalid cert content")
	}
}

func TestNewClient_WithHeaderPassthrough(t *testing.T) {
	c, err := NewClient(ProxyConfig{HeaderPassthrough: []string{"X-Custom"}}, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Transport should be a headerPropagator
	_, ok := c.Transport.(*headerPropagator)
	if !ok {
		t.Errorf("transport type = %T, want *headerPropagator", c.Transport)
	}
}

func TestNewClient_WithoutHeaderPassthrough(t *testing.T) {
	c, err := NewClient(ProxyConfig{HeaderPassthrough: nil}, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// NewClient always wraps with headerPropagator; empty allowList forwards all
	hp, ok := c.Transport.(*headerPropagator)
	if !ok {
		t.Fatalf("transport type = %T, want *headerPropagator", c.Transport)
	}
	if len(hp.allowList) != 0 {
		t.Errorf("allowList = %v, want empty", hp.allowList)
	}
}

// generateSelfSignedCert creates a self-signed cert+key pair for testing.
func generateSelfSignedCert(t *testing.T) (certPEM, keyPEM []byte) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		IsCA:         true,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	return
}
