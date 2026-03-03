// Package optimizer uses an OpenAI-compatible LLM endpoint to merge, deduplicate,
// and consolidate multiple instruction files into a single coherent output.
package optimizer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// LLMConfig configures the OpenAI-compatible LLM used to consolidate content.
type LLMConfig struct {
	// Endpoint is the base URL (e.g. "https://api.openai.com/v1").
	Endpoint string `yaml:"endpoint"`
	// Model is the model name (e.g. "gpt-4o-mini").
	Model string `yaml:"model"`
	// APIKey is the bearer token.
	APIKey string `yaml:"-"`
	// Enabled controls whether LLM optimization is on by default.
	Enabled bool `yaml:"enabled"`
}

// Optimizer calls an OpenAI-compatible chat completions endpoint to consolidate content.
type Optimizer struct {
	cfg    LLMConfig
	client *http.Client
}

// New creates an Optimizer. Returns nil if LLM is not configured.
// An optional *http.Client can be passed to control proxy/TLS behavior;
// if nil, a default client with a 60 s timeout is used.
func New(cfg LLMConfig, client ...*http.Client) *Optimizer {
	if cfg.Endpoint == "" || cfg.APIKey == "" {
		return nil
	}
	var c *http.Client
	if len(client) > 0 && client[0] != nil {
		c = client[0]
	} else {
		c = &http.Client{Timeout: 60 * time.Second}
	}
	return &Optimizer{cfg: cfg, client: c}
}

// Enabled returns true if the optimizer is configured and available.
func (o *Optimizer) Enabled() bool {
	return o != nil
}

const systemPrompt = `You are an instruction consolidator. You receive multiple AI assistant instruction files from different sources. Your job is to merge them into a single, coherent set of instructions that:

1. Removes duplicate or redundant directives
2. Resolves conflicts by preferring the more specific instruction
3. Preserves all unique, actionable guidance
4. Organizes the output logically by topic (build commands, architecture, conventions, etc.)
5. Keeps the output in Markdown format
6. Does NOT add generic advice or filler — only preserve what was in the originals

Return ONLY the consolidated instructions in Markdown. No preamble or explanation.`

// Optimize takes multiple instruction contents and returns a single consolidated version.
func (o *Optimizer) Optimize(ctx context.Context, instructions []ContentInput) (string, error) {
	if !o.Enabled() {
		return ConcatRaw(instructions), nil
	}

	var userMsg strings.Builder
	for _, inst := range instructions {
		fmt.Fprintf(&userMsg, "## Source: %s (%s)\n\n%s\n\n---\n\n",
			inst.Source, inst.Path, inst.Content)
	}

	body := chatRequest{
		Model: o.cfg.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMsg.String()},
		},
		Temperature: floatPtr(0.2),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	endpoint := strings.TrimSuffix(o.cfg.Endpoint, "/")
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint += "/chat/completions"
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.cfg.APIKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling LLM endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decoding LLM response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("LLM returned no choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// ContentInput is the input to the optimizer.
type ContentInput struct {
	Source  string
	Path    string
	Content string
}

// ConcatRaw is the fallback when LLM is not available — simple concatenation.
func ConcatRaw(instructions []ContentInput) string {
	var sb strings.Builder
	for _, inst := range instructions {
		fmt.Fprintf(&sb, "# Instructions from %s (%s)\n\n%s\n\n---\n\n",
			inst.Source, inst.Path, inst.Content)
	}
	return sb.String()
}

// OpenAI-compatible request/response types

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature *float64      `json:"temperature,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

func floatPtr(f float64) *float64 { return &f }
