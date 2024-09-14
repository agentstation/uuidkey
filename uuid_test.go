//go:build test
// +build test

package uuidkey_test

import (
	"testing"

	"github.com/agentstation/uuidkey"

	frsuuid "github.com/gofrs/uuid" // test-only dependency
	guuid "github.com/google/uuid"  // test-only dependency
)

func TestValid(t *testing.T) {
	validKeys := []uuidkey.Key{
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X",
		"0000000-0000000-0000000-0000000",
		"ZZZZZZZ-ZZZZZZZ-ZZZZZZZ-ZZZZZZZ",
	}
	invalidKeys := []uuidkey.Key{
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0",   // Too short
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0XX", // Too long
		"38qarv0-1ET0G6Z-2CJD9VA-2ZZAR0X",  // Lowercase
		"38QARV0 1ET0G6Z 2CJD9VA 2ZZAR0X",  // Spaces instead of hyphens
		"38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0!",  // Invalid character
		"38QARV0-1ET0G6-2CJD9VA-2ZZAR0X",   // Part too short
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

func TestFromString(t *testing.T) {
	validKey := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	k, err := uuidkey.FromString(validKey)
	if err != nil {
		t.Errorf("FromString() returned an error for valid key: %v", err)
	}
	if k != uuidkey.Key(validKey) {
		t.Errorf("FromString() returned incorrect key. Got %s, want %s", k, validKey)
	}

	invalidKey := "invalid-key"
	_, err = uuidkey.FromString(invalidKey)
	if err == nil {
		t.Errorf("FromString() did not return an error for invalid key")
	}
}

func TestUUIDString(t *testing.T) {
	validKey := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	expectedUUID := "d1756360-5da0-40df-9926-a76abff5601d"

	uuidStr, err := validKey.UUIDString()
	if err != nil {
		t.Errorf("UUIDString() returned an error for valid key: %v", err)
	}
	if uuidStr != expectedUUID {
		t.Errorf("UUIDString() returned incorrect UUID string. Got %s, want %s", uuidStr, expectedUUID)
	}

	invalidKey := uuidkey.Key("invalid-key")
	_, err = invalidKey.UUIDString()
	if err == nil {
		t.Errorf("UUIDString() did not return an error for invalid key")
	}
}

func TestEncodeDecode(t *testing.T) {
	uuidStr := "d1756360-5da0-40df-9926-a76abff5601d"
	key := uuidkey.Encode(uuidStr)
	decodedUUID := key.Decode()

	if decodedUUID != uuidStr {
		t.Errorf("Encode/Decode roundtrip failed. Got %s, want %s", decodedUUID, uuidStr)
	}
}

func TestGoogleUUIDRoundtrip(t *testing.T) {
	for i := 0; i < 1000; i++ { // Test with 1000 random UUIDs
		// Generate a random UUID using Google's library
		originalUUID := guuid.New()
		uuidString := originalUUID.String()

		// Encode the UUID to our custom key format
		key := uuidkey.Encode(uuidString)

		// Ensure the key is valid
		if !key.Valid() {
			t.Errorf("Generated key is not valid: %s", key)
			continue
		}

		// Decode the key back to a UUID string
		decodedUUIDString, err := key.UUIDString()
		if err != nil {
			t.Errorf("Error decoding key %s: %v", key, err)
			continue
		}

		// Parse the decoded UUID string back into a UUID object
		decodedUUID, err := guuid.Parse(decodedUUIDString)
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
		originalUUID, err := frsuuid.NewV4()
		if err != nil {
			t.Fatalf("Failed to generate UUID: %v", err)
		}
		uuidString := originalUUID.String()

		// Encode the UUID to our custom key format
		key := uuidkey.Encode(uuidString)

		// Ensure the key is valid
		if !key.Valid() {
			t.Errorf("Generated key is not valid: %s", key)
			continue
		}

		// Decode the key back to a UUID string
		decodedUUIDString, err := key.UUIDString()
		if err != nil {
			t.Errorf("Error decoding key %s: %v", key, err)
			continue
		}

		// Parse the decoded UUID string back into a UUID object
		decodedUUID, err := frsuuid.FromString(decodedUUIDString)
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