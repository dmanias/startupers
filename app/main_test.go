package main

import (
	"bytes"
	"encoding/json"
	"github.com/dmanias/startupers/app/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIAskQuestion(t *testing.T) {
	// Load configuration from config.json
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Create a request to pass to our handler.
	question := Question{Query: "What is the capital of France?"}
	body, err := json.Marshal(question)
	if err != nil {
		t.Fatalf("Could not marshal question: %v", err)
	}
	req, err := http.NewRequest("POST", "/ask", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use the API key and URL from the loaded configuration
		handleAsk(cfg.AI.APIKey, w, r)
	})

	// Serve the HTTP request to our handler.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}
