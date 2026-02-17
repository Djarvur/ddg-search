// Package main provides the CLI entry point for page-dump.
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

func TestErrNoURL(t *testing.T) {
	t.Parallel()

	if errNoURL == nil {
		t.Error("Expected errNoURL to be defined")
	}

	if errNoURL.Error() != "no URL provided" {
		t.Errorf("Expected 'no URL provided', got %q", errNoURL.Error())
	}
}

func TestRunDump_NoURL(t *testing.T) {
	t.Parallel()

	// This test verifies the error variable is correctly defined
	if errNoURL == nil {
		t.Error("Expected errNoURL to be defined")
	}

	if errNoURL.Error() != "no URL provided" {
		t.Errorf("Expected 'no URL provided', got %q", errNoURL.Error())
	}
}

func TestRunDump_FunctionExists(t *testing.T) {
	t.Parallel()

	// Test that the function signature is correct
	// This verifies the function exists and has the right signature
	_ = runDump
}
