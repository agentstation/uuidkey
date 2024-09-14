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

func BenchmarkFromString(b *testing.B) {
	s := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.FromString(s)
	}
}

func BenchmarkFromKey(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.UUIDString()
	}
}

func BenchmarkEncode(b *testing.B) {
	uuid := "d1756360-5da0-40df-9926-a76abff5601d"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uuidkey.Encode(uuid)
	}
}

func BenchmarkDecode(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key.Decode()
	}
}

func BenchmarkValidateInvalid(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0") // Invalid key
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key.Valid()
	}
}

func BenchmarkFromStringValid(b *testing.B) {
	s := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.FromString(s)
	}
}

func BenchmarkFromStringInvalid(b *testing.B) {
	s := "38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0" // Invalid key
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uuidkey.FromString(s)
	}
}

func BenchmarkUUIDStringValid(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.UUIDString()
	}
}

func BenchmarkUUIDStringInvalid(b *testing.B) {
	key := uuidkey.Key("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0") // Invalid key
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = key.UUIDString()
	}
}
