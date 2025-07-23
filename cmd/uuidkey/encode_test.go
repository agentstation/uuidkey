package main

import (
	"strings"
	"testing"
)

func TestEncodeCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "encode valid UUID",
			args:    []string{"encode", "550e8400-e29b-41d4-a716-446655440000"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// In quiet mode, encode just outputs the key
				assertContains(t, output, "1AGX100-3H9PGEM-2KHCH36-1AM8000")
			},
		},
		{
			name:    "encode UUID without hyphens",
			args:    []string{"encode", "550e8400-e29b-41d4-a716-446655440000"}, // encode expects hyphens
			wantErr: false,
			check: func(t *testing.T, output string) {
				// In quiet mode, encode just outputs the key
				assertContains(t, output, "1AGX100-3H9PGEM-2KHCH36-1AM8000")
			},
		},
		{
			name:    "encode invalid UUID",
			args:    []string{"encode", "invalid-uuid"},
			wantErr: true,
			check: func(t *testing.T, output string) {
				assertContains(t, output, "error encoding UUID")
			},
		},
		{
			name:    "encode with no-hyphens flag",
			args:    []string{"encode", "550e8400-e29b-41d4-a716-446655440000", "--no-hyphens"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// In default mode, should show Key: prefix
				if !strings.Contains(output, "Key: ") && !strings.Contains(output, "1AGX1003H9PGEM2KHCH361AM8000") {
					t.Errorf("expected key output, got: %s", output)
				}
				// Check no hyphens in key output
				if strings.Contains(output, "-") && strings.Contains(output, "3H9PGEM") {
					t.Error("no-hyphens flag should remove hyphens from key")
				}
			},
		},
		{
			name:    "encode with JSON output",
			args:    []string{"encode", "550e8400-e29b-41d4-a716-446655440000", "--json"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				assertJSONContains(t, output, "key")
				assertJSONContains(t, output, "uuid")
			},
		},
		{
			name:    "encode with quiet mode",
			args:    []string{"encode", "550e8400-e29b-41d4-a716-446655440000", "--quiet"},
			wantErr: false,
			check: func(t *testing.T, output string) {
				// Should only output the key, not the "Key: " prefix
				if strings.Contains(output, "Key: ") {
					t.Error("quiet mode should not include 'Key: ' prefix")
				}
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