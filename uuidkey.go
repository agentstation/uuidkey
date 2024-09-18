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
	// KeyLength is the total length of a valid UUID Key, including hyphens.
	KeyLength = 31 // 7 + 1 + 7 + 1 + 7 + 1 + 7 = 31 characters

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
	if !k.Valid() {
		return "", errors.New("invalid UUID Key")
	}
	return k, nil
}

// Valid verifies if a given Key follows the correct format.
// The format should be:
//   - 31 characters long
//   - Uppercase
//   - Contains only alphanumeric characters
//   - Contains 3 hyphens
//   - Each part is 7 characters long
//
// Examples of valid keys:
//   - 38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X
//   - ABCDEFG-1234567-HIJKLMN-OPQRSTU
//
// Examples of invalid keys:
//   - 38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0  (too short)
//   - 38qarv0-1ET0G6Z-2CJD9VA-2ZZAR0X (contains lowercase)
//   - 38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X- (extra hyphen)
//   - 38QARV0-1ET0G6Z-2CJD9VA2ZZAR0X (missing hyphen)
//   - 38QARV0-1ET0G6-2CJD9VA-2ZZAR0X (part too short)
func (k Key) Valid() bool {
	if len(k) != KeyLength { // check if the key is 31 characters long
		return false
	}
	hyphenCount := 0
	partLen := 0
	for _, char := range k {
		switch {
		case char == '-':
			hyphenCount++                 // collect the number of hyphens
			if partLen != KeyPartLength { // check parts are 7 characters long
				return false
			}
			partLen = 0 // reset the part length
		// check if the key is uppercase
		// check if the key contains only alphanumeric characters
		case char < '0' || (char > '9' && char < 'A') || char > 'Z':
			return false
		default:
			partLen++
		}
	}
	// check if the key contains 3 hyphens and the last part is 7 characters long
	return hyphenCount == KeyHyphenCount && partLen == KeyPartLength
}

// UUID will validate and convert a given Key into a UUID string.
func (k Key) UUID() (string, error) {
	if !k.Valid() {
		return "", errors.New("invalid UUID key")
	}
	return k.Decode()
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

// Encode will encode a given UUID string into a Key with basic length validation.
func Encode(uuid string) (Key, error) {
	if len(uuid) != UUIDLength { // basic length validation to ensure we can encode
		return "", fmt.Errorf("invalid UUID length: expected %d characters, got %d", UUIDLength, len(uuid))
	}

	// select the 5 parts of the UUID string
	s1 := uuid[0:8]   // [d1756360]-5da0-40df-9926-a76abff5601d
	s2 := uuid[9:13]  // d1756360-[5da0]-40df-9926-a76abff5601d
	s3 := uuid[14:18] // d1756360-5da0-[40df]-9926-a76abff5601d
	s4 := uuid[19:23] // d1756360-5da0-40df-[9926]-a76abff5601d
	s5 := uuid[24:36] // d1756360-5da0-40df-9926-[a76abff5601d]

	// decode each string part into uint64
	n1, _ := strconv.ParseUint(s1, 16, 32)
	n2, _ := strconv.ParseUint(s2+s3, 16, 32)     // combine s2 and s3
	n3, _ := strconv.ParseUint(s4+s5[:4], 16, 32) // combine s4 and the first 4 chars of s5
	n4, _ := strconv.ParseUint(s5[4:], 16, 32)    // the last 8 chars of s5

	// encode each uint64 into base32 crockford encoding format
	e1 := encode(n1)
	e2 := encode(n2)
	e3 := encode(n3)
	e4 := encode(n4)

	// build and return key
	return Key(e1 + "-" + e2 + "-" + e3 + "-" + e4), nil
}

// EncodeBytes encodes a [16]byte UUID into a Key.
func EncodeBytes(uuid [16]byte) (Key, error) {
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
	// The key is constructed by joining the encoded parts with hyphens
	return Key(e1 + "-" + e2 + "-" + e3 + "-" + e4), nil
}

// Decode will decode a given Key into a UUID string with basic length validation.
func (k Key) Decode() (string, error) {
	if len(k) != KeyLength { // basic length validation to ensure we can decode
		return "", fmt.Errorf("invalid Key length: expected %d characters, got %d", KeyLength, len(k))
	}

	// select the 4 parts of the key string
	key := string(k) // convert the type from a Key to string
	s1 := key[0:7]   // [38QARV0]-1ET0G6Z-2CJD9VA-2ZZAR0X
	s2 := key[8:15]  // 38QARV0-[1ET0G6Z]-2CJD9VA-2ZZAR0X
	s3 := key[16:23] // 38QARV0-1ET0G6Z-[2CJD9VA]-2ZZAR0X
	s4 := key[24:31] // 38QARV0-1ET0G6Z-2CJD9VA-[2ZZAR0X]

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

	// build and return UUID string
	return (n1 + "-" + n2a + "-" + n2b + "-" + n3a + "-" + n3b + n4), nil
}

// Bytes converts a Key to a [16]byte UUID.
func (k Key) Bytes() ([16]byte, error) {
	keyStr := string(k)

	// Check the length of the Key
	if len(keyStr) != KeyLength {
		return [16]byte{}, fmt.Errorf("invalid Key length: expected %d characters, got %d", KeyLength, len(keyStr))
	}

	var uuid [16]byte
	var err error
	var n uint64

	// Process each part of the key
	for i, part := range [4]string{keyStr[:7], keyStr[8:15], keyStr[16:23], keyStr[24:]} {
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
