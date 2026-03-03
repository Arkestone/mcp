// Package httputil provides shared HTTP infrastructure for proxy support,
// custom TLS, and header propagation across all MCP servers.
//
// All outbound HTTP clients (GitHub API, LLM endpoints) use this package
// so that proxy, CA certificate, and header-forwarding settings apply uniformly.
package httputil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

// ProxyConfig holds network settings shared across all HTTP clients.
type ProxyConfig struct {
	// ProxyURL overrides the proxy. When empty, Go's default behavior applies:
	// HTTP_PROXY / HTTPS_PROXY / NO_PROXY environment variables are honored.
	ProxyURL string `yaml:"proxy_url"`

	// CACertFile is a PEM-encoded CA bundle appended to the system pool.
	// Typically used for corporate TLS-intercepting proxies.
	CACertFile string `yaml:"ca_cert"`

	// TLSInsecureSkipVerify disables certificate verification. Use for testing only.
	TLSInsecureSkipVerify bool `yaml:"tls_insecure_skip_verify"`

	// HeaderPassthrough lists HTTP header names that should be forwarded from
	// incoming MCP requests to outbound calls (GitHub, LLM). Case-insensitive.
	// When empty, no headers are forwarded.
	HeaderPassthrough []string `yaml:"header_passthrough"`
}

// NewTransport creates an *http.Transport configured with proxy, TLS, and
// header-propagation settings. The returned transport automatically forwards
// headers stored in the request context (see WithHeaders / HeadersFromContext).
func NewTransport(cfg ProxyConfig) (*http.Transport, error) {
	base := http.DefaultTransport.(*http.Transport).Clone()

	// Proxy
	if cfg.ProxyURL != "" {
		u, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL %q: %w", cfg.ProxyURL, err)
		}
		base.Proxy = http.ProxyURL(u)
	}
	// else: default proxy func reads HTTP_PROXY / HTTPS_PROXY / NO_PROXY

	// TLS
	if cfg.CACertFile != "" || cfg.TLSInsecureSkipVerify {
		tlsCfg := &tls.Config{
			InsecureSkipVerify: cfg.TLSInsecureSkipVerify, //nolint:gosec // operator opt-in
		}
		if cfg.CACertFile != "" {
			pem, err := os.ReadFile(cfg.CACertFile)
			if err != nil {
				return nil, fmt.Errorf("reading CA cert %q: %w", cfg.CACertFile, err)
			}
			pool, err := x509.SystemCertPool()
			if err != nil {
				pool = x509.NewCertPool()
			}
			if !pool.AppendCertsFromPEM(pem) {
				return nil, fmt.Errorf("no valid certificates in %q", cfg.CACertFile)
			}
			tlsCfg.RootCAs = pool
		}
		base.TLSClientConfig = tlsCfg
	}

	return base, nil
}

// NewClient creates an *http.Client with proxy/TLS settings and the given timeout.
// Outbound requests automatically propagate context headers (see WithHeaders).
func NewClient(cfg ProxyConfig, timeout time.Duration) (*http.Client, error) {
	transport, err := NewTransport(cfg)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &headerPropagator{
			base:      transport,
			allowList: cfg.HeaderPassthrough,
		},
		Timeout: timeout,
	}, nil
}

// DefaultClient creates a minimal client that still honors env-based proxy
// and propagates context headers. Use when no explicit ProxyConfig is needed.
func DefaultClient(timeout time.Duration) *http.Client {
	c, _ := NewClient(ProxyConfig{}, timeout)
	return c
}

// ---------- context-based header propagation ----------

type ctxKey struct{}

// WithHeaders stores HTTP headers in context for downstream propagation.
func WithHeaders(ctx context.Context, h http.Header) context.Context {
	return context.WithValue(ctx, ctxKey{}, h)
}

// HeadersFromContext retrieves previously stored headers. Returns nil if none.
func HeadersFromContext(ctx context.Context) http.Header {
	if h, ok := ctx.Value(ctxKey{}).(http.Header); ok {
		return h
	}
	return nil
}

// headerPropagator is an http.RoundTripper that copies context headers to
// outbound requests before delegating to the base transport.
type headerPropagator struct {
	base      http.RoundTripper
	allowList []string // if empty, forward all context headers
}

func (p *headerPropagator) RoundTrip(req *http.Request) (*http.Response, error) {
	incoming := HeadersFromContext(req.Context())
	if incoming != nil {
		// Clone the request so we don't mutate the caller's headers.
		clone := req.Clone(req.Context())
		if len(p.allowList) == 0 {
			for k, vals := range incoming {
				if clone.Header.Get(k) == "" { // don't override explicitly set headers
					for _, v := range vals {
						clone.Header.Add(k, v)
					}
				}
			}
		} else {
			for _, k := range p.allowList {
				if vals := incoming.Values(k); len(vals) > 0 && clone.Header.Get(k) == "" {
					for _, v := range vals {
						clone.Header.Add(k, v)
					}
				}
			}
		}
		req = clone
	}
	return p.base.RoundTrip(req)
}

// ---------- header-capture middleware (for incoming HTTP transport) ----------

// HeaderCaptureMiddleware returns an http.Handler that captures selected headers
// from inbound requests into the request context so that downstream MCP tool
// handlers can propagate them to outbound calls.
//
// When allowList is empty, all headers are captured.
func HeaderCaptureMiddleware(allowList []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured := make(http.Header)
		if len(allowList) == 0 {
			for k, vals := range r.Header {
				captured[k] = vals
			}
		} else {
			for _, k := range allowList {
				if vals := r.Header.Values(k); len(vals) > 0 {
					captured[k] = vals
				}
			}
		}
		ctx := WithHeaders(r.Context(), captured)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ---------- helpers ----------

// DialTimeout is a sensible default for custom transports.
const DialTimeout = 30 * time.Second

func init() {
	// Ensure default transport has reasonable timeouts for proxied environments.
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		if t.DialContext == nil {
			t.DialContext = (&net.Dialer{Timeout: DialTimeout}).DialContext
		}
	}
}
