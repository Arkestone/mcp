# Network Configuration Guide

This document describes the network architecture, outbound connections, proxy
support, TLS settings, and firewall rules for the **mcp-instructions** and
**mcp-skills** MCP servers. It is intended for infrastructure operators and
network engineers.

---

## 1. Network Architecture

Both MCP servers share the same network stack (`pkg/httputil`). They make
**outbound HTTPS calls only**—there are no peer-to-peer connections between the
two servers.

```
                      ┌──────────────────────────────┐
                      │        MCP Clients            │
                      │  (IDEs, CLI tools, agents)    │
                      └──────────┬───────────────────┘
                                 │
                   ┌─────────────┴─────────────┐
                   │  Inbound (HTTP mode only)  │
                   │  :8080 instructions        │
                   │  :8081 skills              │
                   └─────────────┬─────────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              │          MCP Server Process          │
              │  ┌────────────┐  ┌───────────────┐  │
              │  │ GitHub     │  │ LLM Optimizer  │  │
              │  │ Client     │  │ Client         │  │
              │  └─────┬──────┘  └──────┬─────────┘  │
              └────────┼────────────────┼────────────┘
                       │                │
            ┌──────────┴──┐     ┌───────┴──────────┐
            │  (optional) │     │    (optional)     │
            │  HTTP/HTTPS │     │    HTTP/HTTPS     │
            │  Proxy      │     │    Proxy          │
            └──────┬──────┘     └───────┬───────────┘
                   │                    │
          ┌────────┴───────┐   ┌────────┴──────────┐
          │ api.github.com │   │ LLM Endpoint      │
          │ or GHE Server  │   │ (OpenAI-compat.)  │
          │ :443 HTTPS     │   │ :443 HTTPS        │
          └────────────────┘   └───────────────────┘
```

**Transport modes:**

| Mode    | Inbound listener | Use case |
|---------|-----------------|----------|
| `stdio` | **None** — communicates over stdin/stdout | IDE integrations, single-client |
| `http`  | Listens on configurable address (default `:8080` / `:8081`) | Multi-client, remote access |

In `stdio` mode the server opens **no listening ports**. All network activity is
outbound only.

---

## 2. Outbound Connections

| Destination | Port | Protocol | Purpose | When used |
|-------------|------|----------|---------|-----------|
| `api.github.com` | 443 | HTTPS | GitHub Contents API — fetch instruction/skill files from repositories | When `sources.repos` is configured |
| GitHub Enterprise Server (custom URL) | 443 | HTTPS | Same as above, for GHE | When GitHub client `BaseURL` points to GHE |
| LLM endpoint (user-configured) | Varies (typically 443) | HTTPS | OpenAI-compatible chat completions — content optimization | When `llm.endpoint` and `LLM_API_KEY` are set |
| HTTP/HTTPS proxy | Varies (typically 8080 or 3128) | HTTP or HTTPS | Tunnels all outbound traffic through corporate proxy | When `proxy.proxy_url` or `HTTP_PROXY`/`HTTPS_PROXY` is set |

All outbound HTTP clients share the same transport configuration created by
`httputil.NewClient()`. This ensures proxy, TLS, and header-passthrough settings
apply uniformly to every outbound request.

**Default timeouts:**

| Client | Timeout |
|--------|---------|
| GitHub API client | 30 seconds |
| LLM optimizer client | 60 seconds |
| TCP dial timeout | 30 seconds |

---

## 3. Proxy Configuration

The servers support four methods of proxy configuration, applied in the
following priority order (highest priority wins):

### 3.1 CLI Flags (highest priority)

```bash
mcp-instructions --proxy-url http://proxy.corp.example:8080 \
                 --ca-cert /etc/ssl/corp-ca.pem
```

Available flags:

| Flag | Description |
|------|-------------|
| `--proxy-url` | HTTP/HTTPS proxy URL |
| `--ca-cert` | Path to PEM CA certificate bundle |

### 3.2 Environment Variables (per-server prefix)

Each server uses its own environment variable prefix:

| Server | Prefix |
|--------|--------|
| mcp-instructions | `INSTRUCTIONS` |
| mcp-skills | `SKILLS` |

| Variable | Description | Example |
|----------|-------------|---------|
| `<PREFIX>_PROXY_URL` | HTTP/HTTPS proxy URL | `INSTRUCTIONS_PROXY_URL=http://proxy:8080` |
| `<PREFIX>_CA_CERT` | Path to PEM CA certificate | `INSTRUCTIONS_CA_CERT=/etc/ssl/ca.pem` |
| `<PREFIX>_TLS_INSECURE_SKIP_VERIFY` | Skip TLS verification (`true`/`1`) | `INSTRUCTIONS_TLS_INSECURE_SKIP_VERIFY=false` |
| `<PREFIX>_HEADER_PASSTHROUGH` | Comma-separated header names to forward | `INSTRUCTIONS_HEADER_PASSTHROUGH=X-Request-ID,X-Correlation-ID` |

### 3.3 Config File (YAML)

```yaml
proxy:
  proxy_url: "http://proxy.corp.example:8080"
  ca_cert: "/etc/ssl/corp-ca.pem"
  tls_insecure_skip_verify: false
  header_passthrough:
    - X-Request-ID
    - X-Correlation-ID
```

### 3.4 Standard Go Proxy Environment Variables (lowest priority / fallback)

| Variable | Description |
|----------|-------------|
| `HTTP_PROXY` | Proxy for HTTP requests |
| `HTTPS_PROXY` | Proxy for HTTPS requests |
| `NO_PROXY` | Comma-separated list of hosts to bypass the proxy |

These standard variables are **always honoured** by Go's `net/http` default
transport as a fallback. However, when `proxy_url` is set via any of the
methods above (config file, env var, or CLI flag), it **takes precedence** over
`HTTP_PROXY` / `HTTPS_PROXY` for all outbound requests.

> **Precedence summary:** CLI flag → environment variable → config file →
> `HTTP_PROXY`/`HTTPS_PROXY`/`NO_PROXY`.

---

## 4. TLS / Certificate Configuration

### Custom CA Certificates

Corporate environments that perform TLS inspection (MITM proxies) require the
proxy's CA certificate to be trusted by the MCP server. Configure this with
the `ca_cert` setting:

```yaml
proxy:
  ca_cert: "/etc/ssl/certs/corp-proxy-ca.pem"
```

**How it works:**

1. The system certificate pool is loaded first.
2. The custom CA certificate(s) from the PEM file are **appended** to the pool.
3. All outbound HTTPS connections validate against the combined pool.

If the system pool cannot be loaded (e.g., in a minimal container), an empty
pool is created and only the custom CA certificates are used.

### PEM Bundle Format

The `ca_cert` file can contain multiple certificates. Concatenate your
intermediate and root CAs into a single PEM file:

```
-----BEGIN CERTIFICATE-----
(intermediate CA certificate)
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
(root CA certificate)
-----END CERTIFICATE-----
```

Create the bundle:

```bash
cat intermediate-ca.pem root-ca.pem > /etc/ssl/certs/corp-ca-bundle.pem
```

### Disabling TLS Verification

> ⚠️ **Warning:** `tls_insecure_skip_verify` disables **all** certificate
> verification for outbound HTTPS connections. This exposes traffic to
> man-in-the-middle attacks. **Use only in isolated test environments. Never
> enable in production.**

```yaml
proxy:
  tls_insecure_skip_verify: true   # TESTING ONLY
```

Or via environment variable:

```bash
export INSTRUCTIONS_TLS_INSECURE_SKIP_VERIFY=true   # TESTING ONLY
```

---

## 5. HTTP Header Pass-through

When running in `http` transport mode, the servers can forward selected HTTP
headers from incoming MCP client requests to all outbound calls (GitHub API,
LLM endpoint). This is useful for:

- **Correlation IDs** — trace a request across systems (`X-Request-ID`,
  `X-Correlation-ID`)
- **Audit trails** — pass user identity or session info to downstream APIs
- **Custom authentication** — forward additional auth tokens required by
  proxies or gateways
- **Observability** — propagate tracing headers (e.g., `traceparent`)

### Configuration

Specify the header names to forward in the `header_passthrough` list:

```yaml
proxy:
  header_passthrough:
    - X-Request-ID
    - X-Correlation-ID
    - traceparent
```

Or via environment variable (comma-separated):

```bash
export INSTRUCTIONS_HEADER_PASSTHROUGH=X-Request-ID,X-Correlation-ID,traceparent
```

### Behaviour

- Header matching is **case-insensitive**.
- Headers that are **explicitly set** on outbound requests (e.g.,
  `Authorization`, `Accept`, `Content-Type`) are **never overridden** by
  pass-through. Explicitly set headers always take precedence.
- When `header_passthrough` is empty, **no headers** are forwarded.
- This feature **only applies in `http` transport mode**. In `stdio` mode there
  are no incoming HTTP headers to capture.

### How It Works

1. `HeaderCaptureMiddleware` intercepts incoming HTTP requests and stores the
   listed headers in the Go request context.
2. The `headerPropagator` round-tripper reads headers from the context and
   copies them onto outbound requests before sending.

---

## 6. Firewall Rules

Provide the following table to your network or security team. All entries assume
default ports; adjust if your environment uses non-standard ports.

| Direction | Source | Destination | Port | Protocol | Required | Purpose |
|-----------|--------|-------------|------|----------|----------|---------|
| Inbound | MCP clients | MCP server | 8080 | HTTP | Only in `http` mode | mcp-instructions MCP protocol |
| Inbound | MCP clients | MCP server | 8081 | HTTP | Only in `http` mode | mcp-skills MCP protocol |
| Outbound | MCP server | `api.github.com` | 443 | HTTPS | When `sources.repos` configured | GitHub Contents API |
| Outbound | MCP server | GHE server hostname | 443 | HTTPS | When using GitHub Enterprise | GitHub Enterprise Contents API |
| Outbound | MCP server | LLM endpoint hostname | 443 (typical) | HTTPS | When LLM optimization enabled | OpenAI-compatible chat completions |
| Outbound | MCP server | Proxy server | Varies | HTTP/HTTPS | When proxy configured | All outbound traffic via proxy |

**Notes:**

- When a proxy is configured, **only the proxy destination needs to be
  reachable** from the MCP server. The proxy handles onward connectivity to
  GitHub and LLM endpoints.
- In `stdio` mode, no inbound firewall rules are needed.
- If `sources.repos` is empty (local directories only), no outbound GitHub
  access is required.
- The LLM endpoint is only contacted when `llm.endpoint` and `LLM_API_KEY` are
  both set.

---

## 7. Deployment Scenarios

### 7.1 Direct Internet Access

Minimal configuration — the server can reach GitHub and the LLM endpoint
directly.

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
mcp-instructions --repos github/awesome-copilot
```

No proxy or CA certificate configuration is needed.

### 7.2 Corporate Proxy

Route all outbound traffic through a corporate HTTP proxy with TLS inspection.

```yaml
# config.yaml
sources:
  repos:
    - github/awesome-copilot

proxy:
  proxy_url: "http://proxy.corp.example:8080"
  ca_cert: "/etc/ssl/certs/corp-ca-bundle.pem"
```

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
mcp-instructions --config config.yaml
```

### 7.3 Air-Gapped / No Internet

Use only local directories. No outbound network connections are made.

```yaml
sources:
  dirs:
    - /opt/instructions/repo-a
    - /opt/instructions/repo-b
  # repos: []  ← omit or leave empty
```

No `GITHUB_TOKEN`, proxy, or CA certificate is needed. The server operates
entirely offline.

### 7.4 GitHub Enterprise Server

Point the GitHub client at your GHE instance by setting a custom base URL.
The GitHub client defaults to `https://api.github.com` but accepts any base URL.

```bash
# GHE typically uses https://github.example.com/api/v3
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
mcp-instructions --repos my-org/my-repo
```

Ensure firewall rules allow HTTPS to your GHE hostname on port 443.

### 7.5 Kubernetes

Deploy as a container with configuration via environment variables and
Kubernetes secrets.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mcp-instructions
spec:
  template:
    spec:
      containers:
        - name: mcp-instructions
          image: mcp-instructions:latest
          args: ["--transport", "http"]
          ports:
            - containerPort: 8080
          env:
            - name: GITHUB_TOKEN
              valueFrom:
                secretKeyRef:
                  name: mcp-secrets
                  key: github-token
            - name: INSTRUCTIONS_REPOS
              value: "my-org/copilot-instructions"
            - name: INSTRUCTIONS_PROXY_URL
              value: "http://egress-proxy.internal:3128"
            - name: INSTRUCTIONS_CA_CERT
              value: "/etc/ssl/certs/corp-ca.pem"
          volumeMounts:
            - name: ca-certs
              mountPath: /etc/ssl/certs/corp-ca.pem
              subPath: corp-ca.pem
              readOnly: true
      volumes:
        - name: ca-certs
          configMap:
            name: corp-ca-bundle
```

**Network Policy** (restrict egress to required destinations only):

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: mcp-instructions-egress
spec:
  podSelector:
    matchLabels:
      app: mcp-instructions
  policyTypes:
    - Egress
  egress:
    - to:
        - ipBlock:
            cidr: 0.0.0.0/0    # Adjust to proxy or GitHub IP ranges
      ports:
        - port: 443
          protocol: TCP
        - port: 3128
          protocol: TCP
```

---

## 8. Troubleshooting

### Common Errors

| Error | Cause | Resolution |
|-------|-------|------------|
| `x509: certificate signed by unknown authority` | Outbound HTTPS connection encounters a certificate not in the system trust store (typically a corporate TLS-inspecting proxy). | Set `ca_cert` to a PEM file containing your corporate proxy's CA certificate. See [§4 TLS / Certificate Configuration](#4-tls--certificate-configuration). |
| `proxyconnect tcp: dial tcp: lookup proxy.corp: no such host` | The configured proxy hostname cannot be resolved. | Verify `proxy_url` is correct and the proxy hostname resolves from the server's network. Check DNS configuration. |
| `HTTP 401 for .github/copilot-instructions.md` | GitHub API returned Unauthorized. | A `GITHUB_TOKEN` is **optional** for public repositories. For private repos, set it to a valid personal access token with `repo` (or `contents:read`) scope. Can also be set via `--github-token` flag, `INSTRUCTIONS_GITHUB_TOKEN` env var, or `github_token` in YAML config. |
| `HTTP 403 for .github/copilot-instructions.md` | GitHub API rate limit exceeded or insufficient permissions. | Unauthenticated requests are limited to 60/hour. Set `GITHUB_TOKEN` to raise the limit to 5000/hour. Check rate limits with `curl -H "Authorization: Bearer $GITHUB_TOKEN" https://api.github.com/rate_limit`. |
| `context deadline exceeded` | Outbound connection timed out (default 30s for GitHub, 60s for LLM). | Check that the firewall allows outbound HTTPS. If using a proxy, verify the proxy is reachable and forwarding correctly. |
| `LLM returned HTTP 401` | LLM API key is invalid or expired. | Verify `LLM_API_KEY` is set and valid. |

### Debug Logging

The servers log to stderr via Go's standard `log` package. Key log messages:

- `starting MCP server on stdio` / `starting MCP server on :8080` — confirms transport mode
- `LLM optimization enabled` — confirms the optimizer is active
- `LLM optimization failed, falling back: <error>` — LLM call failed; server falls back to concatenation
- `creating HTTP client: <error>` — fatal error in proxy/TLS setup (check `ca_cert` path and PEM validity)

Increase visibility by capturing stderr:

```bash
mcp-instructions --config config.yaml 2>mcp-debug.log
```

### Testing Connectivity

Verify outbound connectivity using `curl` commands that mirror what the servers
do internally.

**GitHub API:**

```bash
# Test GitHub.com
curl -sf -o /dev/null -w "%{http_code}" \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.raw+json" \
  "https://api.github.com/repos/OWNER/REPO/contents/.github/copilot-instructions.md"

# Through a proxy
curl -sf -o /dev/null -w "%{http_code}" \
  --proxy http://proxy.corp.example:8080 \
  --cacert /etc/ssl/certs/corp-ca-bundle.pem \
  -H "Authorization: Bearer $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.raw+json" \
  "https://api.github.com/repos/OWNER/REPO/contents/.github/copilot-instructions.md"
```

**LLM Endpoint:**

```bash
curl -sf -o /dev/null -w "%{http_code}" \
  -H "Authorization: Bearer $LLM_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"ping"}]}' \
  "https://api.openai.com/v1/chat/completions"
```

**Proxy reachability:**

```bash
curl -sf -o /dev/null -w "%{http_code}" \
  --proxy http://proxy.corp.example:8080 \
  "https://api.github.com"
```

---

## 9. Full Configuration Reference

Complete YAML configuration with all fields and their defaults. Fields marked
`# env:` show the corresponding environment variable.

```yaml
# ──────────────────────────────────────────────
# Sources — where to find instruction/skill files
# ──────────────────────────────────────────────
sources:
  # Local filesystem directories.
  # env: INSTRUCTIONS_DIRS or SKILLS_DIRS (comma-separated)
  # CLI: --dirs
  dirs:
    - /path/to/local/repo

  # GitHub repositories in "owner/repo" or "owner/repo@ref" format.
  # env: INSTRUCTIONS_REPOS or SKILLS_REPOS (comma-separated)
  # CLI: --repos
  repos:
    - github/awesome-copilot
    - github/awesome-copilot@main

# ──────────────────────────────────────────────
# Cache — local cache for remote repository content
# ──────────────────────────────────────────────
cache:
  # Directory for cached files.
  # Default: ~/.cache/mcp-instructions or ~/.cache/mcp-skills
  # env: INSTRUCTIONS_CACHE_DIR or SKILLS_CACHE_DIR
  # CLI: --cache-dir
  dir: ~/.cache/mcp-instructions

  # How often to re-sync from GitHub.
  # Default: 5m
  # env: INSTRUCTIONS_SYNC_INTERVAL or SKILLS_SYNC_INTERVAL
  # CLI: --sync-interval
  sync_interval: 5m

# ──────────────────────────────────────────────
# Proxy — outbound HTTP proxy and TLS settings
# ──────────────────────────────────────────────
proxy:
  # HTTP/HTTPS proxy URL. Takes precedence over HTTP_PROXY/HTTPS_PROXY.
  # env: INSTRUCTIONS_PROXY_URL or SKILLS_PROXY_URL
  # CLI: --proxy-url
  # Default: "" (no proxy; falls back to HTTP_PROXY/HTTPS_PROXY/NO_PROXY)
  proxy_url: ""

  # Path to a PEM-encoded CA certificate bundle.
  # Appended to the system certificate pool.
  # env: INSTRUCTIONS_CA_CERT or SKILLS_CA_CERT
  # CLI: --ca-cert
  # Default: "" (system certs only)
  ca_cert: ""

  # Disable TLS certificate verification. TESTING ONLY.
  # env: INSTRUCTIONS_TLS_INSECURE_SKIP_VERIFY or SKILLS_TLS_INSECURE_SKIP_VERIFY
  # Default: false
  tls_insecure_skip_verify: false

  # HTTP headers to forward from incoming MCP requests to outbound calls.
  # Only effective in http transport mode.
  # env: INSTRUCTIONS_HEADER_PASSTHROUGH or SKILLS_HEADER_PASSTHROUGH (comma-separated)
  # Default: [] (no headers forwarded)
  header_passthrough:
    - X-Request-ID
    - X-Correlation-ID

# ──────────────────────────────────────────────
# LLM — OpenAI-compatible endpoint for content optimization
# ──────────────────────────────────────────────
llm:
  # Base URL of the OpenAI-compatible API.
  # "/chat/completions" is appended automatically if not present.
  # env: LLM_ENDPOINT
  # CLI: --llm-endpoint
  endpoint: "https://api.openai.com/v1"

  # Model name for chat completions.
  # env: LLM_MODEL
  # CLI: --llm-model
  model: "gpt-4o-mini"

  # Whether LLM optimization is on by default.
  # Can be overridden per-request via the "optimize" parameter.
  # env: LLM_ENABLED (true/1)
  # Default: false
  enabled: false

  # API key — set via environment variable only (not stored in YAML).
  # env: LLM_API_KEY
  # apikey: (not configurable in YAML for security)

# ──────────────────────────────────────────────
# Transport — how clients connect to this server
# ──────────────────────────────────────────────

# MCP transport mode: "stdio" or "http".
# Default: stdio
# env: INSTRUCTIONS_TRANSPORT or SKILLS_TRANSPORT
# CLI: --transport
transport: stdio

# HTTP listen address (host:port). Only used when transport is "http".
# Default: ":8080" (instructions) or ":8081" (skills)
# env: INSTRUCTIONS_ADDR or SKILLS_ADDR
# CLI: --addr
addr: ":8080"

# ──────────────────────────────────────────────
# GitHub Token — set via environment variable only
# ──────────────────────────────────────────────
# env: GITHUB_TOKEN
# Required when sources.repos is configured.
# Not configurable in YAML for security.
```

### Environment Variables Summary

| Variable | Scope | Description |
|----------|-------|-------------|
| `GITHUB_TOKEN` | Both servers | GitHub personal access token |
| `LLM_ENDPOINT` | Both servers | OpenAI-compatible API base URL |
| `LLM_MODEL` | Both servers | LLM model name |
| `LLM_API_KEY` | Both servers | LLM API bearer token |
| `LLM_ENABLED` | Both servers | Enable LLM optimization by default (`true`/`1`) |
| `INSTRUCTIONS_*` | mcp-instructions | Server-specific overrides (see table in [§3.2](#32-environment-variables-per-server-prefix)) |
| `SKILLS_*` | mcp-skills | Server-specific overrides (see table in [§3.2](#32-environment-variables-per-server-prefix)) |
| `HTTP_PROXY` | Both servers | Standard Go HTTP proxy (fallback) |
| `HTTPS_PROXY` | Both servers | Standard Go HTTPS proxy (fallback) |
| `NO_PROXY` | Both servers | Hosts to bypass proxy (fallback) |

### CLI Flags Summary

| Flag | Description |
|------|-------------|
| `--config` | Path to YAML config file |
| `--dirs` | Comma-separated local directories |
| `--repos` | Comma-separated GitHub repos (`owner/repo[@ref]`) |
| `--transport` | Transport mode: `stdio` or `http` |
| `--addr` | HTTP listen address |
| `--cache-dir` | Local cache directory |
| `--sync-interval` | Cache sync interval (e.g., `5m`, `1h`) |
| `--llm-endpoint` | LLM endpoint URL |
| `--llm-model` | LLM model name |
| `--proxy-url` | HTTP/HTTPS proxy URL |
| `--ca-cert` | Path to PEM CA certificate bundle |
