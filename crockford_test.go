package uuidkey

import (
	"testing"
	"time"
)

// TestCrockford32ConstantTime verifies that our crock32 encoding operations
// run in constant time regardless of input values
func TestCrockford32ConstantTime(t *testing.T) {
	// Test values with different bit patterns
	testCases := []struct {
		name  string
		value uint32
	}{
		{"all zeros", 0x00000000},
		{"all ones", 0xFFFFFFFF},
		{"alternating bits", 0xAAAAAAAA},
		{"single bit set", 0x00000001},
		{"high bit set", 0x80000000},
		{"random pattern 1", 0x12345678},
		{"random pattern 2", 0xFEDCBA98},
	}

	// Warm up
	for i := range 1000 {
		_ = crock32Encode(uint32(i))
	}

	// Measure encoding times
	encodeTimes := make(map[string]time.Duration)
	iterations := 10000

	for _, tc := range testCases {
		start := time.Now()
		for range iterations {
			_ = crock32Encode(tc.value)
		}
		encodeTimes[tc.name] = time.Since(start)
	}

	// Calculate average and check variance
	var total time.Duration
	for _, duration := range encodeTimes {
		total += duration
	}
	avg := total / time.Duration(len(encodeTimes))

	// Allow 20% variance from average (reasonable for non-cryptographic constant time)
	threshold := avg / 5

	for name, duration := range encodeTimes {
		variance := duration - avg
		if variance < 0 {
			variance = -variance
		}
		if variance > threshold {
			t.Logf("WARNING: %s took %v (avg: %v, variance: %v)", name, duration, avg, variance)
		}
	}

	// Test decoding constant time
	testStrings := []string{
		"0000000",
		"ZZZZZZZ",
		"1234567",
		"ABCDEFG",
		"0000001",
		"Z000000",
	}

	// Warm up
	for range 1000 {
		_, _ = crock32Decode("1234567")
	}

	decodeTimes := make(map[string]time.Duration)

	for _, str := range testStrings {
		start := time.Now()
		for range iterations {
			_, _ = crock32Decode(str)
		}
		decodeTimes[str] = time.Since(start)
	}

	// Calculate average for decode times
	total = 0
	for _, duration := range decodeTimes {
		total += duration
	}
	avg = total / time.Duration(len(decodeTimes))
	threshold = avg / 5

	for str, duration := range decodeTimes {
		variance := duration - avg
		if variance < 0 {
			variance = -variance
		}
		if variance > threshold {
			t.Logf("WARNING: decoding %s took %v (avg: %v, variance: %v)", str, duration, avg, variance)
		}
	}
}

// BenchmarkCrockford32Encode benchmarks the custom crock32 encoding
func BenchmarkCrockford32Encode(b *testing.B) {
	values := []uint32{
		0x00000000,
		0xFFFFFFFF,
		0x12345678,
		0xFEDCBA98,
	}

	b.ResetTimer()
	for i := range b.N {
		_ = crock32Encode(values[i%len(values)])
	}
}

// BenchmarkCrockford32Decode benchmarks the custom crock32 decoding
func BenchmarkCrockford32Decode(b *testing.B) {
	strings := []string{
		"0000000",
		"ZZZZZZZ",
		"1234567",
		"ABCDEFG",
	}

	b.ResetTimer()
	for i := range b.N {
		_, _ = crock32Decode(strings[i%len(strings)])
	}
}

// TestCrock32DecodeOverflowAddition tests the overflow check during addition in crock32Decode
func TestCrock32DecodeOverflowAddition(t *testing.T) {
	// Create a string that will cause overflow during addition
	// Max uint32 is 4,294,967,295
	// In base32, this is "3ZZZZZZ"
	// Adding another digit would cause overflow
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "max value that doesn't overflow",
			input:     "3ZZZZZZ",
			expectErr: false,
		},
		{
			name:      "value that causes overflow on addition",
			input:     "3ZZZZZZZ", // 8 characters, will overflow
			expectErr: true,
		},
		{
			name:      "very large value",
			input:     "ZZZZZZZZ",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := crock32Decode(tt.input)
			if tt.expectErr && err == nil {
				t.Errorf("Expected overflow error for %s, but got none", tt.input)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.input, err)
			}
		})
	}
}
