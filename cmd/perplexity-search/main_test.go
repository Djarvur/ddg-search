// Package main provides the CLI entry point for perplexity-search.
package main

import (
	"testing"
)

func TestVersionConstant(t *testing.T) {
	t.Parallel()

	if version != "dev" {
		t.Errorf("Expected version 'dev', got %q", version)
	}
}

func TestDefaultConstants(t *testing.T) {
	t.Parallel()

	if defaultMaxResults != 5 {
		t.Errorf("Expected defaultMaxResults 5, got %d", defaultMaxResults)
	}

	if defaultModel != "sonar-medium-online" {
		t.Errorf("Expected defaultModel 'sonar-medium-online', got %q", defaultModel)
	}
}

func TestErrNoQuery(t *testing.T) {
	t.Parallel()

	if errNoQuery == nil {
		t.Error("Expected errNoQuery to be defined")
	}

	if errNoQuery.Error() != "no search query provided" {
		t.Errorf("Expected 'no search query provided', got %q", errNoQuery.Error())
	}
}

func TestErrNoAPIKey(t *testing.T) {
	t.Parallel()

	if errNoAPIKey == nil {
		t.Error("Expected errNoAPIKey to be defined")
	}

	expectedMsg := "PERPLEXITY_API_KEY environment variable not set. Please set it in your .env file or shell environment."
	if errNoAPIKey.Error() != expectedMsg {
		t.Errorf("Expected %q, got %q", expectedMsg, errNoAPIKey.Error())
	}
}

func TestRunSearch_NoQuery(t *testing.T) {
	t.Parallel()

	// This test verifies the error variable is correctly defined
	if errNoQuery == nil {
		t.Error("Expected errNoQuery to be defined")
	}

	if errNoQuery.Error() != "no search query provided" {
		t.Errorf("Expected 'no search query provided', got %q", errNoQuery.Error())
	}
}

func TestRunSearch_FunctionExists(t *testing.T) {
	t.Parallel()

	// Test that the function signature is correct
	// This verifies the function exists and has the right signature
	_ = runSearch
}
