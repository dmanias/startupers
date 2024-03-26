package config

import (
	"testing"
	"time"
)

// TestLoadConfig tests loading the configuration from a file.
func TestLoadConfig(t *testing.T) {
	// Assuming config.json is in the same directory as your test file.
	// In a real-world scenario, consider using a path that makes sense for your project structure.
	const configFile = "config.json"

	// Load the configuration
	config, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Here you test the loaded configuration against the expected values.
	// This is a basic example. Expand this based on your actual configuration structure and needs.

	// Test if the APIHost under web matches the expected value.
	if config.Web.APIHost != "0.0.0.0:3000" {
		t.Errorf("Expected APIHost to be %s, got %s", "0.0.0.0:3000", config.Web.APIHost)
	}

	// Convert milliseconds to time.Duration for comparison
	if config.Web.ReadTimeout != 5000*time.Millisecond {
		t.Errorf("Expected ReadTimeout to be %d, got %d", 5000*time.Millisecond, config.Web.ReadTimeout)
	}

	// Add more checks as necessary for each field...
}
