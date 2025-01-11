// Package uuidkey encodes UUIDs to a readable Key format via the Base32-Crockford codec.
package uuidkey

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/richardlehane/crock32"
)

// Key validation constraint constants
const (
	// KeyLengthWithHyphens is the total length of a valid UUID Key, including hyphens.
	KeyLengthWithHyphens = 31 // 7 + 1 + 7 + 1 + 7 + 1 + 7 = 31 characters

	// KeyLengthWithoutHyphens is the total length of a valid UUID Key, excluding hyphens.
	KeyLengthWithoutHyphens = 28 // 7 + 7 + 7 + 7 = 28 characters

	// KeyPartLength is the length of each part in a UUID Key.
	// A UUID Key consists of 4 parts separated by hyphens.
	KeyPartLength = 7

	// KeyHyphenCount is the number of hyphens in a valid UUID Key.
	KeyHyphenCount = 3

	// KeyPartsCount is the number of parts in a valid UUID Key.
	KeyPartsCount = KeyHyphenCount + 1

	// UUIDLength is the standard length of a UUID string, including hyphens.
	// Reference: RFC 4122 (https://tools.ietf.org/html/rfc4122)
	UUIDLength = 36
)

// Key is a UUID Key string.
type Key string

// String will convert your Key into a string.
func (k Key) String() string {
	return string(k)
}

// Parse converts a Key formatted string into a Key type.
func Parse(key string) (Key, error) {
	k := Key(key)
	if !k.IsValid() {
		return "", errors.New("invalid UUID Key")
	}
	return k, nil
}

// IsValid verifies if a given Key follows the correct format.
// The format should be:
//   - 31 characters long (with hyphens) or 28 characters (without hyphens)
//   - Uppercase
//   - Contains only alphanumeric characters
//   - Contains 3 hyphens (if hyphenated)
//   - Each part is 7 characters long
//   - Each part contains only valid crockford base32 characters (I, L, O, U are not allowed)
func (k Key) IsValid() bool {
	length := len(k)
	if length == KeyLengthWithHyphens {
		// Check hyphens first (faster than character validation)
		if k[7] != '-' || k[15] != '-' || k[23] != '-' {
			return false
		}
		// Use direct string indexing instead of slicing
		return isValidPart(string(k[0:7])) && isValidPart(string(k[8:15])) &&
			isValidPart(string(k[16:23])) && isValidPart(string(k[24:31]))
	}
	if length == KeyLengthWithoutHyphens {
		// Use direct string indexing instead of slicing
		return isValidPart(string(k[0:7])) && isValidPart(string(k[7:14])) &&
			isValidPart(string(k[14:21])) && isValidPart(string(k[21:28]))
	}
	return false
}

// isValidPart checks if a 7-character part of the key is valid:
// - Must be exactly 7 characters
// - Must be uppercase alphanumeric
// - Must not contain I, L, O, U (invalid in crockford base32)
func isValidPart(part string) bool {
	if len(part) != KeyPartLength {
		return false
	}
	for i := 0; i < KeyPartLength; i++ {
		c := part[i]
		// Combine conditions to reduce branching
		if c > 'Z' || (c < '0' || (c > '9' && c < 'A')) ||
			c == 'I' || c == 'L' || c == 'O' || c == 'U' {
			return false
		}
	}
	return true
}

// UUID will validate and convert a given Key into a UUID string.
func (k Key) UUID() (string, error) {
	if !k.IsValid() {
		return "", errors.New("invalid UUID key")
	}
	return k.Decode()
}

// config is a struct that contains configuration settings
type config struct {
	hyphens bool
}

// defaultConfig is the default configuration for the uuidkey package
var defaultConfig = config{
	hyphens: true,
}

// Option is a function that configures options
type Option func(c *config)

// apply will apply the options to the default options
func apply(opts ...Option) config {
	c := defaultConfig
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// WithoutHyphens expects no hyphens in the Key
var WithoutHyphens Option = func(c *config) {
	c.hyphens = false
}

// encode will convert your given int64 into base32 crockford encoding format
func encode(n uint64) string {
	encoded := crock32.Encode(n)
	padding := 7 - len(encoded)
	return strings.ToUpper((strings.Repeat("0", padding) + encoded))
}

// decode will convert your given string into original UUID part string
func decode(s string) string {
	i, _ := crock32.Decode(s)
	decoded := strconv.FormatUint(i, 16)
	padding := 8 - len(decoded)
	return (strings.Repeat("0", padding) + decoded)
}

// Encode will encode a given UUID string into a Key.
// It pre-allocates the exact string capacity needed for better performance.
func Encode(uuid string, opts ...Option) (Key, error) {
	options := apply(opts...)

	if len(uuid) != UUIDLength {
		return "", fmt.Errorf("invalid UUID length: expected %d characters, got %d", UUIDLength, len(uuid))
	}

	// Pre-allocate a single slice for the result
	result := make([]byte, KeyLengthWithHyphens)

	// Process parts directly into the result slice
	// UUID format: 8-4-4-4-12 = xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	processUUIDPart(uuid[0:8], result[0:7])                 // First 8 chars
	processUUIDPart(uuid[9:13]+uuid[14:18], result[8:15])   // 4-4 middle section
	processUUIDPart(uuid[19:23]+uuid[24:28], result[16:23]) // 4-4 middle section
	processUUIDPart(uuid[28:36], result[24:31])             // Last 12 chars

	// Add hyphens if needed
	if options.hyphens {
		result[7] = '-'
		result[15] = '-'
		result[23] = '-'
		return Key(result), nil
	}

	// Return without hyphens
	return Key(append(result[:7], append(result[8:15], append(result[16:23], result[24:31]...)...)...)), nil
}

// processUUIDPart converts a hex UUID part directly to base32 and writes to the destination
func processUUIDPart(src string, dst []byte) {
	n, _ := strconv.ParseUint(src, 16, 64)
	encoded := crock32.Encode(n)
	padding := 7 - len(encoded)

	// Write padding zeros
	for i := 0; i < padding; i++ {
		dst[i] = '0'
	}

	// Write encoded part
	copy(dst[padding:], strings.ToUpper(encoded))
}

// EncodeBytes encodes a [16]byte UUID into a Key.
func EncodeBytes(uuid [16]byte, opts ...Option) (Key, error) {
	// apply the options to the default options
	options := apply(opts...)

	// Convert byte groups directly to uint64
	// Each group of 4 bytes is combined into a single uint64
	n1 := uint64(uuid[0])<<24 | uint64(uuid[1])<<16 | uint64(uuid[2])<<8 | uint64(uuid[3])
	n2 := uint64(uuid[4])<<24 | uint64(uuid[5])<<16 | uint64(uuid[6])<<8 | uint64(uuid[7])
	n3 := uint64(uuid[8])<<24 | uint64(uuid[9])<<16 | uint64(uuid[10])<<8 | uint64(uuid[11])
	n4 := uint64(uuid[12])<<24 | uint64(uuid[13])<<16 | uint64(uuid[14])<<8 | uint64(uuid[15])

	// Encode each uint64 into base32 crockford encoding format
	e1 := encode(n1) // Encodes bytes 0-3
	e2 := encode(n2) // Encodes bytes 4-7
	e3 := encode(n3) // Encodes bytes 8-11
	e4 := encode(n4) // Encodes bytes 12-15

	// Build and return key
	if options.hyphens {
		return Key(e1 + "-" + e2 + "-" + e3 + "-" + e4), nil
	}
	return Key(e1 + e2 + e3 + e4), nil
}

// Decode will decode a given Key into a UUID string with basic length validation.
func (k Key) Decode() (string, error) {
	// determine if we should expect hyphens given the length of the key
	hyphens := false
	length := len(k)
	if length != KeyLengthWithoutHyphens {
		if length != KeyLengthWithHyphens {
			return "", fmt.Errorf("invalid Key length: expected %d or %d characters, got %d", KeyLengthWithoutHyphens, KeyLengthWithHyphens, length)
		}
		hyphens = true
	}

	// convert the type from a Key to string
	key := string(k)
	var s1, s2, s3, s4 string

	// select the 4 parts of the key string
	if hyphens {
		s1 = key[0:7]   // [38QARV0]-1ET0G6Z-2CJD9VA-2ZZAR0X
		s2 = key[8:15]  // 38QARV0-[1ET0G6Z]-2CJD9VA-2ZZAR0X
		s3 = key[16:23] // 38QARV0-1ET0G6Z-[2CJD9VA]-2ZZAR0X
		s4 = key[24:31] // 38QARV0-1ET0G6Z-2CJD9VA-[2ZZAR0X]
	} else {
		s1 = key[0:7]   // [38QARV0]1ET0G6Z2CJD9VA2ZZAR0X
		s2 = key[7:14]  // 38QARV0[1ET0G6Z]2CJD9VA2ZZAR0X
		s3 = key[14:21] // 38QARV01ET0G6Z[2CJD9VA]2ZZAR0X
		s4 = key[21:28] // 38QARV01ET0G6Z2CJD9VA[2ZZAR0X]
	}

	// decode each string part into original UUID part string
	n1 := decode(s1)
	n2 := decode(s2)
	n3 := decode(s3)
	n4 := decode(s4)

	// select the 4 parts of the decoded parts
	n2a := n2[0:4]
	n2b := n2[4:8]
	n3a := n3[0:4]
	n3b := n3[4:8]

	var builder strings.Builder
	builder.Grow(UUIDLength) // Pre-allocate exact size needed

	// Write parts with proper formatting
	builder.WriteString(n1)
	builder.WriteByte('-')
	builder.WriteString(n2a)
	builder.WriteByte('-')
	builder.WriteString(n2b)
	builder.WriteByte('-')
	builder.WriteString(n3a)
	builder.WriteByte('-')
	builder.WriteString(n3b)
	builder.WriteString(n4)

	return builder.String(), nil
}

// Bytes converts a Key to a [16]byte UUID.
func (k Key) Bytes() ([16]byte, error) {
	// convert the type from a Key to string
	keyStr := string(k)

	// determine if we should expect hyphens given the length of the key
	hyphens := false
	length := len(keyStr)
	if length != KeyLengthWithoutHyphens {
		if length != KeyLengthWithHyphens {
			return [16]byte{}, fmt.Errorf("invalid Key length: expected %d or %d characters, got %d", KeyLengthWithoutHyphens, KeyLengthWithHyphens, length)
		}
		hyphens = true
	}

	// initialize the UUID array
	var uuid [16]byte
	var err error
	var n uint64

	// select the 4 parts of the key string
	keyParts := [4]string{keyStr[0:7], keyStr[7:14], keyStr[14:21], keyStr[21:28]}
	if hyphens {
		keyParts = [4]string{keyStr[0:7], keyStr[8:15], keyStr[16:23], keyStr[24:31]}
	}

	// Process each part of the key
	for i, part := range keyParts {
		if n, err = crock32.Decode(strings.ToLower(part)); err != nil {
			return [16]byte{}, fmt.Errorf("failed to decode Key part: %v", err)
		}

		// Write 4 bytes for each part
		uuid[i*4] = byte(n >> 24)
		uuid[i*4+1] = byte(n >> 16)
		uuid[i*4+2] = byte(n >> 8)
		uuid[i*4+3] = byte(n)
	}

	return uuid, nil
}
