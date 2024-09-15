//go:build bench
// +build bench

package uuidkey_test

import (
	"testing"

	"github.com/agentstation/uuidkey"
)

func BenchmarkValidate(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = key.Valid()
	}
}

func BenchmarkParse(b *testing.B) {
	s := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.Parse(s)
	}
}

func BenchmarkFromKey(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.UUID()
	}
}

func BenchmarkEncode(b *testing.B) {
	uuid := "d1756360-5da0-40df-9926-a76abff5601d"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.Encode(uuid)
	}
}

func BenchmarkDecode(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.Decode()
	}
}

func BenchmarkValidateInvalid(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0") // Invalid key
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = key.Valid()
	}
}

func BenchmarkParseValid(b *testing.B) {
	s := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.Parse(s)
	}
}

func BenchmarkParseInvalid(b *testing.B) {
	s := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0" // Invalid key
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.Parse(s)
	}
}

func BenchmarkUUIDValid(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.UUID()
	}
}

func BenchmarkUUIDInvalid(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0") // Invalid key
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.UUID()
	}
}
