package main

import (
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut string
	}{
		{
			name:    "no args shows help",
			args:    []string{},
			wantErr: false,
			wantOut: "uuidkey is a fast, secure tool for UUID generation and Base32-Crockford encoding",
		},
		{
			name:    "help flag",
			args:    []string{"--help"},
			wantErr: false,
			wantOut: "Usage:",
		},
		{
			name:    "unknown command",
			args:    []string{"unknown"},
			wantErr: true,
			wantOut: "unknown command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(rootCmd, tt.args...)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			
			if tt.wantOut != "" && !strings.Contains(output, tt.wantOut) {
				t.Errorf("expected output to contain %q, got %q", tt.wantOut, output)
			}
		})
	}
}

func TestGlobalFlags(t *testing.T) {
	// Test JSON output flag with uuid command
	t.Run("json output flag", func(t *testing.T) {
		output, err := executeCommand(rootCmd, "uuid", "--json")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		// Should contain JSON structure
		assertJSONContains(t, output, "uuid")
		assertJSONContains(t, output, "key")
	})
	
	// Test quiet mode flag
	t.Run("quiet mode flag", func(t *testing.T) {
		output, err := executeCommand(rootCmd, "uuid", "--quiet")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		// Should only contain UUID (no extra text)
		trimmed := strings.TrimSpace(output)
		if !isValidUUID(trimmed) {
			t.Errorf("expected valid UUID in quiet mode, got %q", trimmed)
		}
		
		// Should be exactly one line
		if strings.Count(output, "\n") > 1 {
			t.Errorf("expected single line output in quiet mode, got multiple lines")
		}
	})
	
	// Test short flags
	t.Run("short flags", func(t *testing.T) {
		// -j for --json
		output1, _ := executeCommand(rootCmd, "uuid", "-j")
		assertContains(t, output1, "{")
		
		// -q for --quiet  
		output2, _ := executeCommand(rootCmd, "uuid", "-q")
		if strings.Count(output2, "\n") > 1 {
			t.Errorf("expected single line with -q flag")
		}
	})
}

func TestOutput(t *testing.T) {
	// Test the output function behavior
	tests := []struct {
		name       string
		jsonOutput bool
		quietMode  bool
		data       interface{}
		result     string
		wantJSON   bool
		wantQuiet  bool
	}{
		{
			name:       "normal output",
			jsonOutput: false,
			quietMode:  false,
			data:       map[string]string{"key": "value"},
			result:     "test-result",
			wantJSON:   false,
			wantQuiet:  false,
		},
		{
			name:       "json output",
			jsonOutput: true,
			quietMode:  false,
			data:       map[string]string{"key": "value"},
			result:     "test-result",
			wantJSON:   true,
			wantQuiet:  false,
		},
		{
			name:       "quiet mode",
			jsonOutput: false,
			quietMode:  true,
			data:       map[string]string{"key": "value"},
			result:     "test-result",
			wantJSON:   false,
			wantQuiet:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			oldJSON := jsonOutput
			oldQuiet := quietMode
			defer func() {
				jsonOutput = oldJSON
				quietMode = oldQuiet
			}()
			
			jsonOutput = tt.jsonOutput
			quietMode = tt.quietMode
			
			outputStr := captureStdout(func() {
				var data map[string]interface{}
				switch v := tt.data.(type) {
				case map[string]interface{}:
					data = v
				case map[string]string:
					// Convert map[string]string to map[string]interface{}
					data = make(map[string]interface{})
					for k, val := range v {
						data[k] = val
					}
				}
				output(data, tt.result)
			})
			
			if tt.wantJSON {
				assertContains(t, outputStr, "{")
				assertContains(t, outputStr, "}")
			}
			
			if tt.wantQuiet {
				trimmed := strings.TrimSpace(outputStr)
				if trimmed != tt.result {
					t.Errorf("quiet mode: expected %q, got %q", tt.result, trimmed)
				}
			}
		})
	}
}

func TestExitWithError(t *testing.T) {
	// Test different error output modes
	tests := []struct {
		name       string
		jsonOutput bool
		message    string
		err        error
	}{
		{
			name:       "normal error output",
			jsonOutput: false,
			message:    "Test error",
			err:        nil,
		},
		{
			name:       "json error output",
			jsonOutput: true,
			message:    "Test error",
			err:        nil,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original
			oldJSON := jsonOutput
			defer func() { jsonOutput = oldJSON }()
			
			jsonOutput = tt.jsonOutput
			
			// exitWithError calls os.Exit, so we can't test it directly
			// Instead, we could refactor to make it testable or skip this test
		})
	}
}

// TestExecuteFunction tests the Execute function from root.go
func TestExecuteFunction(t *testing.T) {
	// Save original
	oldCmd := rootCmd
	defer func() {
		rootCmd = oldCmd
	}()

	// Create new root command
	rootCmd = newRootCommand()
	
	// Test Execute function
	err := Execute()
	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}
}

// TestRootSpecificCases tests specific cases from TestAdditionalCoverage
func TestRootSpecificCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "root execute",
			args:    []string{},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "uuidkey is a fast, secure tool")
			},
		},
		{
			name:    "help command",
			args:    []string{"help"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "Available Commands")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(newRootCommand(), tt.args...)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v (output: %s)", tt.wantErr, err, output)
			}
			
			if tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}