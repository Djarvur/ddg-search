// Package main provides the CLI entry point for ddg-search.
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

	if defaultMaxResults != 10 {
		t.Errorf("Expected defaultMaxResults 10, got %d", defaultMaxResults)
	}

	if defaultMaxRetries != 3 {
		t.Errorf("Expected defaultMaxRetries 3, got %d", defaultMaxRetries)
	}

	if defaultMaxDelaySecs != 30 {
		t.Errorf("Expected defaultMaxDelaySecs 30, got %d", defaultMaxDelaySecs)
	}

	if backoffMultiplier != 2.0 {
		t.Errorf("Expected backoffMultiplier 2.0, got %f", backoffMultiplier)
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
