package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMainFunc(t *testing.T) {
	
	// Save original args and exit function
	oldArgs := os.Args
	oldExit := osExit
	defer func() { 
		os.Args = oldArgs
		osExit = oldExit
	}()

	// Mock osExit to capture exit code
	var exitCode int
	exitCalled := false
	osExit = func(code int) {
		exitCode = code
		exitCalled = true
	}

	// Test that main handles errors correctly
	os.Args = []string{"uuidkey", "invalid-command"}
	
	// Capture all output to prevent interference with coverage
	// Use originalStdout from testutil_test.go to ensure we restore the real stdout
	oldCmd := rootCmd
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	
	// Run main
	main()
	
	// Restore to the real original stdout, not a local copy
	_ = w.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr
	rootCmd = oldCmd
	_, _ = io.ReadAll(r)

	// Should have called exit with code 1
	if !exitCalled {
		t.Error("expected os.Exit to be called")
	}
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}

// TestMainFunction tests main.go execution with various scenarios
func TestMainFunction(t *testing.T) {
	
	// Save original args and exit function
	oldArgs := os.Args
	oldExit := osExit
	oldStderr := os.Stderr
	defer func() {
		os.Args = oldArgs
		osExit = oldExit
		os.Stderr = oldStderr
	}()

	// Redirect stderr to discard output
	os.Stderr, _ = os.Open(os.DevNull)

	// Mock osExit to capture exit code
	var exitCode int
	osExit = func(code int) {
		exitCode = code
		panic(fmt.Sprintf("os.Exit(%d)", code))
	}

	tests := []struct {
		name         string
		args         []string
		expectExit   bool
		expectedCode int
	}{
		{
			name:         "invalid command",
			args:         []string{"uuidkey", "invalid-command"},
			expectExit:   true,
			expectedCode: 1,
		},
		{
			name:         "help command",
			args:         []string{"uuidkey", "help"},
			expectExit:   false,
			expectedCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			exitCode = 0

			if tt.expectExit {
				defer func() {
					if r := recover(); r != nil {
						if msg, ok := r.(string); ok && strings.Contains(msg, "os.Exit") {
							if exitCode != tt.expectedCode {
								t.Errorf("expected exit code %d, got %d", tt.expectedCode, exitCode)
							}
						} else {
							t.Errorf("unexpected panic: %v", r)
						}
					}
				}()
			}

			// Capture stdout and stderr to prevent test output pollution
			oldCmd := rootCmd
			r, w, _ := os.Pipe()
			os.Stdout = w
			
			// Run main
			main()
			
			// Restore to the real original stdout
			_ = w.Close()
			os.Stdout = originalStdout
			rootCmd = oldCmd
			_, _ = io.ReadAll(r)

			if tt.expectExit {
				t.Error("expected os.Exit to be called")
			}
		})
	}
}

// TestMainExecute tests main execution scenarios
func TestMainExecute(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "help command",
			args:    []string{"help"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			args:    []string{"invalid-command"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := executeCommand(newRootCommand(), tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestBuildVariables(t *testing.T) {
	// Test that build variables are set to default values in test environment
	if version != "dev" {
		t.Errorf("expected version to be 'dev', got %q", version)
	}
	
	if commit != "none" {
		t.Errorf("expected commit to be 'none', got %q", commit)
	}
	
	if date != "unknown" {
		t.Errorf("expected date to be 'unknown', got %q", date)
	}
	
	if builtBy != "unknown" {
		t.Errorf("expected builtBy to be 'unknown', got %q", builtBy)
	}
}

func TestGoVersionAndPlatform(t *testing.T) {
	// These should be populated at init time
	if goVersion == "" {
		t.Error("expected goVersion to be populated")
	}
	
	if platform == "" {
		t.Error("expected platform to be populated")
	}
}