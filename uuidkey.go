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

// Key is a UUID Key string
type Key string

// String will convert your Key into a string
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

// Valid will verify a given Key follows the correct format:
// e.g. 38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X
// 1. 31 characters long
// 2. Uppercase
// 3. Contains only alphanumeric characters
// 4. Contains 3 hyphens
// 5. Each part is 7 characters long
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

// UUIDString will validate and convert a given Key into a UUID string
func (k Key) UUIDString() (string, error) {
	if !k.Valid() {
		return "", errors.New("invalid UUID key")
	}
	return k.Decode(), nil
}
