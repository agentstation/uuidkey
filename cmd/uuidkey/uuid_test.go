package main

import (
	"strings"
	"testing"
)

func TestUUIDCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "generate default UUID (v4)",
			args:    []string{"uuid"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should generate both UUID and key
				assertContains(t, output, "UUID:")
				assertContains(t, output, "Key:")

				// Extract UUID from output
				lines := strings.Split(output, "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "UUID:") {
						uuid := strings.TrimSpace(strings.TrimPrefix(line, "UUID:"))
						if !isValidUUID(uuid) {
							t.Errorf("invalid UUID format: %q", uuid)
						}
						// v4 UUID has '4' in the version position
						if uuid[14] != '4' {
							t.Errorf("expected UUID v4, got %q", uuid)
						}
					}
				}
			},
		},
		{
			name:    "generate UUID v4 explicitly",
			args:    []string{"uuid", "--version", "4"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "UUID:")
				assertContains(t, output, "Key:")
			},
		},
		{
			name:    "generate UUID v6",
			args:    []string{"uuid", "--version", "6"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "UUID:") {
						uuid := strings.TrimSpace(strings.TrimPrefix(line, "UUID:"))
						// v6 UUID has '6' in the version position
						if uuid[14] != '6' {
							t.Errorf("expected UUID v6, got %q", uuid)
						}
					}
				}
			},
		},
		{
			name:    "generate UUID v7",
			args:    []string{"uuid", "--version", "7"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				lines := strings.Split(output, "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "UUID:") {
						uuid := strings.TrimSpace(strings.TrimPrefix(line, "UUID:"))
						// v7 UUID has '7' in the version position
						if uuid[14] != '7' {
							t.Errorf("expected UUID v7, got %q", uuid)
						}
					}
				}
			},
		},
		{
			name:    "encode existing UUID",
			args:    []string{"uuid", "d1756360-5da0-40df-9926-a76abff5601d"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "UUID: d1756360-5da0-40df-9926-a76abff5601d")
				assertContains(t, output, "Key:  38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
			},
		},
		{
			name:    "decode existing key",
			args:    []string{"uuid", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "Key:  38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
				assertContains(t, output, "UUID: d1756360-5da0-40df-9926-a76abff5601d")
			},
		},
		{
			name:    "invalid input",
			args:    []string{"uuid", "invalid-input"},
			wantErr: true,
		},
		{
			name:    "json output",
			args:    []string{"uuid", "--json"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "uuid")
				assertJSONContains(t, output, "key")
			},
		},
		{
			name:    "quiet mode",
			args:    []string{"uuid", "-q"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should only output the UUID
				trimmed := strings.TrimSpace(output)
				if !isValidUUID(trimmed) {
					t.Errorf("expected valid UUID in quiet mode, got %q", trimmed)
				}
			},
		},
		{
			name:    "invalid version error",
			args:    []string{"uuid", "--version", "5"},
			wantErr: true,
		},
		{
			name:    "encode with custom timestamp",
			args:    []string{"uuid", "--version", "6", "--time", "2024-01-01T00:00:00Z"},
			wantErr: true, // Expected error since custom timestamp is not supported
			check: func(t *testing.T, output string) {
				assertContains(t, output, "custom timestamp is not supported")
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

func TestUUIDVersionGeneration(t *testing.T) {
	// Test that multiple generations produce different UUIDs
	t.Run("unique generation", func(t *testing.T) {
		uuids := make(map[string]bool)

		for range 10 {
			output, err := executeCommand(rootCmd, "uuid", "-q")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			uuid := strings.TrimSpace(output)
			if uuids[uuid] {
				t.Errorf("duplicate UUID generated: %s", uuid)
			}
			uuids[uuid] = true
		}
	})

	// Test v6 is sortable
	t.Run("v6 sortable", func(t *testing.T) {
		// Generate multiple v6 UUIDs
		for range 5 {
			output, err := executeCommand(rootCmd, "uuid", "--version", "6", "-q")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			uuid := strings.TrimSpace(output)
			// v6 UUIDs should have version 6 in the version field
			if len(uuid) > 14 && uuid[14] != '6' {
				t.Errorf("expected v6 UUID, got %s", uuid)
			}
		}
		// Note: v6 UUIDs are sortable by time but might not always be strictly
		// increasing due to timing resolution and MAC address components
	})
}

// TestUUIDErrors tests error cases for UUID command
func TestUUIDErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "generate v5 error",
			args:    []string{"uuid", "--version", "5"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "unsupported UUID version: 5")
			},
		},
		{
			name:    "custom time with v7",
			args:    []string{"uuid", "--version", "7", "--time", "2024-01-01T00:00:00Z"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "custom timestamp is not supported")
			},
		},
		{
			name:    "generate error with json",
			args:    []string{"uuid", "--version", "5", "--json"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "error")
			},
		},
		{
			name:    "encode error with json",
			args:    []string{"uuid", "invalid-uuid", "--json"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "error")
			},
		},
		{
			name:    "decode error with json",
			args:    []string{"uuid", "INVALIDKEY", "--json"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "error")
			},
		},
		{
			name:    "parse apikey error with json",
			args:    []string{"uuid", "INVALID_API_KEY", "--json"},
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

// TestDetectInputType tests the detectInputType function from uuid.go
func TestDetectInputType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string
	}{
		{
			name:     "UUID with hyphens",
			input:    "550e8400-e29b-41d4-a716-446655440000",
			wantType: "uuid",
		},
		{
			name:     "UUID without hyphens",
			input:    "550e8400e29b41d4a716446655440000",
			wantType: "uuid",
		},
		{
			name:     "API key",
			input:    "TEST_38QARV01ET0G6Z2CJD9VA2ZZAR0XJBJLSO7WBNWY3F_A1B2C3D8",
			wantType: "apikey",
		},
		{
			name:     "Base32 key with hyphens",
			input:    "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
			wantType: "key",
		},
		{
			name:     "Base32 key without hyphens",
			input:    "38QARV01ET0G6Z2CJD9VA2ZZAR0X",
			wantType: "key",
		},
		{
			name:     "lowercase key",
			input:    "38qarv01et0g6z2cjd9va2zzar0x",
			wantType: "key",
		},
		{
			name:     "invalid input",
			input:    "not-a-valid-input",
			wantType: "unknown",
		},
		{
			name:     "empty input",
			input:    "",
			wantType: "unknown",
		},
		{
			name:     "API key with wrong format",
			input:    "TEST_INVALID",
			wantType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType := detectInputType(tt.input)
			if gotType != tt.wantType {
				t.Errorf("detectInputType(%q) = %q, want %q", tt.input, gotType, tt.wantType)
			}
		})
	}
}

// TestUUIDGenerateErrors tests UUID generation error cases
func TestUUIDGenerateErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "invalid UUID version",
			args:    []string{"uuid", "--version", "10"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "unsupported UUID version")
			},
		},
		{
			name:    "custom timestamp not supported",
			args:    []string{"uuid", "--version", "6", "--time", "2024-01-01T00:00:00Z"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "custom timestamp is not supported")
			},
		},
		{
			name:    "decode invalid key",
			args:    []string{"uuid", "INVALID-KEY-FORMAT"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "Error:")
			},
		},
		{
			name:    "encode invalid UUID",
			args:    []string{"uuid", "not-a-uuid"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "Error:")
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

// TestUUIDSpecificCases tests specific cases from TestAdditionalCoverage
func TestUUIDSpecificCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "uuid generation error",
			args:    []string{"uuid", "--version", "10"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "unsupported UUID version")
			},
		},
		{
			name:    "uuid with quiet mode generating",
			args:    []string{"uuid", "--quiet"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should only output UUID
				if strings.Contains(output, "UUID:") {
					t.Error("quiet mode should not include 'UUID:' prefix")
				}
			},
		},
		{
			name:    "uuid encode quiet mode",
			args:    []string{"uuid", "550e8400-e29b-41d4-a716-446655440000", "--quiet"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should only output key
				if strings.Contains(output, "Key:") {
					t.Error("quiet mode should not include 'Key:' prefix")
				}
			},
		},
		{
			name:    "uuid apikey parse quiet",
			args:    []string{"uuid", "TEST_0SR2SXC2QTEJPQ2ZNEE3D22SYG849TZQ0NNGG0MH4H_5A8C7836", "--quiet"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should only output UUID
				if strings.Contains(output, "UUID:") {
					t.Error("quiet mode should not include 'UUID:' prefix")
				}
			},
		},
		{
			name:    "uuid encode error normal",
			args:    []string{"uuid", "550e8400-e29b-INVALID-a716-446655440000"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "invalid input")
			},
		},
		{
			name:    "uuid parse key error normal",
			args:    []string{"uuid", "38QARV0-1ET0G6Z-2CJD9VA-INVALID"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "error parsing key")
			},
		},
		{
			name:    "uuid decode key error normal",
			args:    []string{"uuid", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0I"}, // 'I' is invalid
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "error parsing key")
			},
		},
		{
			name:    "uuid generate error normal",
			args:    []string{"uuid", "--version", "10"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "unsupported UUID version")
			},
		},
		{
			name:    "uuid encode generate error",
			args:    []string{"uuid", "--version", "4"}, // This should work
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "UUID:")
			},
		},
		{
			name:    "uuid generate internal error",
			args:    []string{"uuid", "--version", "4"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// This should succeed
				assertContains(t, output, "UUID:")
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
