// Package uuidkey provides Crockford Base32 encoding that's compatible with the original crock32 number-based approach
package uuidkey

import (
	"encoding/base32"
	"fmt"
)

// crockford is the standard library base32 encoding with Crockford's alphabet (for entropy generation)
var crockford = base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base32.NoPadding)

// crock32Encode encodes a uint32 as a Crockford Base32 string (number-based encoding)
func crock32Encode(n uint32) string {
	const digits = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	if n == 0 {
		return "0"
	}
	
	// Use a fixed-size array to avoid allocations
	var buf [8]byte // Maximum size needed for uint32 in base32
	idx := len(buf)
	
	for n > 0 {
		idx--
		buf[idx] = digits[n%32]
		n /= 32
	}
	
	return string(buf[idx:])
}

// Pre-computed lookup table for fast character to digit conversion
// 255 indicates invalid character
var decodeTable = func() [256]byte {
	table := [256]byte{}
	// Initialize all values to 255 (invalid)
	for i := range table {
		table[i] = 255
	}
	
	// Numbers 0-9
	for i := byte('0'); i <= '9'; i++ {
		table[i] = i - '0'
	}
	
	// Uppercase letters
	table['A'] = 10
	table['B'] = 11
	table['C'] = 12
	table['D'] = 13
	table['E'] = 14
	table['F'] = 15
	table['G'] = 16
	table['H'] = 17
	// Skip I
	table['J'] = 18
	table['K'] = 19
	// Skip L
	table['M'] = 20
	table['N'] = 21
	// Skip O
	table['P'] = 22
	table['Q'] = 23
	table['R'] = 24
	table['S'] = 25
	table['T'] = 26
	// Skip U
	table['V'] = 27
	table['W'] = 28
	table['X'] = 29
	table['Y'] = 30
	table['Z'] = 31
	
	// Lowercase letters (same values as uppercase)
	for c := byte('a'); c <= 'z'; c++ {
		if table[c-'a'+'A'] != 255 {
			table[c] = table[c-'a'+'A']
		}
	}
	
	// Special mappings per Crockford spec
	table['O'] = 0
	table['o'] = 0
	table['I'] = 1
	table['i'] = 1
	table['L'] = 1
	table['l'] = 1
	
	return table
}()

// crock32Decode decodes a Crockford Base32 string to a uint32 (number-based decoding)
func crock32Decode(s string) (uint32, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("crock32.Decode: empty string")
	}
	
	var result uint32
	for i := 0; i < len(s); i++ {
		// Use lookup table for fast conversion
		digit := decodeTable[s[i]]
		if digit == 255 {
			return 0, fmt.Errorf("crock32.Decode: invalid character %c", s[i])
		}
		
		// Check for overflow before multiplication
		if result > (^uint32(0))/32 {
			return 0, fmt.Errorf("crock32.Decode: integer overflow")
		}
		result = result*32 + uint32(digit)
	}
	
	return result, nil
}