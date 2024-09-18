//go:build !integration
// +build !integration

package uuidkey

import (
	"bytes"
	"reflect"
	"testing"

	// test-only dependencies
	gofrsUUID "github.com/gofrs/uuid"
	googleUUID "github.com/google/uuid"
)

// test-only dependency

func TestValid(t *testing.T) {
	validKeys := []Key{
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
		"0000000-0000000-0000000-0000000",
		"ZZZZZZZ-ZZZZZZZ-ZZZZZZZ-ZZZZZZZ",
	}
	invalidKeys := []Key{
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0",   // Too short
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0XX", // Too long
		"38qarv0-1ET0G6Z-2CJD9VA-2ZZAR0X",  // Lowercase
		"38QARV0 1ET0G6Z 2CJD9VA 2ZZAR0X",  // Spaces instead of hyphens
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0!",  // Invalid character
		"38QARV0-1ET0G6-2CJD9VA-2ZZAR0X",   // Part too short
		"38QARV0-1ET0G6Z-2CJD9VAA-2ZZAR0",  // Third part too long
	}

	for _, k := range validKeys {
		if !k.Valid() {
			t.Errorf("Validate() incorrectly reported valid key as invalid: %s", k)
		}
	}

	for _, k := range invalidKeys {
		if k.Valid() {
			t.Errorf("Validate() incorrectly reported invalid key as valid: %s", k)
		}
	}
}

func TestParse(t *testing.T) {
	validKey := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	k, err := Parse(validKey)
	if err != nil {
		t.Errorf("Parse() returned an error for valid key: %v", err)
	}
	if k != Key(validKey) {
		t.Errorf("Parse() returned incorrect key. Got %s, want %s", k, validKey)
	}

	invalidKey := "invalid-key"
	_, err = Parse(invalidKey)
	if err == nil {
		t.Errorf("Parse() did not return an error for invalid key")
	}
}

func TestEncodeDecode(t *testing.T) {
	uuidStr := "d1756360-5da0-40df-9926-a76abff5601d"
	key, err := Encode(uuidStr)
	if err != nil {
		t.Fatalf("Encode() returned an unexpected error: %v", err)
	}
	decodedUUID, err := key.Decode()
	if err != nil {
		t.Fatalf("Decode() returned an unexpected error: %v", err)
	}

	if decodedUUID != uuidStr {
		t.Errorf("Encode/Decode roundtrip failed. Got %s, want %s", decodedUUID, uuidStr)
	}

	// Test invalid UUID length
	invalidUUID := "invalid-uuid"
	_, err = Encode(invalidUUID)
	if err == nil {
		t.Errorf("Encode() did not return an error for invalid UUID length")
	}
}

func TestUUIDString(t *testing.T) {
	validKey := Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	expectedUUID := "d1756360-5da0-40df-9926-a76abff5601d"

	uuidStr, err := validKey.UUID()
	if err != nil {
		t.Errorf("UUID() returned an error for valid key: %v", err)
	}
	if uuidStr != expectedUUID {
		t.Errorf("UUID() returned incorrect UUID string. Got %s, want %s", uuidStr, expectedUUID)
	}

	invalidKey := Key("invalid-key")
	_, err = invalidKey.UUID()
	if err == nil {
		t.Errorf("UUID() did not return an error for invalid key")
	}
}

func TestGoogleUUIDRoundtrip(t *testing.T) {
	for i := 0; i < 1000; i++ { // Test with 1000 random UUIDs
		// Generate a random UUID using Google's library
		originalUUID := googleUUID.New()
		uuidString := originalUUID.String()

		// Encode the UUID to our custom key format
		key, err := Encode(uuidString)
		if err != nil {
			t.Errorf("Error encoding UUID %s: %v", uuidString, err)
			continue
		}

		// Ensure the key is valid
		if !key.Valid() {
			t.Errorf("Generated key is not valid: %s", key)
			continue
		}

		// Decode the key back to a UUID string
		decodedUUIDString, err := key.UUID()
		if err != nil {
			t.Errorf("Error decoding key %s: %v", key, err)
			continue
		}

		// Parse the decoded UUID string back into a UUID object
		decodedUUID, err := googleUUID.Parse(decodedUUIDString)
		if err != nil {
			t.Errorf("Error parsing decoded UUID string %s: %v", decodedUUIDString, err)
			continue
		}

		// Compare the original and decoded UUIDs
		if originalUUID != decodedUUID {
			t.Errorf("UUID mismatch. Original: %s, Decoded: %s", originalUUID, decodedUUID)
		}
	}
}

func TestGofrsUUIDRoundtrip(t *testing.T) {
	for i := 0; i < 1000; i++ { // Test with 1000 random UUIDs
		// Generate a random UUID using gofrs/uuid library
		originalUUID, err := gofrsUUID.NewV4()
		if err != nil {
			t.Fatalf("Failed to generate UUID: %v", err)
		}
		uuidString := originalUUID.String()

		// Encode the UUID to our custom key format
		key, err := Encode(uuidString)
		if err != nil {
			t.Errorf("Error encoding UUID %s: %v", uuidString, err)
			continue
		}

		// Ensure the key is valid
		if !key.Valid() {
			t.Errorf("Generated key is not valid: %s", key)
			continue
		}

		// Decode the key back to a UUID string
		decodedUUIDString, err := key.UUID()
		if err != nil {
			t.Errorf("Error decoding key %s: %v", key, err)
			continue
		}

		// Parse the decoded UUID string back into a UUID object
		decodedUUID, err := gofrsUUID.FromString(decodedUUIDString)
		if err != nil {
			t.Errorf("Error parsing decoded UUID string %s: %v", decodedUUIDString, err)
			continue
		}

		// Compare the original and decoded UUIDs
		if originalUUID != decodedUUID {
			t.Errorf("UUID mismatch. Original: %s, Decoded: %s", originalUUID, decodedUUID)
		}
	}
}

func TestKeyString(t *testing.T) {
	key := Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	expected := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"

	result := key.String()

	if result != expected {
		t.Errorf("Key.String() returned incorrect value. Got %s, want %s", result, expected)
	}
}

func TestEncodeBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   [16]byte
		want    Key
		wantErr bool
	}{
		{
			name:    "Valid UUID",
			input:   [16]byte{0xd1, 0x75, 0x63, 0x60, 0x5d, 0xa0, 0x40, 0xdf, 0x99, 0x26, 0xa7, 0x6a, 0xbf, 0xf5, 0x60, 0x1d},
			want:    "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
			wantErr: false,
		},
		{
			name:    "All zeros",
			input:   [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:    "0000000-0000000-0000000-0000000",
			wantErr: false,
		},
		{
			name:    "All ones",
			input:   [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			want:    "3ZZZZZZ-3ZZZZZZ-3ZZZZZZ-3ZZZZZZ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeBytesRoundTripGoogle(t *testing.T) {
	for i := 0; i < 1000; i++ { // Test with 1000 random UUIDs
		// Generate a random UUID using Google's library
		originalUUID := googleUUID.New()
		var uuidBytes [16]byte
		copy(uuidBytes[:], originalUUID[:])

		// Encode the UUID bytes to our custom key format
		key, err := EncodeBytes(uuidBytes)
		if err != nil {
			t.Errorf("Error encoding UUID bytes %v: %v", uuidBytes, err)
			continue
		}

		// Ensure the key is valid
		if !key.Valid() {
			t.Errorf("Generated key is not valid: %s", key)
			continue
		}

		// Convert the key back to bytes
		decodedBytes, err := key.Bytes()
		if err != nil {
			t.Errorf("Error converting key %s to bytes: %v", key, err)
			continue
		}

		// Compare the original and decoded UUID bytes
		if !bytes.Equal(uuidBytes[:], decodedBytes[:]) {
			t.Errorf("UUID bytes mismatch. Original: %v, Decoded: %v", uuidBytes, decodedBytes)
		}
	}
}

func TestEncodeBytesRoundTripGofrs(t *testing.T) {
	for i := 0; i < 1000; i++ { // Test with 1000 random UUIDs
		// Generate a random UUID using gofrs/uuid library
		originalUUID, err := gofrsUUID.NewV4()
		if err != nil {
			t.Fatalf("Failed to generate UUID: %v", err)
		}
		var uuidBytes [16]byte
		copy(uuidBytes[:], originalUUID[:])

		// Encode the UUID bytes to our custom key format
		key, err := EncodeBytes(uuidBytes)
		if err != nil {
			t.Errorf("Error encoding UUID bytes %v: %v", uuidBytes, err)
			continue
		}

		// Ensure the key is valid
		if !key.Valid() {
			t.Errorf("Generated key is not valid: %s", key)
			continue
		}

		// Convert the key back to bytes
		decodedBytes, err := key.Bytes()
		if err != nil {
			t.Errorf("Error converting key %s to bytes: %v", key, err)
			continue
		}

		// Compare the original and decoded UUID bytes
		if !bytes.Equal(uuidBytes[:], decodedBytes[:]) {
			t.Errorf("UUID bytes mismatch. Original: %v, Decoded: %v", uuidBytes, decodedBytes)
		}
	}
}

func TestKeyBytes(t *testing.T) {
	tests := []struct {
		name    string
		key     Key
		want    [16]byte
		wantErr bool
	}{
		{
			name:    "Valid Key",
			key:     "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
			want:    [16]byte{0xd1, 0x75, 0x63, 0x60, 0x5d, 0xa0, 0x40, 0xdf, 0x99, 0x26, 0xa7, 0x6a, 0xbf, 0xf5, 0x60, 0x1d},
			wantErr: false,
		},
		{
			name:    "All Zeros",
			key:     "0000000-0000000-0000000-0000000",
			want:    [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "Invalid Key",
			key:     "INVALID-KEY",
			want:    [16]byte{},
			wantErr: true,
		},
		// Add this new test case
		{
			name:    "Invalid Key (Too Short)",
			key:     "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0",
			want:    [16]byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.key.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key.Bytes() = %v, want %v", got, tt.want)
			}
			if err != nil && !tt.wantErr {
				t.Errorf("Key.Bytes() unexpected error: %v", err)
			}
			if tt.wantErr && err == nil {
				t.Errorf("Key.Bytes() expected error, got nil")
			}
		})
	}
}
