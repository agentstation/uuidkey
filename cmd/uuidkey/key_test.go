package main

import (
	"strings"
	"testing"
)

func TestKeyCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "generate new key (no input)",
			args:    []string{"key"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "UUID:")
				assertContains(t, output, "Key:")
				
				// Extract key from output
				lines := strings.Split(output, "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "Key:") {
						key := strings.TrimSpace(strings.TrimPrefix(line, "Key:"))
						if !isValidKey(key) {
							t.Errorf("invalid key format: %q", key)
						}
					}
				}
			},
		},
		{
			name:    "encode UUID to key",
			args:    []string{"key", "d1756360-5da0-40df-9926-a76abff5601d"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "UUID: d1756360-5da0-40df-9926-a76abff5601d")
				assertContains(t, output, "Key:  38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
			},
		},
		{
			name:    "decode key to UUID",
			args:    []string{"key", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "Key:  38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
				assertContains(t, output, "UUID: d1756360-5da0-40df-9926-a76abff5601d")
			},
		},
		{
			name:    "decode key without hyphens",
			args:    []string{"key", "38QARV01ET0G6Z2CJD9VA2ZZAR0X"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "Key:  38QARV01ET0G6Z2CJD9VA2ZZAR0X")
				assertContains(t, output, "UUID: d1756360-5da0-40df-9926-a76abff5601d")
			},
		},
		{
			name:    "invalid input",
			args:    []string{"key", "not-a-uuid-or-key"},
			wantErr: true,
		},
		{
			name:    "json output",
			args:    []string{"key", "d1756360-5da0-40df-9926-a76abff5601d", "--json"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "uuid")
				assertJSONContains(t, output, "key")
				assertContains(t, output, "d1756360-5da0-40df-9926-a76abff5601d")
				assertContains(t, output, "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
			},
		},
		{
			name:    "quiet mode - generate",
			args:    []string{"key", "-q"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should only output the key
				trimmed := strings.TrimSpace(output)
				if !isValidKey(trimmed) {
					t.Errorf("expected valid key in quiet mode, got %q", trimmed)
				}
			},
		},
		{
			name:    "quiet mode - encode",
			args:    []string{"key", "d1756360-5da0-40df-9926-a76abff5601d", "-q"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				trimmed := strings.TrimSpace(output)
				if trimmed != "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X" {
					t.Errorf("expected specific key, got %q", trimmed)
				}
			},
		},
		{
			name:    "quiet mode - decode",
			args:    []string{"key", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X", "-q"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				trimmed := strings.TrimSpace(output)
				if trimmed != "d1756360-5da0-40df-9926-a76abff5601d" {
					t.Errorf("expected specific UUID, got %q", trimmed)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(rootCmd, tt.args...)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v (output: %s)", tt.wantErr, err, output)
			}
			
			if tt.check != nil && err == nil {
				tt.check(t, output)
			}
		})
	}
}

func TestKeySmartDetection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string // "uuid" or "key"
	}{
		{
			name:     "detect UUID with hyphens",
			input:    "d1756360-5da0-40df-9926-a76abff5601d",
			wantType: "uuid",
		},
		{
			name:     "detect key with hyphens",
			input:    "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
			wantType: "key",
		},
		{
			name:     "detect key without hyphens",
			input:    "38QARV01ET0G6Z2CJD9VA2ZZAR0X",
			wantType: "key",
		},
		{
			name:     "detect lowercase UUID",
			input:    strings.ToLower("d1756360-5da0-40df-9926-a76abff5601d"),
			wantType: "uuid",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(rootCmd, "key", tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			
			// Check the output contains both UUID and Key
			assertContains(t, output, "UUID:")
			assertContains(t, output, "Key:")
		})
	}
}

// TestKeyErrors tests error cases for key command
func TestKeyErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "generate with v5",
			args:    []string{"key", "--version", "5"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "unsupported UUID version: 5")
			},
		},
		{
			name:    "encode error with json",
			args:    []string{"key", "invalid-uuid", "--json"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "error")
			},
		},
		{
			name:    "decode error",
			args:    []string{"key", "INVALIDKEY"},
			wantErr: true,
		},
		{
			name:    "generate error with json",
			args:    []string{"key", "--version", "5", "--json"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "error")
			},
		},
		{
			name:    "decode invalid with json",
			args:    []string{"key", "INVALIDKEY", "--json"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(newRootCommand(), tt.args...)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			
			if tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

// TestKeyCommandEdgeCases tests edge cases for the key command
func TestKeyCommandEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "decode invalid key format",
			args:    []string{"key", "38QARV0-1ET0G6Z-2CJD9VA-TOOLONG123"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "Error:")
			},
		},
		{
			name:    "encode invalid UUID",
			args:    []string{"key", "not-a-valid-uuid-format"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "invalid input")
			},
		},
		{
			name:    "with no-hyphens flag",
			args:    []string{"key", "d1756360-5da0-40df-9926-a76abff5601d", "--no-hyphens"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "38QARV01ET0G6Z2CJD9VA2ZZAR0X") // No hyphens
				assertNotContains(t, output, "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X") // With hyphens
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeCommand(rootCmd, tt.args...)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
			
			if tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

// TestKeySpecificCases tests specific cases from TestAdditionalCoverage
func TestKeySpecificCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "key with timestamp",
			args:    []string{"key", "--version", "6"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should generate v6 UUID
				assertContains(t, output, "UUID:")
			},
		},
		{
			name:    "key decode quiet mode",
			args:    []string{"key", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X", "--quiet"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should only output UUID
				if strings.Contains(output, "UUID:") {
					t.Error("quiet mode should not include 'UUID:' prefix")
				}
			},
		},
		{
			name:    "key encode error normal",
			args:    []string{"key", "550e8400-e29b-INVALID-a716-446655440000"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "invalid input")
			},
		},
		{
			name:    "key parse error normal",
			args:    []string{"key", "38QARV0-1ET0G6Z-2CJD9VA-INVALID"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "error parsing key")
			},
		},
		{
			name:    "key decode error normal",
			args:    []string{"key", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0I"}, // 'I' is invalid in Crockford
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "error parsing key")
			},
		},
		{
			name:    "key encode quiet",
			args:    []string{"key", "550e8400-e29b-41d4-a716-446655440000", "--quiet"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should only output key
				if strings.Contains(output, "Key:") {
					t.Error("quiet mode should not include 'Key:' prefix")
				}
			},
		},
		{
			name:    "key generate error normal",
			args:    []string{"key", "--version", "10"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "unsupported UUID version")
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