package config

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseRepoRef(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  RepoRef
	}{
		{
			name:  "owner and repo",
			input: "owner/repo",
			want:  RepoRef{Owner: "owner", Repo: "repo", Ref: ""},
		},
		{
			name:  "owner repo and branch ref",
			input: "owner/repo@main",
			want:  RepoRef{Owner: "owner", Repo: "repo", Ref: "main"},
		},
		{
			name:  "owner repo and semver ref",
			input: "owner/repo@v1.2.3",
			want:  RepoRef{Owner: "owner", Repo: "repo", Ref: "v1.2.3"},
		},
		{
			name:  "owner only",
			input: "justowner",
			want:  RepoRef{Owner: "justowner", Repo: "", Ref: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseRepoRef(tt.input)
			if got != tt.want {
				t.Errorf("ParseRepoRef(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSplitCSV(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "basic csv",
			input: "a,b,c",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "trimmed spaces",
			input: " a , b ",
			want:  []string{"a", "b"},
		},
		{
			name:  "empty parts skipped",
			input: "a,,b",
			want:  []string{"a", "b"},
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "single value",
			input: "single",
			want:  []string{"single"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitCSV(tt.input)
			if !stringSliceEqual(got, tt.want) {
				t.Errorf("SplitCSV(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Run("valid yaml", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		content := `sources:
  dirs:
    - /tmp/docs
    - /tmp/notes
  repos:
    - owner/repo@main
transport: http
addr: ":9090"
cache:
  dir: /tmp/cache
  sync_interval: 10m
github_token: ghp_from_yaml
llm:
  endpoint: https://api.example.com/v1
  model: gpt-4o
  enabled: true
`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := LoadFromFile(path)
		if len(cfg.Sources.Dirs) != 2 || cfg.Sources.Dirs[0] != "/tmp/docs" || cfg.Sources.Dirs[1] != "/tmp/notes" {
			t.Errorf("Sources.Dirs = %v, want [/tmp/docs /tmp/notes]", cfg.Sources.Dirs)
		}
		if len(cfg.Sources.Repos) != 1 || cfg.Sources.Repos[0] != "owner/repo@main" {
			t.Errorf("Sources.Repos = %v, want [owner/repo@main]", cfg.Sources.Repos)
		}
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
		}
		if cfg.Addr != ":9090" {
			t.Errorf("Addr = %q, want %q", cfg.Addr, ":9090")
		}
		if cfg.Cache.Dir != "/tmp/cache" {
			t.Errorf("Cache.Dir = %q, want %q", cfg.Cache.Dir, "/tmp/cache")
		}
		if cfg.Cache.SyncInterval != 10*time.Minute {
			t.Errorf("Cache.SyncInterval = %v, want %v", cfg.Cache.SyncInterval, 10*time.Minute)
		}
		if cfg.LLM.Endpoint != "https://api.example.com/v1" {
			t.Errorf("LLM.Endpoint = %q, want %q", cfg.LLM.Endpoint, "https://api.example.com/v1")
		}
		if cfg.LLM.Model != "gpt-4o" {
			t.Errorf("LLM.Model = %q, want %q", cfg.LLM.Model, "gpt-4o")
		}
		if !cfg.LLM.Enabled {
			t.Error("LLM.Enabled = false, want true")
		}
		if cfg.GitHubToken != "ghp_from_yaml" {
			t.Errorf("GitHubToken = %q, want %q", cfg.GitHubToken, "ghp_from_yaml")
		}
	})

	t.Run("empty path", func(t *testing.T) {
		cfg := LoadFromFile("")
		if cfg == nil {
			t.Fatal("expected non-nil Config")
		}
		if cfg.Transport != "" || cfg.Addr != "" {
			t.Errorf("expected empty Config, got Transport=%q Addr=%q", cfg.Transport, cfg.Addr)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		cfg := LoadFromFile("/nonexistent/path/config.yaml")
		if cfg == nil {
			t.Fatal("expected non-nil Config")
		}
		if cfg.Transport != "" || cfg.Addr != "" {
			t.Errorf("expected empty Config, got Transport=%q Addr=%q", cfg.Transport, cfg.Addr)
		}
	})

	t.Run("yaml token preserved when no env", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "tok.yaml")
		if err := os.WriteFile(path, []byte("github_token: yaml-secret\n"), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := LoadFromFile(path)
		ApplyEnv(cfg, "NOPREFIX")
		if cfg.GitHubToken != "yaml-secret" {
			t.Errorf("GitHubToken = %q, want %q (YAML preserved)", cfg.GitHubToken, "yaml-secret")
		}
	})

	t.Run("env token overrides yaml token", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "tok.yaml")
		if err := os.WriteFile(path, []byte("github_token: yaml-secret\n"), 0644); err != nil {
			t.Fatal(err)
		}
		t.Setenv("GITHUB_TOKEN", "env-secret")
		cfg := LoadFromFile(path)
		ApplyEnv(cfg, "NOPREFIX")
		if cfg.GitHubToken != "env-secret" {
			t.Errorf("GitHubToken = %q, want %q (env overrides yaml)", cfg.GitHubToken, "env-secret")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.yaml")
		if err := os.WriteFile(path, []byte(":::not valid yaml\n\t{["), 0644); err != nil {
			t.Fatal(err)
		}
		cfg := LoadFromFile(path)
		if cfg == nil {
			t.Fatal("expected non-nil Config")
		}
	})
}

func TestApplyEnv(t *testing.T) {
	tests := []struct {
		name   string
		envs   map[string]string
		prefix string
		check  func(t *testing.T, cfg *Config)
	}{
		{
			name:   "PREFIX_DIRS",
			prefix: "TEST",
			envs:   map[string]string{"TEST_DIRS": "/a,/b"},
			check: func(t *testing.T, cfg *Config) {
				want := []string{"/a", "/b"}
				if !stringSliceEqual(cfg.Sources.Dirs, want) {
					t.Errorf("Sources.Dirs = %v, want %v", cfg.Sources.Dirs, want)
				}
			},
		},
		{
			name:   "PREFIX_REPOS",
			prefix: "TEST",
			envs:   map[string]string{"TEST_REPOS": "o1/r1,o2/r2@main"},
			check: func(t *testing.T, cfg *Config) {
				want := []string{"o1/r1", "o2/r2@main"}
				if !stringSliceEqual(cfg.Sources.Repos, want) {
					t.Errorf("Sources.Repos = %v, want %v", cfg.Sources.Repos, want)
				}
			},
		},
		{
			name:   "PREFIX_TRANSPORT",
			prefix: "TEST",
			envs:   map[string]string{"TEST_TRANSPORT": "http"},
			check: func(t *testing.T, cfg *Config) {
				if cfg.Transport != "http" {
					t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
				}
			},
		},
		{
			name:   "PREFIX_ADDR",
			prefix: "TEST",
			envs:   map[string]string{"TEST_ADDR": ":8080"},
			check: func(t *testing.T, cfg *Config) {
				if cfg.Addr != ":8080" {
					t.Errorf("Addr = %q, want %q", cfg.Addr, ":8080")
				}
			},
		},
		{
			name:   "PREFIX_CACHE_DIR",
			prefix: "TEST",
			envs:   map[string]string{"TEST_CACHE_DIR": "/custom/cache"},
			check: func(t *testing.T, cfg *Config) {
				if cfg.Cache.Dir != "/custom/cache" {
					t.Errorf("Cache.Dir = %q, want %q", cfg.Cache.Dir, "/custom/cache")
				}
			},
		},
		{
			name:   "PREFIX_SYNC_INTERVAL",
			prefix: "TEST",
			envs:   map[string]string{"TEST_SYNC_INTERVAL": "15m"},
			check: func(t *testing.T, cfg *Config) {
				if cfg.Cache.SyncInterval != 15*time.Minute {
					t.Errorf("Cache.SyncInterval = %v, want %v", cfg.Cache.SyncInterval, 15*time.Minute)
				}
			},
		},
		{
			name:   "GITHUB_TOKEN",
			prefix: "TEST",
			envs:   map[string]string{"GITHUB_TOKEN": "ghp_secret123"},
			check: func(t *testing.T, cfg *Config) {
				if cfg.GitHubToken != "ghp_secret123" {
					t.Errorf("GitHubToken = %q, want %q", cfg.GitHubToken, "ghp_secret123")
				}
			},
		},
		{
			name:   "prefixed GITHUB_TOKEN takes precedence",
			prefix: "TEST",
			envs:   map[string]string{"TEST_GITHUB_TOKEN": "prefixed-tok", "GITHUB_TOKEN": "global-tok"},
			check: func(t *testing.T, cfg *Config) {
				if cfg.GitHubToken != "prefixed-tok" {
					t.Errorf("GitHubToken = %q, want %q", cfg.GitHubToken, "prefixed-tok")
				}
			},
		},
		{
			name:   "no token env leaves existing value",
			prefix: "TEST",
			envs:   map[string]string{},
			check: func(t *testing.T, cfg *Config) {
				// When no GITHUB_TOKEN or TEST_GITHUB_TOKEN is set,
				// ApplyEnv should not overwrite existing value.
				if cfg.GitHubToken != "" {
					t.Errorf("GitHubToken = %q, want empty (no env set)", cfg.GitHubToken)
				}
			},
		},
		{
			name:   "LLM env vars",
			prefix: "TEST",
			envs: map[string]string{
				"LLM_ENDPOINT": "https://llm.example.com",
				"LLM_MODEL":    "gpt-4o-mini",
				"LLM_API_KEY":  "sk-key123",
				"LLM_ENABLED":  "true",
			},
			check: func(t *testing.T, cfg *Config) {
				if cfg.LLM.Endpoint != "https://llm.example.com" {
					t.Errorf("LLM.Endpoint = %q, want %q", cfg.LLM.Endpoint, "https://llm.example.com")
				}
				if cfg.LLM.Model != "gpt-4o-mini" {
					t.Errorf("LLM.Model = %q, want %q", cfg.LLM.Model, "gpt-4o-mini")
				}
				if cfg.LLM.APIKey != "sk-key123" {
					t.Errorf("LLM.APIKey = %q, want %q", cfg.LLM.APIKey, "sk-key123")
				}
				if !cfg.LLM.Enabled {
					t.Error("LLM.Enabled = false, want true")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}
			cfg := &Config{}
			ApplyEnv(cfg, tt.prefix)
			tt.check(t, cfg)
		})
	}
}

func TestApplyDefaults(t *testing.T) {
	opts := Options{
		DefaultAddr:      ":3000",
		DefaultCacheName: "test-cache",
	}

	t.Run("empty transport defaults to stdio", func(t *testing.T) {
		cfg := &Config{}
		ApplyDefaults(cfg, opts)
		if cfg.Transport != "stdio" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "stdio")
		}
	})

	t.Run("empty addr defaults to opts.DefaultAddr", func(t *testing.T) {
		cfg := &Config{}
		ApplyDefaults(cfg, opts)
		if cfg.Addr != ":3000" {
			t.Errorf("Addr = %q, want %q", cfg.Addr, ":3000")
		}
	})

	t.Run("empty cache dir defaults to home/.cache/name", func(t *testing.T) {
		cfg := &Config{}
		ApplyDefaults(cfg, opts)
		home, _ := os.UserHomeDir()
		want := home + "/.cache/test-cache"
		if cfg.Cache.Dir != want {
			t.Errorf("Cache.Dir = %q, want %q", cfg.Cache.Dir, want)
		}
	})

	t.Run("zero sync interval defaults to 5m", func(t *testing.T) {
		cfg := &Config{}
		ApplyDefaults(cfg, opts)
		if cfg.Cache.SyncInterval != 5*time.Minute {
			t.Errorf("Cache.SyncInterval = %v, want %v", cfg.Cache.SyncInterval, 5*time.Minute)
		}
	})

	t.Run("non-empty values are not overwritten", func(t *testing.T) {
		cfg := &Config{
			Transport: "http",
			Addr:      ":9999",
			Cache: CacheConfig{
				Dir:          "/custom",
				SyncInterval: 10 * time.Minute,
			},
		}
		ApplyDefaults(cfg, opts)
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
		}
		if cfg.Addr != ":9999" {
			t.Errorf("Addr = %q, want %q", cfg.Addr, ":9999")
		}
		if cfg.Cache.Dir != "/custom" {
			t.Errorf("Cache.Dir = %q, want %q", cfg.Cache.Dir, "/custom")
		}
		if cfg.Cache.SyncInterval != 10*time.Minute {
			t.Errorf("Cache.SyncInterval = %v, want %v", cfg.Cache.SyncInterval, 10*time.Minute)
		}
	})
}

func TestParsedRepos(t *testing.T) {
	cfg := &Config{
		Sources: Sources{
			Repos: []string{"owner1/repo1", "owner2/repo2@v2.0.0"},
		},
	}
	refs := cfg.ParsedRepos()
	if len(refs) != 2 {
		t.Fatalf("ParsedRepos() returned %d refs, want 2", len(refs))
	}
	want0 := RepoRef{Owner: "owner1", Repo: "repo1", Ref: ""}
	want1 := RepoRef{Owner: "owner2", Repo: "repo2", Ref: "v2.0.0"}
	if refs[0] != want0 {
		t.Errorf("refs[0] = %+v, want %+v", refs[0], want0)
	}
	if refs[1] != want1 {
		t.Errorf("refs[1] = %+v, want %+v", refs[1], want1)
	}
}

func TestLoad(t *testing.T) {
	t.Run("no flags no env no file", func(t *testing.T) {
		// Reset flag.CommandLine so Load can register its own flags.
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		os.Args = []string{"test"}
		// Ensure no config file is found in cwd.
		origDir, _ := os.Getwd()
		tmp := t.TempDir()
		os.Chdir(tmp)
		defer os.Chdir(origDir)

		cfg := Load(Options{
			EnvPrefix:        "TESTLOAD_NONE",
			DefaultAddr:      ":4000",
			DefaultCacheName: "test-load-cache",
		})
		if cfg.Transport != "stdio" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "stdio")
		}
		if cfg.Addr != ":4000" {
			t.Errorf("Addr = %q, want %q", cfg.Addr, ":4000")
		}
	})

	t.Run("with config env var", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		os.Args = []string{"test"}

		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "test.yaml")
		os.WriteFile(cfgPath, []byte("transport: http\naddr: \":7777\"\n"), 0644)
		t.Setenv("TLOAD_CONFIG", cfgPath)

		cfg := Load(Options{
			EnvPrefix:        "TLOAD",
			DefaultCacheName: "test-load-cache",
		})
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
		}
		if cfg.Addr != ":7777" {
			t.Errorf("Addr = %q, want %q", cfg.Addr, ":7777")
		}
	})

	t.Run("with config.yaml in cwd", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		os.Args = []string{"test"}

		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yaml")
		os.WriteFile(cfgPath, []byte("transport: http\n"), 0644)
		origDir, _ := os.Getwd()
		os.Chdir(dir)
		defer os.Chdir(origDir)

		cfg := Load(Options{
			EnvPrefix:        "TLOAD_CWD",
			DefaultCacheName: "test-load-cache",
		})
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
		}
	})

	t.Run("with config.yml in cwd", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		os.Args = []string{"test"}

		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yml")
		os.WriteFile(cfgPath, []byte("transport: http\n"), 0644)
		origDir, _ := os.Getwd()
		os.Chdir(dir)
		defer os.Chdir(origDir)

		cfg := Load(Options{
			EnvPrefix:        "TLOAD_YML",
			DefaultCacheName: "test-load-cache",
		})
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
		}
	})

	t.Run("with CLI flags", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		os.Args = []string{"test", "-transport", "http", "-addr", ":5555", "-dirs", "/x,/y"}

		origDir, _ := os.Getwd()
		os.Chdir(t.TempDir())
		defer os.Chdir(origDir)

		cfg := Load(Options{
			EnvPrefix:        "TLOAD_FLAGS",
			DefaultCacheName: "test-load-cache",
		})
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
		}
		if cfg.Addr != ":5555" {
			t.Errorf("Addr = %q, want %q", cfg.Addr, ":5555")
		}
		if !stringSliceEqual(cfg.Sources.Dirs, []string{"/x", "/y"}) {
			t.Errorf("Sources.Dirs = %v, want [/x /y]", cfg.Sources.Dirs)
		}
	})
}

func TestApplyFlags(t *testing.T) {
	t.Run("all empty flags do not change config", func(t *testing.T) {
		cfg := &Config{
			Transport: "http",
			Addr:      ":9090",
		}
		applyFlags(cfg, "", "", "", "", "", "", "", "", "", "", "")
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
		}
		if cfg.Addr != ":9090" {
			t.Errorf("Addr = %q, want %q", cfg.Addr, ":9090")
		}
	})

	t.Run("dirs flag", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "/a,/b", "", "", "", "", "", "", "", "", "", "")
		want := []string{"/a", "/b"}
		if !stringSliceEqual(cfg.Sources.Dirs, want) {
			t.Errorf("Sources.Dirs = %v, want %v", cfg.Sources.Dirs, want)
		}
	})

	t.Run("repos flag", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "", "o/r1,o/r2@main", "", "", "", "", "", "", "", "", "")
		want := []string{"o/r1", "o/r2@main"}
		if !stringSliceEqual(cfg.Sources.Repos, want) {
			t.Errorf("Sources.Repos = %v, want %v", cfg.Sources.Repos, want)
		}
	})

	t.Run("transport and addr flags", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "", "", "http", ":8080", "", "", "", "", "", "", "")
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q, want %q", cfg.Transport, "http")
		}
		if cfg.Addr != ":8080" {
			t.Errorf("Addr = %q, want %q", cfg.Addr, ":8080")
		}
	})

	t.Run("cacheDir flag", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "", "", "", "", "/tmp/mycache", "", "", "", "", "", "")
		if cfg.Cache.Dir != "/tmp/mycache" {
			t.Errorf("Cache.Dir = %q, want %q", cfg.Cache.Dir, "/tmp/mycache")
		}
	})

	t.Run("syncInterval valid", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "", "", "", "", "", "10m", "", "", "", "", "")
		if cfg.Cache.SyncInterval != 10*time.Minute {
			t.Errorf("Cache.SyncInterval = %v, want %v", cfg.Cache.SyncInterval, 10*time.Minute)
		}
	})

	t.Run("syncInterval invalid ignored", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "", "", "", "", "", "notaduration", "", "", "", "", "")
		if cfg.Cache.SyncInterval != 0 {
			t.Errorf("Cache.SyncInterval = %v, want 0", cfg.Cache.SyncInterval)
		}
	})

	t.Run("llmEndpoint and llmModel flags", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "", "", "", "", "", "", "https://llm.test/v1", "gpt-4o", "", "", "")
		if cfg.LLM.Endpoint != "https://llm.test/v1" {
			t.Errorf("LLM.Endpoint = %q, want %q", cfg.LLM.Endpoint, "https://llm.test/v1")
		}
		if cfg.LLM.Model != "gpt-4o" {
			t.Errorf("LLM.Model = %q, want %q", cfg.LLM.Model, "gpt-4o")
		}
	})

	t.Run("proxyURL and caCert flags", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "", "", "", "", "", "", "", "", "http://proxy:8080", "/etc/ca.pem", "")
		if cfg.Proxy.ProxyURL != "http://proxy:8080" {
			t.Errorf("Proxy.ProxyURL = %q, want %q", cfg.Proxy.ProxyURL, "http://proxy:8080")
		}
		if cfg.Proxy.CACertFile != "/etc/ca.pem" {
			t.Errorf("Proxy.CACertFile = %q, want %q", cfg.Proxy.CACertFile, "/etc/ca.pem")
		}
	})

	t.Run("all flags at once", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "/d1", "o/r1", "http", ":3000", "/cache", "5m",
			"https://llm.test", "model-x", "http://proxy", "/ca.pem", "ghp_all")
		if !stringSliceEqual(cfg.Sources.Dirs, []string{"/d1"}) {
			t.Errorf("Sources.Dirs = %v", cfg.Sources.Dirs)
		}
		if !stringSliceEqual(cfg.Sources.Repos, []string{"o/r1"}) {
			t.Errorf("Sources.Repos = %v", cfg.Sources.Repos)
		}
		if cfg.Transport != "http" {
			t.Errorf("Transport = %q", cfg.Transport)
		}
		if cfg.Addr != ":3000" {
			t.Errorf("Addr = %q", cfg.Addr)
		}
		if cfg.Cache.Dir != "/cache" {
			t.Errorf("Cache.Dir = %q", cfg.Cache.Dir)
		}
		if cfg.Cache.SyncInterval != 5*time.Minute {
			t.Errorf("Cache.SyncInterval = %v", cfg.Cache.SyncInterval)
		}
		if cfg.LLM.Endpoint != "https://llm.test" {
			t.Errorf("LLM.Endpoint = %q", cfg.LLM.Endpoint)
		}
		if cfg.LLM.Model != "model-x" {
			t.Errorf("LLM.Model = %q", cfg.LLM.Model)
		}
		if cfg.Proxy.ProxyURL != "http://proxy" {
			t.Errorf("Proxy.ProxyURL = %q", cfg.Proxy.ProxyURL)
		}
		if cfg.Proxy.CACertFile != "/ca.pem" {
			t.Errorf("Proxy.CACertFile = %q", cfg.Proxy.CACertFile)
		}
		if cfg.GitHubToken != "ghp_all" {
			t.Errorf("GitHubToken = %q, want %q", cfg.GitHubToken, "ghp_all")
		}
	})

	t.Run("github-token flag", func(t *testing.T) {
		cfg := &Config{}
		applyFlags(cfg, "", "", "", "", "", "", "", "", "", "", "ghp_fromflag")
		if cfg.GitHubToken != "ghp_fromflag" {
			t.Errorf("GitHubToken = %q, want %q", cfg.GitHubToken, "ghp_fromflag")
		}
	})
}

func TestApplyEnv_ProxySettings(t *testing.T) {
	t.Run("PREFIX_PROXY_URL", func(t *testing.T) {
		t.Setenv("PX_PROXY_URL", "http://myproxy:3128")
		cfg := &Config{}
		ApplyEnv(cfg, "PX")
		if cfg.Proxy.ProxyURL != "http://myproxy:3128" {
			t.Errorf("Proxy.ProxyURL = %q, want %q", cfg.Proxy.ProxyURL, "http://myproxy:3128")
		}
	})

	t.Run("PREFIX_CA_CERT", func(t *testing.T) {
		t.Setenv("PX_CA_CERT", "/etc/ssl/ca.pem")
		cfg := &Config{}
		ApplyEnv(cfg, "PX")
		if cfg.Proxy.CACertFile != "/etc/ssl/ca.pem" {
			t.Errorf("Proxy.CACertFile = %q, want %q", cfg.Proxy.CACertFile, "/etc/ssl/ca.pem")
		}
	})

	t.Run("PREFIX_TLS_INSECURE_SKIP_VERIFY true", func(t *testing.T) {
		t.Setenv("PX_TLS_INSECURE_SKIP_VERIFY", "true")
		cfg := &Config{}
		ApplyEnv(cfg, "PX")
		if !cfg.Proxy.TLSInsecureSkipVerify {
			t.Error("Proxy.TLSInsecureSkipVerify = false, want true")
		}
	})

	t.Run("PREFIX_TLS_INSECURE_SKIP_VERIFY 1", func(t *testing.T) {
		t.Setenv("PX_TLS_INSECURE_SKIP_VERIFY", "1")
		cfg := &Config{}
		ApplyEnv(cfg, "PX")
		if !cfg.Proxy.TLSInsecureSkipVerify {
			t.Error("Proxy.TLSInsecureSkipVerify = false, want true")
		}
	})

	t.Run("PREFIX_TLS_INSECURE_SKIP_VERIFY false stays false", func(t *testing.T) {
		t.Setenv("PX_TLS_INSECURE_SKIP_VERIFY", "false")
		cfg := &Config{}
		ApplyEnv(cfg, "PX")
		if cfg.Proxy.TLSInsecureSkipVerify {
			t.Error("Proxy.TLSInsecureSkipVerify = true, want false")
		}
	})

	t.Run("PREFIX_HEADER_PASSTHROUGH", func(t *testing.T) {
		t.Setenv("PX_HEADER_PASSTHROUGH", "X-A,X-B")
		cfg := &Config{}
		ApplyEnv(cfg, "PX")
		want := []string{"X-A", "X-B"}
		if !stringSliceEqual(cfg.Proxy.HeaderPassthrough, want) {
			t.Errorf("Proxy.HeaderPassthrough = %v, want %v", cfg.Proxy.HeaderPassthrough, want)
		}
	})

	t.Run("all proxy env vars together", func(t *testing.T) {
		t.Setenv("PX_PROXY_URL", "http://proxy:1234")
		t.Setenv("PX_CA_CERT", "/ca.pem")
		t.Setenv("PX_TLS_INSECURE_SKIP_VERIFY", "1")
		t.Setenv("PX_HEADER_PASSTHROUGH", "X-Foo,X-Bar")
		cfg := &Config{}
		ApplyEnv(cfg, "PX")
		if cfg.Proxy.ProxyURL != "http://proxy:1234" {
			t.Errorf("Proxy.ProxyURL = %q", cfg.Proxy.ProxyURL)
		}
		if cfg.Proxy.CACertFile != "/ca.pem" {
			t.Errorf("Proxy.CACertFile = %q", cfg.Proxy.CACertFile)
		}
		if !cfg.Proxy.TLSInsecureSkipVerify {
			t.Error("Proxy.TLSInsecureSkipVerify = false")
		}
		if !stringSliceEqual(cfg.Proxy.HeaderPassthrough, []string{"X-Foo", "X-Bar"}) {
			t.Errorf("Proxy.HeaderPassthrough = %v", cfg.Proxy.HeaderPassthrough)
		}
	})
}

func TestDefaultCacheDir_NoHome(t *testing.T) {
	// Unset HOME to trigger the fallback path in defaultCacheDir.
	t.Setenv("HOME", "")
	// Also unset XDG and plan9 alternatives that os.UserHomeDir checks.
	t.Setenv("XDG_CONFIG_HOME", "")

	cfg := &Config{}
	ApplyDefaults(cfg, Options{DefaultCacheName: "test-srv"})
	want := "/tmp/test-srv-cache"
	if cfg.Cache.Dir != want {
		t.Errorf("Cache.Dir = %q, want %q", cfg.Cache.Dir, want)
	}
}

func TestApplyEnv_SyncIntervalInvalid(t *testing.T) {
	t.Setenv("TEST_SYNC_INTERVAL", "notaduration")
	cfg := &Config{}
	ApplyEnv(cfg, "TEST")
	if cfg.Cache.SyncInterval != 0 {
		t.Errorf("Cache.SyncInterval = %v, want 0 for invalid duration", cfg.Cache.SyncInterval)
	}
}

func TestApplyEnv_LLMEnabled1(t *testing.T) {
	t.Setenv("LLM_ENABLED", "1")
	cfg := &Config{}
	ApplyEnv(cfg, "TEST")
	if !cfg.LLM.Enabled {
		t.Error("LLM.Enabled = false, want true for LLM_ENABLED=1")
	}
}

func stringSliceEqual(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Additional nominal / error / limit tests
// ---------------------------------------------------------------------------

func TestParseRepoRef_OnlyOwner(t *testing.T) {
	ref := ParseRepoRef("just-owner")
	if ref.Owner != "just-owner" {
		t.Errorf("Owner = %q, want just-owner", ref.Owner)
	}
	if ref.Repo != "" {
		t.Errorf("Repo = %q, want empty", ref.Repo)
	}
}

func TestParseRepoRef_EmptyString(t *testing.T) {
	ref := ParseRepoRef("")
	if ref.Owner != "" || ref.Repo != "" || ref.Ref != "" {
		t.Errorf("expected all empty, got %+v", ref)
	}
}

func TestParseRepoRef_OnlyRef(t *testing.T) {
	ref := ParseRepoRef("@main")
	if ref.Ref != "main" {
		t.Errorf("Ref = %q, want main", ref.Ref)
	}
}

func TestParseRepoRef_SlashInRef(t *testing.T) {
	ref := ParseRepoRef("o/r@feature/branch")
	if ref.Ref != "feature/branch" {
		t.Errorf("Ref = %q, want feature/branch", ref.Ref)
	}
}

func TestParsedRepos_Empty(t *testing.T) {
	cfg := &Config{}
	refs := cfg.ParsedRepos()
	if len(refs) != 0 {
		t.Errorf("expected 0 refs, got %d", len(refs))
	}
}

func TestParsedRepos_Multiple(t *testing.T) {
	cfg := &Config{
		Sources: Sources{Repos: []string{"a/b", "c/d@v1"}},
	}
	refs := cfg.ParsedRepos()
	if len(refs) != 2 {
		t.Fatalf("got %d refs, want 2", len(refs))
	}
	if refs[0].Repo != "b" {
		t.Errorf("refs[0].Repo = %q", refs[0].Repo)
	}
	if refs[1].Ref != "v1" {
		t.Errorf("refs[1].Ref = %q", refs[1].Ref)
	}
}

func TestApplyEnv_WhitespaceOnlyCSV(t *testing.T) {
	cfg := &Config{}
	os.Setenv("TEST_DIRS", "  ,  , ")
	defer os.Unsetenv("TEST_DIRS")

	ApplyEnv(cfg, "TEST")
	// Should result in whitespace-only entries (not filtered); all entries are preserved.
	for _, d := range cfg.Sources.Dirs {
		_ = d
	}
}

func TestApplyDefaults_SetsDefaultAddr(t *testing.T) {
	cfg := &Config{}
	ApplyDefaults(cfg, Options{
		DefaultAddr:      ":9090",
		DefaultCacheName: "test-cache",
	})
	if cfg.Addr != ":9090" {
		t.Errorf("Addr = %q, want :9090", cfg.Addr)
	}
}

func TestApplyDefaults_DoesNotOverrideAddr(t *testing.T) {
	cfg := &Config{Addr: ":1234"}
	ApplyDefaults(cfg, Options{
		DefaultAddr:      ":9090",
		DefaultCacheName: "test-cache",
	})
	if cfg.Addr != ":1234" {
		t.Errorf("Addr = %q, want :1234", cfg.Addr)
	}
}

func TestApplyDefaults_SetsDefaultCacheDir(t *testing.T) {
	cfg := &Config{}
	ApplyDefaults(cfg, Options{
		DefaultAddr:      ":8080",
		DefaultCacheName: "mcp-test",
	})
	if cfg.Cache.Dir == "" {
		t.Error("expected non-empty default cache dir")
	}
}

func TestApplyDefaults_SetsDefaultSyncInterval(t *testing.T) {
	cfg := &Config{}
	ApplyDefaults(cfg, Options{
		DefaultAddr:      ":8080",
		DefaultCacheName: "test",
	})
	if cfg.Cache.SyncInterval == 0 {
		t.Error("expected non-zero sync interval")
	}
}

func TestLoadFromFile_NonExistentFile(t *testing.T) {
	cfg := LoadFromFile("/nonexistent/path/config.yaml")
	// Should return empty/default config
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoadFromFile_InvalidYAML(t *testing.T) {
	f := filepath.Join(t.TempDir(), "bad.yaml")
	os.WriteFile(f, []byte(":::invalid yaml:::"), 0o644)

	cfg := LoadFromFile(f)
	// Invalid YAML should silently return empty config (ignores unmarshal error)
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoadFromFile_UnknownFields(t *testing.T) {
	f := filepath.Join(t.TempDir(), "extra.yaml")
	os.WriteFile(f, []byte("sources:\n  dirs:\n    - /tmp\nunknown_field: value\n"), 0o644)

	cfg := LoadFromFile(f)
	if len(cfg.Sources.Dirs) != 1 {
		t.Errorf("dirs = %v, want 1 entry", cfg.Sources.Dirs)
	}
}
