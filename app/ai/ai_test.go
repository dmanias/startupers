//go:build integration
// +build integration

package ai

import (
	"github.com/dmanias/startupers/app/config" // Adjust the import path to where your config package is located
	"testing"
)

func TestAskAIIntegration(t *testing.T) {
	// Load the configuration
	cfg, err := config.LoadConfig("../config/config.json") // Adjust the path to your actual config.json file location
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Use a known prompt
	prompt := "What is the capital of France?"

	// Call AskAI with the loaded API key
	response, err := AskAI(cfg.AI.APIKey, prompt)
	if err != nil {
		t.Fatalf("AskAI returned an error: %v", err)
	}

	// Perform basic validation on the response
	// Note: Since OpenAI's responses can vary, and without knowing the exact model you're using or its configuration,
	// it's challenging to predict the exact response content.
	// You might want to check for non-empty responses or other indicators of a successful interaction.
	if response == "" {
		t.Errorf("Expected a non-empty response from AskAI")
	}

	// Additional validation can be performed based on expected characteristics of the response
}
