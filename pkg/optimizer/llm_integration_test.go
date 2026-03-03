//go:build integration

package optimizer_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Arkestone/mcp/pkg/optimizer"
)

// TestLLMEndpointIntegration tests the optimizer against a real OpenAI-compatible
// endpoint.  Set the following environment variables to enable:
//
//	LLM_ENDPOINT  – base URL, e.g. https://api.openai.com/v1
//	LLM_API_KEY   – bearer token
//	LLM_MODEL     – model ID, e.g. gpt-4o-mini
//
// Example (Ollama):
//
//	LLM_ENDPOINT=http://localhost:11434/v1 \
//	LLM_API_KEY=ollama \
//	LLM_MODEL=llama3.2 \
//	go test -tags integration -v -run TestLLMEndpointIntegration ./pkg/optimizer/
func TestLLMEndpointIntegration(t *testing.T) {
	endpoint := os.Getenv("LLM_ENDPOINT")
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")

	if endpoint == "" || apiKey == "" {
		t.Skip("LLM_ENDPOINT and LLM_API_KEY must be set to run this test")
	}
	if model == "" {
		model = "llama-3.1-8b-instruct"
	}

	opt := optimizer.New(optimizer.LLMConfig{
		Endpoint: endpoint,
		APIKey:   apiKey,
		Model:    model,
		Enabled:  true,
	})
	if opt == nil {
		t.Fatal("optimizer returned nil — check endpoint and api key")
	}

	inputs := []optimizer.ContentInput{
		{
			Source:  "team-standards",
			Path:    ".github/copilot-instructions.md",
			Content: "Use Go 1.24. Always run gofmt. Prefer table-driven tests.",
		},
		{
			Source:  "project-conventions",
			Path:    ".github/copilot-instructions.md",
			Content: "Always run gofmt before committing. Write tests for every function. Use context.Context in all handlers.",
		},
	}

	result, err := opt.Optimize(context.Background(), inputs)
	if err != nil {
		t.Fatalf("Optimize failed: %v", err)
	}
	if len(strings.TrimSpace(result)) == 0 {
		t.Fatal("result is empty")
	}
	t.Logf("Consolidated output (%d chars):\n%s", len(result), result)
}
