//go:build !integration
// +build !integration

package uuidkey

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	// test-only dependencies
	gofrsUUID "github.com/gofrs/uuid"
	googleUUID "github.com/google/uuid"
)

// test-only dependency

// TestValid tests the IsValid method for both with and without hyphens
func TestValid(t *testing.T) {
	validKeys := []Key{
		// with hyphens
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
		"0000000-0000000-0000000-0000000",
		"ZZZZZZZ-ZZZZZZZ-ZZZZZZZ-ZZZZZZZ",
		// no hyphens
		"38QARV01ET0G6Z2CJD9VA2ZZAR0X",
		"0000000000000000000000000000",
		"ZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
	}
	invalidKeys := []Key{
		// with hyphens
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0",   // Too short
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0XX", // Too long
		"38qarv0-1ET0G6Z-2CJD9VA-2ZZAR0X",  // Lowercase
		"38QARV0 1ET0G6Z 2CJD9VA 2ZZAR0X",  // Spaces instead of hyphens
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0!",  // Invalid character
		"38QARV0-1ET0G6-2CJD9VA-2ZZAR0X",   // Part too short
		"38QARV0-1ET0G6Z-2CJD9VAA-2ZZAR0",  // Third part too long
		"38QARV0-LET0G6Z-2CJD9VA-2ZZAROX",  // Contains non-crockford base32 characters
		// no hyphens
		"38QARV01ET0G6Z2CJD9VA2ZZAR0",   // Too short
		"38QARV01ET0G6Z2CJD9VA2ZZAR0XL", // Too long
		"38qarv01et0g6z2cjd9va2zzar0",   // Lowercase
		"38QARV01ET0G6Z2CJD9VA2ZZAR0!",  // Invalid character
		"38QARV01ET0G6Z2CJD9VA2ZZAR0XL", // Too long
		"38QARV0LET0G6Z2CJD9VA2ZZAR0X",  // Contains non-crockford base32 characters
		// Additional hyphen cases
		"38QARV0-1ET0G6Z2CJD9VA2ZZAR0X", // Only first hyphen

		// Additional non-hyphen cases
		"38QARV01ET0G6Z-2CJD9VA2ZZAR0X", // Unexpected hyphen in middle
		"38QARV01ET0G6Z2CJD9VA-2ZZAR0X", // Unexpected hyphen near end
		"38QARV01ET0G6Z2CJD9VA2ZZAR0X-", // Unexpected hyphen at end
		// Additional length validation cases
		"38QARV0-1ET0G6-2CJD9VA-2ZZAR0X",   // Second part too short (6 chars)
		"38QARV0-1ET0G6ZZ-2CJD9VA-2ZZAR0X", // Second part too long (8 chars)
		"38QAR-1ET0G6Z-2CJD9VA-2ZZAR0X",    // First part too short (5 chars)
		"38QARV0Z-1ET0G6Z-2CJD9VA-2ZZAR0X", // First part too long (8 chars)

		// Without hyphens length validation
		"38QAR01ET0G6Z2CJD9VA2ZZAR0X",   // First part too short (5 chars)
		"38QARV0Z1ET0G6Z2CJD9VA2ZZAR0X", // First part too long (8 chars)
		"38QARV01ET0G2CJD9VA2ZZAR0X",    // Second part too short (5 chars)
		"38QARV01ET0G6ZZ2CJD9VA2ZZAR0X", // Second part too long (8 chars)
	}

	for _, k := range validKeys {
		if !k.IsValid() {
			t.Errorf("Validate() incorrectly reported valid key as invalid: %s", k)
		}
	}

	for _, k := range invalidKeys {
		if k.IsValid() {
			t.Errorf("Validate() incorrectly reported invalid key as valid: %s", k)
		}
	}
}

// TestParse tests the Parse method for both with and without hyphens
func TestParse(t *testing.T) {
	validKeyWithHyphens := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	k, err := Parse(validKeyWithHyphens)
	if err != nil {
		t.Errorf("Parse() returned an error for valid key: %v", err)
	}
	if k != Key(validKeyWithHyphens) {
		t.Errorf("Parse() returned incorrect key. Got %s, want %s", k, validKeyWithHyphens)
	}

	validKeyWithoutHyphens := "38QARV01ET0G6Z2CJD9VA2ZZAR0X"
	k, err = Parse(validKeyWithoutHyphens)
	if err != nil {
		t.Errorf("Parse() returned an error for valid key: %v", err)
	}
	if k != Key(validKeyWithoutHyphens) {
		t.Errorf("Parse() returned incorrect key. Got %s, want %s", k, validKeyWithoutHyphens)
	}

	invalidKey := "invalid-key"
	_, err = Parse(invalidKey)
	if err == nil {
		t.Errorf("Parse() did not return an error for invalid key")
	}
}

// TestEncodeDecode tests the Encode and Decode methods for both with and without hyphens
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

	key, err = Encode(uuidStr, WithoutHyphens)
	if err != nil {
		t.Fatalf("Encode() returned an unexpected error: %v", err)
	}
	decodedUUID, err = key.Decode()
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

// TestUUIDString tests the UUID method for both with and without hyphens
func TestUUIDString(t *testing.T) {
	validKeyWithHyphens := Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	validKeyWithoutHyphens := Key("38QARV01ET0G6Z2CJD9VA2ZZAR0X")
	expectedUUID := "d1756360-5da0-40df-9926-a76abff5601d"

	uuidStr, err := validKeyWithHyphens.UUID()
	if err != nil {
		t.Errorf("UUID() returned an error for valid key: %v", err)
	}
	if uuidStr != expectedUUID {
		t.Errorf("UUID() returned incorrect UUID string. Got %s, want %s", uuidStr, expectedUUID)
	}

	uuidStr, err = validKeyWithoutHyphens.UUID()
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

// TestGoogleUUIDRoundtrip tests the roundtrip from Google's UUID library to our custom key format and back
func TestGoogleUUIDRoundtrip(t *testing.T) {
	for range 1000 { // Test with 1000 random UUIDs
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
		if !key.IsValid() {
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

// TestGoogleUUIDRoundtripWithoutHyphens tests the roundtrip from Google's UUID library to our custom key format and back without hyphens
func TestGoogleUUIDRoundtripWithoutHyphens(t *testing.T) {
	for range 1000 { // Test with 1000 random UUIDs
		// Generate a random UUID using Google's library
		originalUUID := googleUUID.New()
		uuidString := originalUUID.String()

		// Encode the UUID to our custom key format
		key, err := Encode(uuidString, WithoutHyphens)
		if err != nil {
			t.Errorf("Error encoding UUID %s: %v", uuidString, err)
			continue
		}

		// Ensure the key is valid
		if !key.IsValid() {
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

// TestGofrsUUIDRoundtrip tests the roundtrip from gofrs/uuid library to our custom key format and back
func TestGofrsUUIDRoundtrip(t *testing.T) {
	for range 1000 { // Test with 1000 random UUIDs
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
		if !key.IsValid() {
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

// TestGofrsUUIDRoundtripWithoutHyphens tests the roundtrip from gofrs/uuid library to our custom key format and back without hyphens
func TestGofrsUUIDRoundtripWithoutHyphens(t *testing.T) {
	for range 1000 { // Test with 1000 random UUIDs
		// Generate a random UUID using gofrs/uuid library
		originalUUID, err := gofrsUUID.NewV4()
		if err != nil {
			t.Fatalf("Failed to generate UUID: %v", err)
		}
		uuidString := originalUUID.String()

		// Encode the UUID to our custom key format
		key, err := Encode(uuidString, WithoutHyphens)
		if err != nil {
			t.Errorf("Error encoding UUID %s: %v", uuidString, err)
			continue
		}

		// Ensure the key is valid
		if !key.IsValid() {
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

// TestKeyString tests the String method for both with and without hyphens
func TestKeyString(t *testing.T) {
	keyWithHyphens := Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	keyWithoutHyphens := Key("38QARV01ET0G6Z2CJD9VA2ZZAR0X")
	expectedWithHyphens := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	expectedWithoutHyphens := "38QARV01ET0G6Z2CJD9VA2ZZAR0X"

	result := keyWithHyphens.String()

	if result != expectedWithHyphens {
		t.Errorf("Key.String() returned incorrect value. Got %s, want %s", result, expectedWithHyphens)
	}

	result = keyWithoutHyphens.String()

	if result != expectedWithoutHyphens {
		t.Errorf("Key.String() returned incorrect value. Got %s, want %s", result, expectedWithoutHyphens)
	}
}

// TestEncodeBytes tests the EncodeBytes method for both with and without hyphens
func TestEncodeBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   [16]byte
		want    Key
		wantErr bool
	}{
		{
			name:    "Valid UUID with hyphens",
			input:   [16]byte{0xd1, 0x75, 0x63, 0x60, 0x5d, 0xa0, 0x40, 0xdf, 0x99, 0x26, 0xa7, 0x6a, 0xbf, 0xf5, 0x60, 0x1d},
			want:    "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
			wantErr: false,
		},
		{
			name:    "All zeros with hyphens",
			input:   [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:    "0000000-0000000-0000000-0000000",
			wantErr: false,
		},
		{
			name:    "All ones with hyphens",
			input:   [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			want:    "3ZZZZZZ-3ZZZZZZ-3ZZZZZZ-3ZZZZZZ",
			wantErr: false,
		},
		{
			name:    "Valid UUID without hyphens",
			input:   [16]byte{0xd1, 0x75, 0x63, 0x60, 0x5d, 0xa0, 0x40, 0xdf, 0x99, 0x26, 0xa7, 0x6a, 0xbf, 0xf5, 0x60, 0x1d},
			want:    "38QARV01ET0G6Z2CJD9VA2ZZAR0X",
			wantErr: false,
		},
		{
			name:    "All zeros without hyphens",
			input:   [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:    "0000000000000000000000000000",
			wantErr: false,
		},
		{
			name:    "All ones without hyphens",
			input:   [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			want:    "3ZZZZZZ3ZZZZZZ3ZZZZZZ3ZZZZZZ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []Option
			if strings.Contains(tt.name, "without hyphens") {
				opts = append(opts, WithoutHyphens)
			}
			got, err := EncodeBytes(tt.input, opts...)
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

// TestEncodeBytesRoundTripGoogle tests the roundtrip from Google's UUID library to our custom key format and back
func TestEncodeBytesRoundTripGoogle(t *testing.T) {
	for range 1000 { // Test with 1000 random UUIDs
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
		if !key.IsValid() {
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

// TestEncodeBytesRoundTripGoogleWithoutHyphens tests the roundtrip from Google's UUID library to our custom key format and back without hyphens
func TestEncodeBytesRoundTripGoogleWithoutHyphens(t *testing.T) {
	for range 1000 { // Test with 1000 random UUIDs
		// Generate a random UUID using Google's library
		originalUUID := googleUUID.New()
		var uuidBytes [16]byte
		copy(uuidBytes[:], originalUUID[:])

		// Encode the UUID bytes to our custom key format
		key, err := EncodeBytes(uuidBytes, WithoutHyphens)
		if err != nil {
			t.Errorf("Error encoding UUID bytes %v: %v", uuidBytes, err)
			continue
		}

		// Ensure the key is valid
		if !key.IsValid() {
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

// TestEncodeBytesRoundTripGofrs tests the roundtrip from gofrs/uuid library to our custom key format and back
func TestEncodeBytesRoundTripGofrs(t *testing.T) {
	for range 1000 { // Test with 1000 random UUIDs
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
		if !key.IsValid() {
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

// TestEncodeBytesRoundTripGofrsWithoutHyphens tests the roundtrip from gofrs/uuid library to our custom key format and back without hyphens
func TestEncodeBytesRoundTripGofrsWithoutHyphens(t *testing.T) {
	for range 1000 { // Test with 1000 random UUIDs
		// Generate a random UUID using gofrs/uuid library
		originalUUID, err := gofrsUUID.NewV4()
		if err != nil {
			t.Fatalf("Failed to generate UUID: %v", err)
		}
		var uuidBytes [16]byte
		copy(uuidBytes[:], originalUUID[:])

		// Encode the UUID bytes to our custom key format
		key, err := EncodeBytes(uuidBytes, WithoutHyphens)
		if err != nil {
			t.Errorf("Error encoding UUID bytes %v: %v", uuidBytes, err)
			continue
		}

		// Ensure the key is valid
		if !key.IsValid() {
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

// TestKeyBytes tests the Bytes method for both with and without hyphens
func TestKeyBytes(t *testing.T) {
	tests := []struct {
		name    string
		key     Key
		want    [16]byte
		wantErr bool
	}{
		{
			name:    "Valid Key with hyphens",
			key:     "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
			want:    [16]byte{0xd1, 0x75, 0x63, 0x60, 0x5d, 0xa0, 0x40, 0xdf, 0x99, 0x26, 0xa7, 0x6a, 0xbf, 0xf5, 0x60, 0x1d},
			wantErr: false,
		},
		{
			name:    "All Zeros with hyphens",
			key:     "0000000-0000000-0000000-0000000",
			want:    [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "Invalid Key (Too Short) with hyphens",
			key:     "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Valid Key without hyphens",
			key:     "38QARV01ET0G6Z2CJD9VA2ZZAR0X",
			want:    [16]byte{0xd1, 0x75, 0x63, 0x60, 0x5d, 0xa0, 0x40, 0xdf, 0x99, 0x26, 0xa7, 0x6a, 0xbf, 0xf5, 0x60, 0x1d},
			wantErr: false,
		},
		{
			name:    "All Zeros without hyphens",
			key:     "0000000000000000000000000000",
			want:    [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name:    "Invalid Key (Too Short) without hyphens",
			key:     "38QARV01ET0G6Z2CJD9VA2ZZAR0",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Key with hyphens",
			key:     "INVALID-KEY",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Base32 characters with hyphens",
			key:     "######-1111111-1111111-1111111", // '#' is invalid in Base32
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Base32 characters without hyphens",
			key:     "#######1111111111111111111111",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Base32 characters in second group with hyphens",
			key:     "0000000-#######-1111111-1111111",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Base32 characters in third group with hyphens",
			key:     "0000000-1111111-#######-1111111",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Base32 characters in fourth group with hyphens",
			key:     "0000000-1111111-1111111-#######",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Base32 characters in second group without hyphens",
			key:     "0000000#######1111111111111",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Base32 characters in third group without hyphens",
			key:     "00000001111111#######111111",
			want:    [16]byte{},
			wantErr: true,
		},
		{
			name:    "Invalid Base32 characters in fourth group without hyphens",
			key:     "0000000111111111111111######",
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

func TestDecodeErrorMessage(t *testing.T) {
	tests := []struct {
		name          string
		key           Key
		expectedError string
	}{
		{
			name:          "Too short key",
			key:           "ABC",
			expectedError: "invalid Key length: expected 28 or 31 characters, got 3",
		},
		{
			name:          "Too long key",
			key:           "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0XX",
			expectedError: "invalid Key length: expected 28 or 31 characters, got 32",
		},
		{
			name:          "Empty key",
			key:           "",
			expectedError: "invalid Key length: expected 28 or 31 characters, got 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.key.Decode()
			if err == nil {
				t.Error("Expected error, got nil")
				return
			}
			if err.Error() != tt.expectedError {
				t.Errorf("Expected error message %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}

func TestKey_Bytes_Errors(t *testing.T) {
	tests := []struct {
		name    string
		key     Key
		wantErr string
	}{
		{
			name:    "invalid length - too short",
			key:     Key("ABC"),
			wantErr: "invalid Key length: expected 28 or 31 characters, got 3",
		},
		{
			name:    "invalid characters in first group",
			key:     Key("@#$%^&*-1111111-2222222-3333333"),
			wantErr: "failed to decode Key part: crock32.Decode: invalid character @",
		},
		{
			name:    "invalid characters in second group",
			key:     Key("1111111-@#$%^&*-2222222-3333333"),
			wantErr: "failed to decode Key part: crock32.Decode: invalid character @",
		},
		{
			name:    "invalid characters in third group",
			key:     Key("1111111-2222222-@#$%^&*-3333333"),
			wantErr: "failed to decode Key part: crock32.Decode: invalid character @",
		},
		{
			name:    "invalid characters in fourth group",
			key:     Key("1111111-2222222-3333333-@#$%^&*"),
			wantErr: "failed to decode Key part: crock32.Decode: invalid character @",
		},
		{
			name:    "invalid characters without hyphens",
			key:     Key("1111111222222233333333@#$%^&"),
			wantErr: "failed to decode Key part: crock32.Decode: invalid character @",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.key.Bytes()
			if err == nil {
				t.Error("expected error, got nil")
				return
			}
			if err.Error() != tt.wantErr {
				t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// TestIsValidPartEdgeCases tests edge cases for isValidPart to improve coverage
func TestIsValidPartEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		part  string
		valid bool
	}{
		// Characters just outside valid ranges
		{"char before 0", "//////1", false},        // '/' is just before '0'
		{"char after Z", "ABCDE[F", false},         // '[' is just after 'Z'
		{"char between 9 and A", "12345:6", false}, // ':' is between '9' and 'A'

		// Invalid Crockford characters
		{"contains I", "ABCDIEF", false},
		{"contains L", "ABCDLEF", false},
		{"contains O", "ABCDOEF", false},
		{"contains U", "ABCDUEF", false},

		// Valid edge cases
		{"all zeros", "0000000", true},
		{"all nines", "9999999", true},
		{"all As", "AAAAAAA", true},
		{"all Zs", "ZZZZZZZ", true},
		{"mixed valid", "0A9Z1B8", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPart(tt.part)
			if result != tt.valid {
				t.Errorf("isValidPart(%s) = %v, want %v", tt.part, result, tt.valid)
			}
		})
	}
}

// TestBytesErrorHandling tests error paths in the Bytes function for coverage
func TestBytesErrorHandling(t *testing.T) {
	// Create keys that will cause processByteGroup to fail
	tests := []struct {
		name    string
		key     Key
		wantErr bool
	}{
		{
			// This key has valid length but contains characters that will fail in processByteGroup
			name:    "invalid characters in first group with hyphens",
			key:     Key("!!!!!!!-1ET0G6Z-2CJD9VA-2ZZAR0X"),
			wantErr: true,
		},
		{
			name:    "invalid characters in second group with hyphens",
			key:     Key("38QARV0-!!!!!!!-2CJD9VA-2ZZAR0X"),
			wantErr: true,
		},
		{
			name:    "invalid characters in third group with hyphens",
			key:     Key("38QARV0-1ET0G6Z-!!!!!!!-2ZZAR0X"),
			wantErr: true,
		},
		{
			name:    "invalid characters in fourth group with hyphens",
			key:     Key("38QARV0-1ET0G6Z-2CJD9VA-!!!!!!!"),
			wantErr: true,
		},
		{
			// Without hyphens
			name:    "invalid characters in first group no hyphens",
			key:     Key("!!!!!!!1ET0G6Z2CJD9VA2ZZAR0X"),
			wantErr: true,
		},
		{
			name:    "invalid characters in second group no hyphens",
			key:     Key("38QARV0!!!!!!!2CJD9VA2ZZAR0X"),
			wantErr: true,
		},
		{
			name:    "invalid characters in third group no hyphens",
			key:     Key("38QARV01ET0G6Z!!!!!!!2ZZAR0X"),
			wantErr: true,
		},
		{
			name:    "invalid characters in fourth group no hyphens",
			key:     Key("38QARV01ET0G6Z2CJD9VA!!!!!!!"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.key.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.Bytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErr {
				// Verify it's the expected error from processByteGroup
				expectedErrMsg := "failed to decode Key part"
				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Errorf("Expected error message to contain %q, got %q", expectedErrMsg, err.Error())
				}
			}
		})
	}
}

// TestDecodeErrorPath tests the error handling path in the decode function
// by using keys that pass length validation but contain values that cause overflow
func TestDecodeErrorPath(t *testing.T) {
	// Create a key with all Z's that will cause overflow in some parts
	// ZZZZZZZZ in base32 would overflow uint32 (max 4,294,967,295)
	overflowKey := Key("ZZZZZZZ-ZZZZZZZ-ZZZZZZZ-ZZZZZZZ")

	// This key has valid length and format, so it will pass initial validation
	// but will trigger the error path in decode() due to overflow
	result, err := overflowKey.Decode()

	// We should get a valid UUID format back (with fallback zeros)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// The result should contain the fallback "00000000" for the overflow parts
	if !strings.Contains(result, "00000000") {
		t.Errorf("Expected result to contain fallback zeros, got %s", result)
	}

	// Test without hyphens
	overflowKeyNoHyphens := Key("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZ")
	result2, err2 := overflowKeyNoHyphens.Decode()

	if err2 != nil {
		t.Errorf("Unexpected error: %v", err2)
	}

	if !strings.Contains(result2, "00000000") {
		t.Errorf("Expected result to contain fallback zeros, got %s", result2)
	}
}

// TestDecodeWithInvalidCrock32 specifically tests decode function with invalid input
func TestDecodeWithInvalidCrock32(t *testing.T) {
	// Test the decode function directly with invalid input that triggers the error path
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "overflow input",
			input:    "ZZZZZZZZ", // 8 Z's would overflow uint32
			expected: "00000000",
		},
		{
			name:     "invalid characters",
			input:    "!@#$%^&",
			expected: "00000000",
		},
		{
			name:     "mixed invalid",
			input:    "ABC!DEF",
			expected: "00000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decode(tt.input)
			if result != tt.expected {
				t.Errorf("decode(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
