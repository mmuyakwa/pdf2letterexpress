package cli

import (
	"testing"
)

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand("TestApp", "1.0.0", "Test Description")

	if cmd == nil {
		t.Fatal("NewRootCommand returned nil")
	}

	if cmd.Use != "TestApp <PDF-file>" {
		t.Errorf("Expected Use to be 'TestApp <PDF-file>', got '%s'", cmd.Use)
	}

	if cmd.Short != "Test Description" {
		t.Errorf("Expected Short to be 'Test Description', got '%s'", cmd.Short)
	}

	if cmd.Version != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got '%s'", cmd.Version)
	}
}

func TestSetupLogging(t *testing.T) {
	config := &Config{
		Verbose:  true,
		LogLevel: "debug",
	}

	setupLogging(config)
}