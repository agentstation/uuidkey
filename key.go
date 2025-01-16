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

// config is a struct that contains configuration settings
type config struct {
	hyphens     bool
	entropySize numOfCrock32Chars
}

// defaultConfig is the default configuration for the uuidkey package
var defaultConfig = config{
	hyphens:     true,
	entropySize: EntropyBits160,
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

// decode will convert your given string into original UUID part string
func decode(s string) string {
	i, _ := crock32.Decode(s)

	var builder strings.Builder
	builder.Grow(8) // We know the result will be 8 chars

	decoded := strconv.FormatUint(i, 16)
	padding := 8 - len(decoded)

	for i := 0; i < padding; i++ {
		builder.WriteByte('0')
	}
	builder.WriteString(decoded)

	return builder.String()
}

// Encode will encode a given UUID string into a Key.
// It pre-allocates the exact string capacity needed for better performance.
func Encode(uuid string, opts ...Option) (Key, error) {
	options := apply(opts...)

	if len(uuid) != UUIDLength {
		return "", fmt.Errorf("invalid UUID length: expected %d characters, got %d", UUIDLength, len(uuid))
	}

	// Pre-allocate a strings.Builder with exact capacity
	var builder strings.Builder
	if options.hyphens {
		builder.Grow(KeyLengthWithHyphens)
	} else {
		builder.Grow(KeyLengthWithoutHyphens)
	}

	// Process each part
	processAndWritePart(&builder, uuid[0:8])
	if options.hyphens {
		builder.WriteByte('-')
	}
	processAndWritePart(&builder, uuid[9:13]+uuid[14:18])
	if options.hyphens {
		builder.WriteByte('-')
	}
	processAndWritePart(&builder, uuid[19:23]+uuid[24:28])
	if options.hyphens {
		builder.WriteByte('-')
	}
	processAndWritePart(&builder, uuid[28:36])

	return Key(builder.String()), nil
}

func processAndWritePart(builder *strings.Builder, src string) {
	n, _ := strconv.ParseUint(src, 16, 64)
	encoded := crock32.Encode(n)
	padding := 7 - len(encoded)

	// Write padding zeros
	for i := 0; i < padding; i++ {
		builder.WriteByte('0')
	}

	// Write encoded part in uppercase
	for i := 0; i < len(encoded); i++ {
		builder.WriteByte(toUpper(encoded[i]))
	}
}

// toUpper converts a single byte to uppercase if it's a lowercase letter
func toUpper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 32
	}
	return c
}

// EncodeBytes encodes a [16]byte UUID into a Key.
func EncodeBytes(uuid [16]byte, opts ...Option) (Key, error) {
	options := apply(opts...)

	var builder strings.Builder
	if options.hyphens {
		builder.Grow(KeyLengthWithHyphens)
	} else {
		builder.Grow(KeyLengthWithoutHyphens)
	}

	// Process each 4-byte group
	writeEncodedPart(&builder, uint64(uuid[0])<<24|uint64(uuid[1])<<16|uint64(uuid[2])<<8|uint64(uuid[3]))
	if options.hyphens {
		builder.WriteByte('-')
	}
	writeEncodedPart(&builder, uint64(uuid[4])<<24|uint64(uuid[5])<<16|uint64(uuid[6])<<8|uint64(uuid[7]))
	if options.hyphens {
		builder.WriteByte('-')
	}
	writeEncodedPart(&builder, uint64(uuid[8])<<24|uint64(uuid[9])<<16|uint64(uuid[10])<<8|uint64(uuid[11]))
	if options.hyphens {
		builder.WriteByte('-')
	}
	writeEncodedPart(&builder, uint64(uuid[12])<<24|uint64(uuid[13])<<16|uint64(uuid[14])<<8|uint64(uuid[15]))

	return Key(builder.String()), nil
}

func writeEncodedPart(builder *strings.Builder, n uint64) {
	encoded := crock32.Encode(n)
	padding := 7 - len(encoded)

	// Write padding zeros
	for i := 0; i < padding; i++ {
		builder.WriteByte('0')
	}

	// Write encoded part in uppercase
	for i := 0; i < len(encoded); i++ {
		builder.WriteByte(toUpper(encoded[i]))
	}
}

// Decode will decode a given Key into a UUID string with basic length validation.
func (k Key) Decode() (string, error) {
	// determine if we should expect hyphens given the length of the key
	hyphens := false
	length := len(k)
	if length != KeyLengthWithoutHyphens {
		if length != KeyLengthWithHyphens {
			return "", fmt.Errorf("invalid Key length: expected %d or %d characters, got %d",
				KeyLengthWithoutHyphens, KeyLengthWithHyphens, length)
		}
		hyphens = true
	}

	var builder strings.Builder
	builder.Grow(UUIDLength) // Pre-allocate exact size needed: 36 bytes

	// select the 4 parts of the key string
	var s1, s2, s3, s4 string
	s := string(k)
	if hyphens {
		s1 = s[0:7]   // [38QARV0]-1ET0G6Z-2CJD9VA-2ZZAR0X
		s2 = s[8:15]  // 38QARV0-[1ET0G6Z]-2CJD9VA-2ZZAR0X
		s3 = s[16:23] // 38QARV0-1ET0G6Z-[2CJD9VA]-2ZZAR0X
		s4 = s[24:31] // 38QARV0-1ET0G6Z-2CJD9VA-[2ZZAR0X]
	} else {
		s1 = s[0:7]   // [38QARV0]1ET0G6Z2CJD9VA2ZZAR0X
		s2 = s[7:14]  // 38QARV0[1ET0G6Z]2CJD9VA2ZZAR0X
		s3 = s[14:21] // 38QARV01ET0G6Z[2CJD9VA]2ZZAR0X
		s4 = s[21:28] // 38QARV01ET0G6Z2CJD9VA[2ZZAR0X]
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
	length := len(k)
	if length != KeyLengthWithoutHyphens && length != KeyLengthWithHyphens {
		return [16]byte{}, fmt.Errorf("invalid Key length: expected %d or %d characters, got %d",
			KeyLengthWithoutHyphens, KeyLengthWithHyphens, length)
	}

	hyphens := length == KeyLengthWithHyphens
	var uuid [16]byte
	var err error

	// Avoid string conversion by working directly with the Key type
	s := string(k)
	if hyphens {
		if err = processByteGroup(s[0:7], &uuid, 0); err != nil {
			return [16]byte{}, err
		}
		if err = processByteGroup(s[8:15], &uuid, 4); err != nil {
			return [16]byte{}, err
		}
		if err = processByteGroup(s[16:23], &uuid, 8); err != nil {
			return [16]byte{}, err
		}
		if err = processByteGroup(s[24:31], &uuid, 12); err != nil {
			return [16]byte{}, err
		}
	} else {
		if err = processByteGroup(s[0:7], &uuid, 0); err != nil {
			return [16]byte{}, err
		}
		if err = processByteGroup(s[7:14], &uuid, 4); err != nil {
			return [16]byte{}, err
		}
		if err = processByteGroup(s[14:21], &uuid, 8); err != nil {
			return [16]byte{}, err
		}
		if err = processByteGroup(s[21:28], &uuid, 12); err != nil {
			return [16]byte{}, err
		}
	}

	return uuid, nil
}

func processByteGroup(part string, uuid *[16]byte, offset int) error {
	n, err := crock32.Decode(strings.ToLower(part))
	if err != nil {
		return fmt.Errorf("failed to decode Key part: %v", err)
	}

	uuid[offset] = byte(n >> 24)
	uuid[offset+1] = byte(n >> 16)
	uuid[offset+2] = byte(n >> 8)
	uuid[offset+3] = byte(n)
	return nil
}
