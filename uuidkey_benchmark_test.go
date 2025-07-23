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
	for range b.N {
		_ = key.IsValid()
	}
}

func BenchmarkValidateInvalid(b *testing.B) {
	key := uuidkey.Key(invalidKey)
	b.ResetTimer()
	for range b.N {
		_ = key.IsValid()
	}
}

func BenchmarkParse(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Parse(validKey)
	}
}

func BenchmarkParseInvalid(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Parse(invalidKey)
	}
}

func BenchmarkUUID(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for range b.N {
		_, _ = key.UUID()
	}
}

func BenchmarkUUIDInvalid(b *testing.B) {
	key := uuidkey.Key(invalidKey)
	b.ResetTimer()
	for range b.N {
		_, _ = key.UUID()
	}
}

func BenchmarkEncode(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Encode(validUUID)
	}
}

func BenchmarkDecode(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for range b.N {
		_, _ = key.Decode()
	}
}

func BenchmarkBytes(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for range b.N {
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
	for range b.N {
		_, _ = uuidkey.EncodeBytes(uuid)
	}
}

func BenchmarkValidateWithHyphens(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for range b.N {
		_ = key.IsValid()
	}
}

func BenchmarkValidateWithoutHyphens(b *testing.B) {
	key := uuidkey.Key("38QARV01ET0G6Z2CJD9VA2ZZAR0X")
	b.ResetTimer()
	for range b.N {
		_ = key.IsValid()
	}
}

func BenchmarkParseWithHyphens(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Parse(validKey)
	}
}

func BenchmarkParseWithoutHyphens(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Parse("38QARV01ET0G6Z2CJD9VA2ZZAR0X")
	}
}

func BenchmarkEncodeWithHyphens(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Encode(validUUID)
	}
}

func BenchmarkEncodeWithoutHyphens(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Encode(validUUID, uuidkey.WithoutHyphens)
	}
}

func BenchmarkDecodeWithHyphens(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for range b.N {
		_, _ = key.Decode()
	}
}

func BenchmarkDecodeWithoutHyphens(b *testing.B) {
	key := uuidkey.Key("38QARV01ET0G6Z2CJD9VA2ZZAR0X")
	b.ResetTimer()
	for range b.N {
		_, _ = key.Decode()
	}
}

func BenchmarkBytesWithHyphens(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for range b.N {
		_, _ = key.Bytes()
	}
}

func BenchmarkBytesWithoutHyphens(b *testing.B) {
	key := uuidkey.Key("38QARV01ET0G6Z2CJD9VA2ZZAR0X")
	b.ResetTimer()
	for range b.N {
		_, _ = key.Bytes()
	}
}

func BenchmarkEncodeBytesWithHyphens(b *testing.B) {
	uuid := [16]byte{
		0xd1, 0x75, 0x63, 0x60,
		0x5d, 0xa0, 0x40, 0xdf,
		0x99, 0x26, 0xa7, 0x6a,
		0xbf, 0xf5, 0x60, 0x1d,
	}
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.EncodeBytes(uuid)
	}
}

func BenchmarkEncodeBytesWithoutHyphens(b *testing.B) {
	uuid := [16]byte{
		0xd1, 0x75, 0x63, 0x60,
		0x5d, 0xa0, 0x40, 0xdf,
		0x99, 0x26, 0xa7, 0x6a,
		0xbf, 0xf5, 0x60, 0x1d,
	}
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.EncodeBytes(uuid, uuidkey.WithoutHyphens)
	}
}

func BenchmarkString(b *testing.B) {
	key := uuidkey.Key(validKey)
	b.ResetTimer()
	for range b.N {
		_ = key.String()
	}
}

func BenchmarkValidateInvalidFormat(b *testing.B) {
	key := uuidkey.Key("INVALID-FORMAT-KEY")
	b.ResetTimer()
	for range b.N {
		_ = key.IsValid()
	}
}

func BenchmarkParseInvalidFormat(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Parse("INVALID-FORMAT-KEY")
	}
}

func BenchmarkDecodeInvalidFormat(b *testing.B) {
	key := uuidkey.Key("INVALID-FORMAT-KEY")
	b.ResetTimer()
	for range b.N {
		_, _ = key.Decode()
	}
}

func BenchmarkEncodeInvalidUUID(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.Encode("not-a-valid-uuid")
	}
}

func BenchmarkBytesInvalidFormat(b *testing.B) {
	key := uuidkey.Key("INVALID-FORMAT-KEY")
	b.ResetTimer()
	for range b.N {
		_, _ = key.Bytes()
	}
}

func BenchmarkNewAPIKey(b *testing.B) {
	prefix := "TEST"
	uuid := "d1756360-5da0-40df-9926-a76abff5601d"
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.NewAPIKey(prefix, uuid)
	}
}

func BenchmarkNewAPIKeyWith128BitEntropy(b *testing.B) {
	prefix := "TEST"
	uuid := "d1756360-5da0-40df-9926-a76abff5601d"
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.NewAPIKey(prefix, uuid, uuidkey.With128BitEntropy)
	}
}

func BenchmarkNewAPIKeyWith256BitEntropy(b *testing.B) {
	prefix := "TEST"
	uuid := "d1756360-5da0-40df-9926-a76abff5601d"
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.NewAPIKey(prefix, uuid, uuidkey.With256BitEntropy)
	}
}

func BenchmarkNewAPIKeyFromBytes(b *testing.B) {
	prefix := "TEST"
	uuid := [16]byte{
		0xd1, 0x75, 0x63, 0x60,
		0x5d, 0xa0, 0x40, 0xdf,
		0x99, 0x26, 0xa7, 0x6a,
		0xbf, 0xf5, 0x60, 0x1d,
	}
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.NewAPIKeyFromBytes(prefix, uuid)
	}
}

func BenchmarkAPIKeyString(b *testing.B) {
	key, _ := uuidkey.NewAPIKey("TEST", "d1756360-5da0-40df-9926-a76abff5601d")
	b.ResetTimer()
	for range b.N {
		_ = key.String()
	}
}

func BenchmarkParseAPIKey(b *testing.B) {
	key, _ := uuidkey.NewAPIKey("TEST", "d1756360-5da0-40df-9926-a76abff5601d")
	apiKey := key.String()
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.ParseAPIKey(apiKey)
	}
}

func BenchmarkParseAPIKeyInvalid(b *testing.B) {
	apiKey := "INVALID_KEY_FORMAT"
	b.ResetTimer()
	for range b.N {
		_, _ = uuidkey.ParseAPIKey(apiKey)
	}
}
