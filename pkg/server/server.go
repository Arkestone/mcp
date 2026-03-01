// Package server provides shared MCP server utilities.
package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Arkestone/mcp/pkg/httputil"
	"github.com/Arkestone/mcp/pkg/optimizer"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ShouldOptimize determines whether to run LLM optimization for a request.
// It considers the optimizer availability, the global default, and an
// optional per-request override ("true"/"false"/"yes"/"no"/"1"/"0").
func ShouldOptimize(opt *optimizer.Optimizer, globalDefault bool, perRequest string) bool {
	if !opt.Enabled() {
		return false
	}
	switch strings.ToLower(perRequest) {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return globalDefault
	}
}

// WrapHandler optionally wraps an http.Handler with header-capture middleware.
// When headerPassthrough is non-empty, the returned handler captures the listed
// headers from inbound requests into the request context.
func WrapHandler(handler http.Handler, headerPassthrough []string) http.Handler {
	if len(headerPassthrough) > 0 {
		return httputil.HeaderCaptureMiddleware(headerPassthrough, handler)
	}
	return handler
}

// RunHTTP starts a StreamableHTTP MCP server and blocks until ctx is canceled.
// When headerPassthrough is non-empty, listed headers from inbound HTTP requests
// are captured into the request context so downstream calls can forward them.
// A GET /healthz endpoint returns 200 {"status":"ok"} and is always available.
func RunHTTP(ctx context.Context, srv *mcp.Server, addr string, headerPassthrough []string) error {
	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server { return srv },
		nil,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.Handle("/", mcpHandler)

	root := WrapHandler(mux, headerPassthrough)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: root,
	}

	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutCtx)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
