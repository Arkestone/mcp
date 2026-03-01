// Package config provides unified configuration loading for all MCP servers.
// Configuration is layered: YAML file → environment variables → CLI flags.
// Each MCP server specifies its own EnvPrefix to namespace env vars.
package config

import (
	"flag"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/Arkestone/mcp/pkg/httputil"
	"github.com/Arkestone/mcp/pkg/optimizer"
)

// Config holds configuration shared by all MCP servers.
type Config struct {
	Sources     Sources              `yaml:"sources"`
	Cache       CacheConfig          `yaml:"cache"`
	Proxy       httputil.ProxyConfig `yaml:"proxy"`
	LLM         optimizer.LLMConfig  `yaml:"llm"`
	Transport   string               `yaml:"transport"`
	Addr        string               `yaml:"addr"`
	GitHubToken string               `yaml:"github_token,omitempty"`
}

// Sources defines where to find content (local dirs and GitHub repos).
type Sources struct {
	Dirs  []string `yaml:"dirs"`
	Repos []string `yaml:"repos"` // "owner/repo" or "owner/repo@ref"
}

// CacheConfig controls local caching of remote repository content.
type CacheConfig struct {
	Dir          string        `yaml:"dir"`
	SyncInterval time.Duration `yaml:"sync_interval"`
}

// RepoRef represents a parsed "owner/repo@ref" reference.
type RepoRef struct {
	Owner string
	Repo  string
	Ref   string // branch/tag/sha; empty = default branch
}

// ParseRepoRef parses "owner/repo" or "owner/repo@ref" into a RepoRef.
func ParseRepoRef(s string) RepoRef {
	var r RepoRef
	if at := strings.IndexByte(s, '@'); at >= 0 {
		r.Ref = s[at+1:]
		s = s[:at]
	}
	parts := strings.SplitN(s, "/", 2)
	if len(parts) >= 1 {
		r.Owner = parts[0]
	}
	if len(parts) >= 2 {
		r.Repo = parts[1]
	}
	return r
}

// ParsedRepos returns parsed RepoRef values from Sources.Repos.
func (c *Config) ParsedRepos() []RepoRef {
	refs := make([]RepoRef, 0, len(c.Sources.Repos))
	for _, s := range c.Sources.Repos {
		refs = append(refs, ParseRepoRef(s))
	}
	return refs
}

// Options configures how Load resolves env vars and defaults per MCP server.
type Options struct {
	// EnvPrefix namespaces environment variables (e.g. "INSTRUCTIONS" → INSTRUCTIONS_DIRS).
	EnvPrefix string
	// DefaultAddr is the default HTTP listen address.
	DefaultAddr string
	// DefaultCacheName is the cache subdirectory name (e.g. "mcp-instructions").
	DefaultCacheName string
}

// Load reads configuration layered: YAML file → env vars → CLI flags → defaults.
func Load(opts Options) *Config {
	var flagConfig, flagDirs, flagRepos, flagTransport, flagAddr string
	var flagCacheDir, flagSyncInterval, flagLLMEndpoint, flagLLMModel string
	var flagProxyURL, flagCACert, flagGitHubToken string

	flag.StringVar(&flagConfig, "config", "", "Path to YAML config file")
	flag.StringVar(&flagDirs, "dirs", "", "Comma-separated local directories")
	flag.StringVar(&flagRepos, "repos", "", "Comma-separated GitHub repos (owner/repo[@ref])")
	flag.StringVar(&flagTransport, "transport", "", "Transport: stdio (default) or http")
	flag.StringVar(&flagAddr, "addr", "", "HTTP listen address")
	flag.StringVar(&flagCacheDir, "cache-dir", "", "Local cache directory for remote repos")
	flag.StringVar(&flagSyncInterval, "sync-interval", "", "Sync interval (e.g. 5m)")
	flag.StringVar(&flagLLMEndpoint, "llm-endpoint", "", "OpenAI-compatible LLM endpoint URL")
	flag.StringVar(&flagLLMModel, "llm-model", "", "LLM model name")
	flag.StringVar(&flagProxyURL, "proxy-url", "", "HTTP/HTTPS proxy URL")
	flag.StringVar(&flagCACert, "ca-cert", "", "Path to PEM CA certificate bundle")
	flag.StringVar(&flagGitHubToken, "github-token", "", "GitHub personal access token (optional, for private repos)")
	flag.Parse()

	configPath := flagConfig
	if configPath == "" {
		configPath = os.Getenv(opts.EnvPrefix + "_CONFIG")
	}
	if configPath == "" {
		for _, p := range []string{"config.yaml", "config.yml"} {
			if _, err := os.Stat(p); err == nil {
				configPath = p
				break
			}
		}
	}

	cfg := LoadFromFile(configPath)
	ApplyEnv(cfg, opts.EnvPrefix)
	applyFlags(cfg, flagDirs, flagRepos, flagTransport, flagAddr,
		flagCacheDir, flagSyncInterval, flagLLMEndpoint, flagLLMModel,
		flagProxyURL, flagCACert, flagGitHubToken)
	ApplyDefaults(cfg, opts)
	return cfg
}

// LoadFromFile reads a YAML config file. Returns empty Config if path is empty or missing.
func LoadFromFile(path string) *Config {
	cfg := &Config{}
	if path != "" {
		if data, err := os.ReadFile(path); err == nil {
			_ = yaml.Unmarshal(data, cfg)
		}
	}
	return cfg
}

// ApplyEnv overlays environment variables with the given prefix onto cfg.
func ApplyEnv(cfg *Config, prefix string) {
	if v := os.Getenv(prefix + "_DIRS"); v != "" {
		cfg.Sources.Dirs = SplitCSV(v)
	}
	if v := os.Getenv(prefix + "_REPOS"); v != "" {
		cfg.Sources.Repos = SplitCSV(v)
	}
	if v := os.Getenv(prefix + "_TRANSPORT"); v != "" {
		cfg.Transport = v
	}
	if v := os.Getenv(prefix + "_ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv(prefix + "_CACHE_DIR"); v != "" {
		cfg.Cache.Dir = v
	}
	if v := os.Getenv(prefix + "_SYNC_INTERVAL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Cache.SyncInterval = d
		}
	}
	// GitHub token: prefer prefixed env var, fall back to GITHUB_TOKEN, keep
	// any value already set via YAML config if neither env var is present.
	if v := os.Getenv(prefix + "_GITHUB_TOKEN"); v != "" {
		cfg.GitHubToken = v
	} else if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		cfg.GitHubToken = v
	}

	if v := os.Getenv("LLM_ENDPOINT"); v != "" {
		cfg.LLM.Endpoint = v
	}
	if v := os.Getenv("LLM_MODEL"); v != "" {
		cfg.LLM.Model = v
	}
	cfg.LLM.APIKey = os.Getenv("LLM_API_KEY")
	if v := os.Getenv("LLM_ENABLED"); v == "true" || v == "1" {
		cfg.LLM.Enabled = true
	}

	// Proxy settings — both prefixed and standard env vars.
	if v := os.Getenv(prefix + "_PROXY_URL"); v != "" {
		cfg.Proxy.ProxyURL = v
	}
	if v := os.Getenv(prefix + "_CA_CERT"); v != "" {
		cfg.Proxy.CACertFile = v
	}
	if v := os.Getenv(prefix + "_TLS_INSECURE_SKIP_VERIFY"); v == "true" || v == "1" {
		cfg.Proxy.TLSInsecureSkipVerify = true
	}
	if v := os.Getenv(prefix + "_HEADER_PASSTHROUGH"); v != "" {
		cfg.Proxy.HeaderPassthrough = SplitCSV(v)
	}
}

// ApplyDefaults fills zero-value fields with sensible defaults.
func ApplyDefaults(cfg *Config, opts Options) {
	if cfg.Transport == "" {
		cfg.Transport = "stdio"
	}
	if cfg.Addr == "" {
		cfg.Addr = opts.DefaultAddr
	}
	if cfg.Cache.Dir == "" {
		cfg.Cache.Dir = defaultCacheDir(opts.DefaultCacheName)
	}
	if cfg.Cache.SyncInterval == 0 {
		cfg.Cache.SyncInterval = 5 * time.Minute
	}
}

func applyFlags(cfg *Config, dirs, repos, transport, addr, cacheDir, syncInterval, llmEndpoint, llmModel, proxyURL, caCert, githubToken string) {
	if dirs != "" {
		cfg.Sources.Dirs = SplitCSV(dirs)
	}
	if repos != "" {
		cfg.Sources.Repos = SplitCSV(repos)
	}
	if transport != "" {
		cfg.Transport = transport
	}
	if addr != "" {
		cfg.Addr = addr
	}
	if cacheDir != "" {
		cfg.Cache.Dir = cacheDir
	}
	if syncInterval != "" {
		if d, err := time.ParseDuration(syncInterval); err == nil {
			cfg.Cache.SyncInterval = d
		}
	}
	if llmEndpoint != "" {
		cfg.LLM.Endpoint = llmEndpoint
	}
	if llmModel != "" {
		cfg.LLM.Model = llmModel
	}
	if proxyURL != "" {
		cfg.Proxy.ProxyURL = proxyURL
	}
	if caCert != "" {
		cfg.Proxy.CACertFile = caCert
	}
	if githubToken != "" {
		cfg.GitHubToken = githubToken
	}
}

// SplitCSV splits a comma-separated string into trimmed, non-empty parts.
func SplitCSV(s string) []string {
	var out []string
	for _, v := range strings.Split(s, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func defaultCacheDir(name string) string {
	if home, err := os.UserHomeDir(); err == nil {
		return home + "/.cache/" + name
	}
	return "/tmp/" + name + "-cache"
}
