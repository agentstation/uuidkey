package uuidkey

import (
	"strings"
	"testing"
)

func TestAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		uuid     string
		options  []Option
		wantLen  int // expected length of entropy
		wantErr  bool
		validate func(t *testing.T, key APIKey)
	}{
		{
			name:    "Default settings (160-bit entropy)",
			prefix:  "TEST",
			uuid:    "d1756360-5da0-40df-9926-a76abff5601d",
			wantLen: 21,
			validate: func(t *testing.T, key APIKey) {
				if key.Prefix != "TEST" {
					t.Errorf("wrong prefix, got %s, want TEST", key.Prefix)
				}
				if len(key.Entropy) != 21 {
					t.Errorf("wrong entropy length, got %d, want 21", len(key.Entropy))
				}
			},
		},
		{
			name:    "128-bit entropy",
			prefix:  "DEMO",
			uuid:    "d1756360-5da0-40df-9926-a76abff5601d",
			options: []Option{With128BitEntropy},
			wantLen: 14,
			validate: func(t *testing.T, key APIKey) {
				if len(key.Entropy) != 14 {
					t.Errorf("wrong entropy length, got %d, want 14", len(key.Entropy))
				}
			},
		},
		{
			name:    "256-bit entropy",
			prefix:  "PROD",
			uuid:    "d1756360-5da0-40df-9926-a76abff5601d",
			options: []Option{With256BitEntropy},
			wantLen: 42,
			validate: func(t *testing.T, key APIKey) {
				if len(key.Entropy) != 42 {
					t.Errorf("wrong entropy length, got %d, want 42", len(key.Entropy))
				}
			},
		},
		{
			name:    "160-bit entropy",
			prefix:  "TEST",
			uuid:    "d1756360-5da0-40df-9926-a76abff5601d",
			options: []Option{With160BitEntropy},
			wantLen: 21,
			validate: func(t *testing.T, key APIKey) {
				if len(key.Entropy) != 21 {
					t.Errorf("wrong entropy length, got %d, want 21", len(key.Entropy))
				}
			},
		},
		{
			name:    "Empty prefix",
			prefix:  "",
			uuid:    "d1756360-5da0-40df-9926-a76abff5601d",
			wantErr: true,
			validate: func(t *testing.T, key APIKey) {
				// We shouldn't reach this point if wantErr is true
				t.Error("expected error for empty prefix, got none")
			},
		},
		{
			name:    "Multiple options",
			prefix:  "TEST",
			uuid:    "d1756360-5da0-40df-9926-a76abff5601d",
			options: []Option{With128BitEntropy, With256BitEntropy}, // Last option should win
			wantLen: 42,
			validate: func(t *testing.T, key APIKey) {
				if len(key.Entropy) != 42 {
					t.Errorf("wrong entropy length, got %d, want 42", len(key.Entropy))
				}
			},
		},
		{
			name:    "Entropy padding check",
			prefix:  "TEST",
			uuid:    "d1756360-5da0-40df-9926-a76abff5601d",
			options: []Option{With128BitEntropy},
			wantLen: 14,
			validate: func(t *testing.T, key APIKey) {
				if len(key.Entropy) != 14 {
					t.Errorf("wrong entropy length, got %d, want 14", len(key.Entropy))
				}
				// Ensure any leading zeros are properly handled
				for _, c := range key.Entropy {
					if !strings.ContainsRune("0123456789ABCDEFGHJKMNPQRSTVWXYZ", c) {
						t.Errorf("invalid character in entropy: %c", c)
					}
				}
				// Parse and verify the key can be reconstructed
				parsed, err := ParseAPIKey(key.String())
				if err != nil {
					t.Errorf("failed to parse key with padded entropy: %v", err)
				}
				if parsed.Entropy != key.Entropy {
					t.Errorf("entropy mismatch after parsing, got %s, want %s", parsed.Entropy, key.Entropy)
				}
			},
		},
		{
			name:    "Entropy padding verification",
			prefix:  "TEST",
			uuid:    "00000000-0000-0000-0000-000000000000", // Using all zeros to maximize chance of padding
			options: []Option{With128BitEntropy},
			wantLen: 14,
			validate: func(t *testing.T, key APIKey) {
				// Check length
				if len(key.Entropy) != 14 {
					t.Errorf("wrong entropy length, got %d, want 14", len(key.Entropy))
				}

				// Verify any leading zeros are preserved
				parsed, err := ParseAPIKey(key.String())
				if err != nil {
					t.Errorf("failed to parse key: %v", err)
				}

				// Verify the entropy matches after a round trip
				if parsed.Entropy != key.Entropy {
					t.Errorf("entropy mismatch after parsing, got %s, want %s", parsed.Entropy, key.Entropy)
				}

				// Verify the length is maintained
				if len(parsed.Entropy) != 14 {
					t.Errorf("parsed entropy length wrong, got %d, want 14", len(parsed.Entropy))
				}

				// Verify all characters are valid Base32-Crockford
				for _, c := range key.Entropy {
					if !strings.ContainsRune("0123456789ABCDEFGHJKMNPQRSTVWXYZ", c) {
						t.Errorf("invalid character in entropy: %c", c)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Catch panics and convert to errors if we expect an error
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			// Test NewAPIKey
			key, err := NewAPIKey(tt.prefix, tt.uuid, tt.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Validate entropy length
			if len(key.Entropy) != tt.wantLen {
				t.Errorf("entropy length mismatch, got %d, want %d", len(key.Entropy), tt.wantLen)
			}

			// Run additional validations if provided
			if tt.validate != nil {
				tt.validate(t, key)
			}

			// Test String() method
			str := key.String()
			parts := strings.Split(str, "_")
			if len(parts) != 3 {
				t.Errorf("String() result should have 3 parts (prefix_key_checksum), got %d parts", len(parts))
			}
			if parts[0] != tt.prefix {
				t.Errorf("String() prefix mismatch, got %s, want %s", parts[0], tt.prefix)
			}
			if len(parts[2]) != 8 {
				t.Errorf("String() checksum length mismatch, got %d, want 8", len(parts[2]))
			}

			// Test ParseAPIKey
			parsed, err := ParseAPIKey(str)
			if err != nil {
				t.Fatalf("ParseAPIKey() failed: %v", err)
			}
			if parsed.String() != str {
				t.Errorf("ParseAPIKey() result mismatch, got %s, want %s", parsed.String(), str)
			}
			if parsed.Prefix != key.Prefix {
				t.Errorf("ParseAPIKey() prefix mismatch, got %s, want %s", parsed.Prefix, key.Prefix)
			}
			if parsed.Key != key.Key {
				t.Errorf("ParseAPIKey() key mismatch, got %s, want %s", parsed.Key, key.Key)
			}
			if parsed.Entropy != key.Entropy {
				t.Errorf("ParseAPIKey() entropy mismatch, got %s, want %s", parsed.Entropy, key.Entropy)
			}
		})
	}
}

func TestNewAPIKeyFromBytes(t *testing.T) {
	uuid := [16]byte{
		0xd1, 0x75, 0x63, 0x60,
		0x5d, 0xa0, 0x40, 0xdf,
		0x99, 0x26, 0xa7, 0x6a,
		0xbf, 0xf5, 0x60, 0x1d,
	}

	tests := []struct {
		name     string
		prefix   string
		options  []Option
		wantLen  int
		validate func(t *testing.T, key APIKey)
	}{
		{
			name:    "Default settings",
			prefix:  "TEST",
			wantLen: 21,
			validate: func(t *testing.T, key APIKey) {
				if key.Prefix != "TEST" {
					t.Errorf("wrong prefix, got %s, want TEST", key.Prefix)
				}
				if len(key.Entropy) != 21 {
					t.Errorf("wrong entropy length, got %d, want 21", len(key.Entropy))
				}
			},
		},
		{
			name:    "128-bit entropy",
			prefix:  "DEMO",
			options: []Option{With128BitEntropy},
			wantLen: 14,
			validate: func(t *testing.T, key APIKey) {
				if len(key.Entropy) != 14 {
					t.Errorf("wrong entropy length, got %d, want 14", len(key.Entropy))
				}
			},
		},
		{
			name:    "256-bit entropy",
			prefix:  "PROD",
			options: []Option{With256BitEntropy},
			wantLen: 42,
			validate: func(t *testing.T, key APIKey) {
				if len(key.Entropy) != 42 {
					t.Errorf("wrong entropy length, got %d, want 42", len(key.Entropy))
				}
				// Validate entropy characters are in valid Base32-Crockford range
				for _, c := range key.Entropy {
					if !strings.ContainsRune("0123456789ABCDEFGHJKMNPQRSTVWXYZ", c) {
						t.Errorf("invalid character in entropy: %c", c)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewAPIKeyFromBytes(tt.prefix, uuid, tt.options...)
			if err != nil {
				t.Fatalf("NewAPIKeyFromBytes() failed: %v", err)
			}

			// Validate entropy length
			if len(key.Entropy) != tt.wantLen {
				t.Errorf("entropy length mismatch, got %d, want %d", len(key.Entropy), tt.wantLen)
			}

			// Run additional validations
			tt.validate(t, key)

			// Verify the key can be parsed back
			str := key.String()
			parts := strings.Split(str, "_")
			if len(parts) != 3 {
				t.Errorf("String() result should have 3 parts (prefix_key_checksum), got %d parts", len(parts))
			}
			if len(parts[2]) != 8 {
				t.Errorf("String() checksum length mismatch, got %d, want 8", len(parts[2]))
			}

			parsed, err := ParseAPIKey(str)
			if err != nil {
				t.Fatalf("ParseAPIKey() failed: %v", err)
			}
			if parsed.String() != str {
				t.Errorf("ParseAPIKey() result mismatch, got %s, want %s", parsed.String(), str)
			}
		})
	}
}

func TestParseAPIKeyErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "Empty string",
			input:   "",
			wantErr: "invalid APIKey format: expected 3 parts, got 1",
		},
		{
			name:    "Missing checksum",
			input:   "TEST_38QARV01ET0G6Z2CJD9VA2ZZAR0X",
			wantErr: "invalid APIKey format: expected 3 parts, got 2",
		},
		{
			name:    "Too many parts",
			input:   "TEST_38QARV01ET0G6Z2CJD9VA2ZZAR0X_ABC_DEF",
			wantErr: "invalid APIKey format: expected 3 parts, got 4",
		},
		{
			name:    "Invalid prefix",
			input:   "_38QARV01ET0G6Z2CJD9VA2ZZAR0X_12345678",
			wantErr: "invalid prefix: cannot be empty",
		},
		{
			name:    "Invalid key length",
			input:   "TEST_INVALID_12345678",
			wantErr: "invalid Key format: insufficient length",
		},
		{
			name: "Invalid characters in key",
			// Using proper length but invalid characters
			input:   "TEST_38QARV01ET0G6Z2CJD9VA2ZZ!R0XAAAAAA_12345678",
			wantErr: "invalid Key format: invalid UUID Key",
		},
		{
			name:    "Invalid checksum length",
			input:   "TEST_38QARV01ET0G6Z2CJD9VA2ZZAR0X_123",
			wantErr: "invalid checksum format: must be 8 hexadecimal characters",
		},
		{
			name:    "Invalid checksum value",
			input:   "TEST_38QARV01ET0G6Z2CJD9VA2ZZAR0X_12345678",
			wantErr: "invalid checksum: expected ", // actual value will vary
		},
		{
			name: "Invalid UUID format in key",
			// Using invalid characters that will fail UUID parsing
			input:   "TEST_IIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIII_12345678",
			wantErr: "invalid Key format: invalid UUID Key",
		},
		{
			name:    "Invalid checksum characters (lowercase)",
			input:   "TEST_38QARV01ET0G6Z2CJD9VA2ZZAR0X_1234abcd",
			wantErr: "invalid checksum format: must be 8 hexadecimal characters",
		},
		{
			name:    "Invalid checksum characters (non-hex)",
			input:   "TEST_38QARV01ET0G6Z2CJD9VA2ZZAR0X_1234WXYZ",
			wantErr: "invalid checksum format: must be 8 hexadecimal characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseAPIKey(tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.HasPrefix(err.Error(), tt.wantErr) {
				t.Errorf("wrong error message, got %q, want prefix %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestEntropyCharacterValidation(t *testing.T) {
	tests := []struct {
		name    string
		prefix  string
		options []Option
	}{
		{
			name:    "128-bit character validation",
			prefix:  "TEST",
			options: []Option{With128BitEntropy},
		},
		{
			name:    "160-bit (default) character validation",
			prefix:  "TEST",
			options: []Option{},
		},
		{
			name:    "256-bit character validation",
			prefix:  "TEST",
			options: []Option{With256BitEntropy},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewAPIKey(tt.prefix, "d1756360-5da0-40df-9926-a76abff5601d", tt.options...)
			if err != nil {
				t.Fatalf("NewAPIKey() failed: %v", err)
			}

			// Validate all characters are in Base32-Crockford alphabet
			for _, c := range key.Entropy {
				if !strings.ContainsRune("0123456789ABCDEFGHJKMNPQRSTVWXYZ", c) {
					t.Errorf("invalid character in entropy: %c", c)
				}
			}

			// Validate checksum is valid uppercase hexadecimal (CRC32)
			str := key.String()
			parts := strings.Split(str, "_")
			checksum := parts[2]
			if len(checksum) != 8 {
				t.Errorf("checksum length must be 8, got %d", len(checksum))
			}
			for _, c := range checksum {
				if !strings.ContainsRune("0123456789ABCDEF", c) {
					t.Errorf("invalid character in checksum: %c (must be uppercase hexadecimal)", c)
				}
			}
		})
	}
}

func TestChecksumConsistency(t *testing.T) {
	// Use a fixed UUID for consistency
	uuid := [16]byte{
		0xd1, 0x75, 0x63, 0x60,
		0x5d, 0xa0, 0x40, 0xdf,
		0x99, 0x26, 0xa7, 0x6a,
		0xbf, 0xf5, 0x60, 0x1d,
	}

	// Create first key
	key1, err := NewAPIKeyFromBytes("TEST", uuid, With128BitEntropy)
	if err != nil {
		t.Fatalf("NewAPIKeyFromBytes() failed: %v", err)
	}

	// Create second key with same components
	key2 := APIKey{
		Prefix:  key1.Prefix,
		Key:     key1.Key,
		Entropy: key1.Entropy,
	}
	key2.Checksum = key2.calculateChecksum()

	// Verify the checksums match
	if key1.Checksum != key2.Checksum {
		t.Errorf("checksums don't match for identical keys: %s != %s", key1.Checksum, key2.Checksum)
	}

	// Verify checksum is recalculated correctly
	str := key1.String()
	parsed, err := ParseAPIKey(str)
	if err != nil {
		t.Fatalf("ParseAPIKey() failed: %v", err)
	}
	if parsed.Checksum != key1.Checksum {
		t.Errorf("checksum mismatch after parsing: %s != %s", parsed.Checksum, key1.Checksum)
	}
}
