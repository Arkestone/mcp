// Package testutil provides shared test helpers across the mcp monorepo.
// Import this package only from test files (*_test.go).
package testutil

import "github.com/Arkestone/mcp/pkg/optimizer"

// LLMConfig returns a LLMConfig suitable for creating a non-nil Optimizer in tests.
func LLMConfig() optimizer.LLMConfig {
	return optimizer.LLMConfig{
		Endpoint: "http://localhost:9999",
		APIKey:   "test-key",
		Model:    "test-model",
	}
}
