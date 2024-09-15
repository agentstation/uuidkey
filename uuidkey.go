// Package uuidkey encodes UUIDs to a readable Key format via the Base32-Crockford codec.
package uuidkey

import (
	"errors"
)

// key validation constraint constants
const (
	key_len      = 31
	key_part_len = 7
	key_hyphens  = 3
)

// Key is a UUID Key string.
type Key string

// String will convert your Key into a string.
func (k Key) String() string {
	return string(k)
}

// FromString will convert a Key formatted string type into a Key type.
func FromString(key string) (Key, error) {
	k := Key(key)
	if !k.Valid() {
		return "", errors.New("invalid UUID Key")
	}
	return Key(key), nil
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
	if len(k) != key_len { // check if the key is 31 characters long
		return false
	}
	hyphenCount := 0
	partLen := 0
	for _, char := range k {
		switch {
		case char == '-':
			hyphenCount++                // collect the number of hyphens
			if partLen != key_part_len { // check parts are 7 characters long
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
	return hyphenCount == key_hyphens && partLen == key_part_len
}

// UUIDString will validate and convert a given Key into a UUID string.
func (k Key) UUIDString() (string, error) {
	if !k.Valid() {
		return "", errors.New("invalid UUID key")
	}
	return k.Decode(), nil
}
