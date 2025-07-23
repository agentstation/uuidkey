package main

import (
	"strings"
	"testing"
)

func TestDecodeCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut string
	}{
		{
			name:    "decode valid key with hyphens",
			args:    []string{"decode", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"},
			wantErr: false,
			wantOut: "d1756360-5da0-40df-9926-a76abff5601d",
		},
		{
			name:    "decode valid key without hyphens",
			args:    []string{"decode", "38QARV01ET0G6Z2CJD9VA2ZZAR0X"},
			wantErr: false,
			wantOut: "d1756360-5da0-40df-9926-a76abff5601d",
		},
		{
			name:    "decode lowercase key (should fail)",
			args:    []string{"decode", strings.ToLower("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")},
			wantErr: true,
			wantOut: "",
		},
		{
			name:    "missing key argument",
			args:    []string{"decode"},
			wantErr: true,
		},
		{
			name:    "invalid key format",
			args:    []string{"decode", "not-a-valid-key"},
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    []string{"decode", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X", "extra"},
			wantErr: true,
		},
		{
			name:    "json output",
			args:    []string{"decode", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X", "--json"},
			wantErr: false,
			wantOut: `"uuid":"d1756360-5da0-40df-9926-a76abff5601d"`,
		},
		{
			name:    "quiet mode",
			args:    []string{"decode", "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X", "-q"},
			wantErr: false,
			wantOut: "d1756360-5da0-40df-9926-a76abff5601d",
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
			
			// Additional checks for specific modes
			if strings.Contains(tt.name, "quiet mode") {
				trimmed := strings.TrimSpace(output)
				if trimmed != tt.wantOut {
					t.Errorf("quiet mode: expected exact output %q, got %q", tt.wantOut, trimmed)
				}
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	// Test that encode -> decode produces the original UUID
	originalUUID := "d1756360-5da0-40df-9926-a76abff5601d"
	
	// Encode
	encodeOut, err := executeCommand(rootCmd, "encode", originalUUID, "-q")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	encodedKey := strings.TrimSpace(encodeOut)
	
	// Decode
	decodeOut, err := executeCommand(rootCmd, "decode", encodedKey, "-q")
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	decodedUUID := strings.TrimSpace(decodeOut)
	
	// Compare
	if decodedUUID != originalUUID {
		t.Errorf("round trip failed: started with %q, ended with %q", originalUUID, decodedUUID)
	}
}

// TestDecodeErrors tests error cases for decode command
func TestDecodeErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "json output with error",
			args:    []string{"decode", "INVALID-KEY", "--json"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "error")
			},
		},
		{
			name:    "decode error case normal",
			args:    []string{"decode", "38QARV0-1ET0G6Z-2CJD9VA-INVALID"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "error parsing key")
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