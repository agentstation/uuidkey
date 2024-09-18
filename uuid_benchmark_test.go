//go:build bench
// +build bench

package uuidkey_test

import (
	"testing"

	"github.com/agentstation/uuidkey"
)

const (
	validKey   = "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	invalidKey = "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0" // Invalid key
	validUUID  = "d1756360-5da0-40df-9926-a76abff5601d"
)

func BenchmarkValidate(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = key.Valid()
	}
}

func BenchmarkValidateInvalid(b *testing.B) {
	key := uuidkey.Key(invalidKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = key.Valid()
	}
}

func BenchmarkParse(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.Parse(validKey)
	}
}

func BenchmarkParseInvalid(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.Parse(invalidKey)
	}
}

func BenchmarkUUID(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.UUID()
	}
}

func BenchmarkUUIDInvalid(b *testing.B) {
	key := uuidkey.Key(invalidKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.UUID()
	}
}

func BenchmarkEncode(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.Encode(validUUID)
	}
}

func BenchmarkDecode(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.Decode()
	}
}

func BenchmarkBytes(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.Bytes()
	}
}

func BenchmarkEncodeBytes(b *testing.B) {
	uuid := [16]byte{
		0xd1, 0x75, 0x63, 0x60,
		0x5d, 0xa0, 0x40, 0xdf,
		0x99, 0x26, 0xa7, 0x6a,
		0xbf, 0xf5, 0x60, 0x1d,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.EncodeBytes(uuid)
	}
}
