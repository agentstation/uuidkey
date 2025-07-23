package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func init() {
	// Set test environment variable to prevent os.Exit during tests
	_ = os.Setenv("GO_TEST", "1")
}

// executeCommand runs a cobra command with the given arguments and returns the output
func executeCommand(_ *cobra.Command, args ...string) (output string, err error) {
	// Create a fresh command instance to avoid state pollution
	cmd := newRootCommand()
	
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	
	// Set this as the current command for output function
	setCurrentCommand(cmd)
	defer func() { setCurrentCommand(nil) }()
	
	// Reset flags to ensure clean state
	resetGlobalFlags()
	
	err = cmd.Execute()
	return buf.String(), err
}


// resetGlobalFlags resets all global flags to their default values
func resetGlobalFlags() {
	jsonOutput = false
	quietMode = false
	// Reset version command flags
	verifyChecksum = false
	verifyOnline = false
	// Reset apikey command flags
	prefix = ""
	entropy = 128
	// Reset uuid command flags
	uuidVersion = 4  // Default is v4
	timeFlag = ""
	// Reset key command flags
	noHyphens = false
}

// captureStdout captures stdout during function execution
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	return string(out)
}


// assertContains checks if the string contains the substring
func assertContains(t *testing.T, output, substr string) {
	t.Helper()
	if !strings.Contains(output, substr) {
		t.Errorf("expected output to contain %q, got %q", substr, output)
	}
}

// assertNotContains checks if the string does not contain the substring
func assertNotContains(t *testing.T, output, substr string) {
	t.Helper()
	if strings.Contains(output, substr) {
		t.Errorf("expected output to not contain %q, got %q", substr, output)
	}
}

// assertJSONContains checks if JSON output contains expected key
func assertJSONContains(t *testing.T, output, key string) {
	t.Helper()
	if !strings.Contains(output, `"`+key+`"`) {
		t.Errorf("expected JSON output to contain key %q, got %q", key, output)
	}
}

// isValidUUID checks if a string is a valid UUID format
func isValidUUID(s string) bool {
	s = strings.TrimSpace(s)
	// Basic UUID format check (8-4-4-4-12)
	parts := strings.Split(s, "-")
	if len(parts) != 5 {
		return false
	}
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		return false
	}
	// Check if all characters are hex
	for _, part := range parts {
		for _, c := range part {
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
				return false
			}
		}
	}
	return true
}

// isValidKey checks if a string is a valid Base32-Crockford key
func isValidKey(s string) bool {
	s = strings.TrimSpace(s)
	// Remove hyphens for validation
	s = strings.ReplaceAll(s, "-", "")
	// Valid key length is 28 characters without hyphens
	if len(s) != 28 {
		return false
	}
	// Check if all characters are valid crockford base32
	for _, c := range s {
		if !isValidCrockfordChar(byte(c)) {
			return false
		}
	}
	return true
}

// isValidCrockfordChar checks if a character is valid in crockford base32
func isValidCrockfordChar(c byte) bool {
	// Valid characters per Crockford spec:
	// Digits: 0-9
	// Letters: A-H, J-K, M-N, P-Q, R-T, V-X, Y-Z
	// Excluded: I, L (confusing with 1), O (confusing with 0), U (obscenity)
	return (c >= '0' && c <= '9') ||
		(c >= 'A' && c <= 'H') ||
		(c >= 'J' && c <= 'K') ||
		(c >= 'M' && c <= 'N') ||
		(c >= 'P' && c <= 'Q') ||
		(c >= 'R' && c <= 'T') ||
		(c >= 'V' && c <= 'Z' && c != 'U')
}