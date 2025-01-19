package uuidkey

import (
	"crypto/rand"
	"fmt"
	"hash"
	"strings"

	"hash/crc32"

	"github.com/richardlehane/crock32"
	"golang.org/x/crypto/blake2b"
)

// Secret Scanning Configuration for GitHub
//
// To detect these API keys in GitHub's secret scanning program, use the following regex pattern:
// [A-Z]{8}[A-Z0-9]?_[A-Z0-9]{26}[A-Z0-9]{14,42}_[A-F0-9]{8}
//
// Pattern breakdown:
// - [A-Z]{8}[A-Z0-9]?    : Prefix (8 uppercase letters, optional alphanumeric)
// - _                      : First separator
// - [A-Z0-9]{26}          : UUID Key (26 chars in Base32-Crockford)
// - [A-Z0-9]{14,42}       : Entropy (14-42 chars depending on entropy bits)
// - _                      : Second separator
// - [A-F0-9]{8}           : CRC32 checksum (8 hex chars)
//
// Example match (128 bits of entropy):
// AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJBJLSO7WBNWY3F_A1B2C3D8
//
// Example match (160 bits of entropy):
// AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJBJLSO7WBNWY3FK7EK2JP_A1B2C3D8
//
// Example match (256 bits of entropy):
// AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJBJLSO7WBNWY3FK7EK2JPXQW23MGHE7X33BNA8R5ET_A1B2C3D8

var (
	// blake2bHasher is a package-level hasher for entropy generation
	blake2bHasher hash.Hash
	// initErr captures any initialization error
	initErr error
)

func init() {
	blake2bHasher, initErr = blake2b.New256(nil)
}

// checkInit verifies the hasher was properly initialized
func checkInit() error {
	if initErr != nil {
		return fmt.Errorf("BLAKE2b hasher initialization failed: %w", initErr)
	}
	if blake2bHasher == nil {
		return fmt.Errorf("BLAKE2b hasher not initialized")
	}
	return nil
}

// numOfCrock32Chars represents the length of characters needed in an
// APIKey entropy segment to provide a given number of bits of entropy.
type numOfCrock32Chars int

const (
	// EntropyBits128 represents the length of characters needed in
	// an APIKey entropy segment to provide 128 bits of entropy.
	// It assumes use of a UUIDv7 Base32-Crockford encoded Key
	// and that the entropy is also encoded using Base32-Crockford.
	EntropyBits128 numOfCrock32Chars = 14

	// EntropyBits160 represents the length of characters needed in
	// an APIKey entropy segment to provide 160 bits of entropy.
	// It assumes use of a UUIDv7 Base32-Crockford encoded Key
	// and that the entropy is also encoded using Base32-Crockford.
	EntropyBits160 numOfCrock32Chars = 21

	// EntropyBits256 represents the length of characters needed in
	// an APIKey entropy segment to provide 256 bits of entropy.
	// It assumes use of a UUIDv7 Base32-Crockford encoded Key
	// and that the entropy is also encoded using Base32-Crockford.
	EntropyBits256 numOfCrock32Chars = 42

	checksumLength         = 8
	entropyBytesMultiplier = 8 / 5 // Used for calculating entropy bytes
	initialEntropyBytes    = 32    // Initial entropy buffer size
)

// With128BitEntropy expects 128 bits of entropy in the APIKey
var With128BitEntropy Option = func(c *config) {
	c.entropySize = EntropyBits128
}

// With160BitEntropy expects 160 bits of entropy in the APIKey
var With160BitEntropy Option = func(c *config) {
	c.entropySize = EntropyBits160
}

// With256BitEntropy expects 256 bits of entropy in the APIKey
var With256BitEntropy Option = func(c *config) {
	c.entropySize = EntropyBits256
}

// APIKey represents a compound key consisting of four parts or segments:
// - Prefix: A company or application identifier (e.g., "AGNTSTNP")
// - Key: A UUID-based identifier encoded in Base32-Crockford
// - Entropy: Additional segment of random data for increased uniqueness
// - Checksum: CRC32 checksum of the previous components (8 characters)
//
// Format:
//
//	[Prefix]_[UUID Key][Entropy]_[Checksum]
//
//	AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJJLSO7WBNWY3F_A1B2C3D8
//	└─────┘ └──────────────────────────┘└────────────┘ └──────┘
//	Prefix     UUID-based identifier       Entropy      Checksum
type APIKey struct {
	Prefix   string
	Key      Key
	Entropy  string
	Checksum string
}

// calculateChecksum generates an 8-character hexadecimal CRC32 checksum
func (a APIKey) calculateChecksum() string {
	// Combine all parts except checksum
	data := a.Prefix + "_" + a.Key.String() + a.Entropy

	// Calculate CRC32 checksum and return as uppercase hex
	crc := crc32.ChecksumIEEE([]byte(data))
	return fmt.Sprintf("%08X", crc) // Always 8 chars, zero-padded, uppercase
}

// String returns the complete API key as a string with all components joined
func (a APIKey) String() string {
	if a.Checksum == "" {
		a.Checksum = a.calculateChecksum()
	}
	// Pre-allocate capacity based on known component lengths
	var sb strings.Builder
	sb.Grow(len(a.Prefix) + 1 + len(a.Key.String()) + len(a.Entropy) + 1 + len(a.Checksum))
	sb.WriteString(a.Prefix)
	sb.WriteByte('_')
	sb.WriteString(a.Key.String())
	sb.WriteString(a.Entropy)
	sb.WriteByte('_')
	sb.WriteString(a.Checksum)
	return sb.String()
}

// NewAPIKey creates a new APIKey from a string prefix, string UUID, and options.
func NewAPIKey(prefix, uuid string, opts ...Option) (APIKey, error) {
	// Check if the prefix is empty
	if prefix == "" {
		return APIKey{}, fmt.Errorf("prefix cannot be empty")
	}

	// Apply options
	options := apply(opts...)

	// Encode the UUID without hyphens
	key, err := Encode(uuid, WithoutHyphens)
	if err != nil {
		return APIKey{}, err
	}

	// Generate entropy of the specified size
	entropy, err := generateEntropy(options.entropySize)
	if err != nil {
		return APIKey{}, err
	}

	// Create the APIKey with the prefix, key, and entropy
	apiKey := APIKey{
		Prefix:  prefix,
		Key:     key,
		Entropy: entropy,
	}

	// Calculate the checksum for the APIKey
	apiKey.Checksum = apiKey.calculateChecksum()

	return apiKey, nil
}

// NewAPIKeyFromBytes creates a new APIKey from a string prefix, [16]byte UUID, and options.
func NewAPIKeyFromBytes(prefix string, uuid [16]byte, opts ...Option) (APIKey, error) {
	// Apply options
	options := apply(opts...)

	// Encode the UUID without hyphens
	key, err := EncodeBytes(uuid, WithoutHyphens)
	if err != nil {
		return APIKey{}, err
	}

	// Generate entropy of the specified size
	entropy, err := generateEntropy(options.entropySize)
	if err != nil {
		return APIKey{}, err
	}

	// Create the APIKey with the prefix, key, and entropy
	apiKey := APIKey{
		Prefix:  prefix,
		Key:     key,
		Entropy: entropy,
	}

	// Calculate the checksum for the APIKey
	apiKey.Checksum = apiKey.calculateChecksum()

	// Return the fully constructed APIKey
	return apiKey, nil
}

// ParseAPIKey will parse a given APIKey string into an APIKey type.
func ParseAPIKey(apikey string) (APIKey, error) {
	// Check if the APIKey is empty
	if apikey == "" {
		return APIKey{}, fmt.Errorf("invalid APIKey format: expected 3 parts, got 1")
	}

	// Split the string to count parts (for error messages)
	parts := strings.Split(apikey, "_")
	if len(parts) != 3 {
		return APIKey{}, fmt.Errorf("invalid APIKey format: expected 3 parts, got %d", len(parts))
	}

	// Now use the more efficient index-based approach
	firstSep := strings.IndexByte(apikey, '_')
	lastSep := strings.LastIndexByte(apikey, '_')

	// Extract the prefix
	prefix := apikey[:firstSep]
	if prefix == "" {
		return APIKey{}, fmt.Errorf("invalid prefix: cannot be empty")
	}

	// Extract the remainder (Key and Entropy)
	remainder := apikey[firstSep+1 : lastSep]

	// Extract the checksum
	checksum := apikey[lastSep+1:]

	// The remainder should contain both the Key and Entropy parts
	if len(remainder) < KeyLengthWithoutHyphens {
		return APIKey{}, fmt.Errorf("invalid Key format: insufficient length")
	}

	// Extract the Key part (first KeyLengthWithoutHyphens characters)
	keyPart := remainder[:KeyLengthWithoutHyphens]

	// Parse the Key part
	key, err := Parse(keyPart)
	if err != nil {
		return APIKey{}, fmt.Errorf("invalid Key format: %v", err)
	}

	// The rest is entropy
	entropy := remainder[KeyLengthWithoutHyphens:]

	// Validate checksum format (must be 8 uppercase hexadecimal characters)
	if len(checksum) != checksumLength {
		return APIKey{}, fmt.Errorf("invalid checksum format: must be 8 hexadecimal characters")
	}
	for _, c := range checksum {
		// Check if character is 0-9 or A-F
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
			return APIKey{}, fmt.Errorf("invalid checksum format: must be 8 hexadecimal characters")
		}
	}

	// Create the APIKey with the prefix, key, entropy, and checksum
	apiKey := APIKey{
		Prefix:   prefix,
		Key:      key,
		Entropy:  entropy,
		Checksum: checksum,
	}

	// Verify checksum
	expectedChecksum := apiKey.calculateChecksum()
	if checksum != expectedChecksum {
		return APIKey{}, fmt.Errorf("invalid checksum: expected %s, got %s", expectedChecksum, checksum)
	}

	return apiKey, nil
}

// generateEntropy generates entropy of a given size
func generateEntropy(size numOfCrock32Chars) (string, error) {
	if err := checkInit(); err != nil {
		return "", err
	}

	numOfRandomBytes := int(size) * entropyBytesMultiplier

	// Generate initial random bytes
	inputBytes := make([]byte, initialEntropyBytes)
	if _, err := rand.Read(inputBytes); err != nil {
		return "", err
	}

	// Pre-allocate entropy buffer
	entropy := make([]byte, numOfRandomBytes)
	for pos := 0; pos < numOfRandomBytes; pos += blake2bHasher.Size() {
		blake2bHasher.Reset()
		blake2bHasher.Write(inputBytes)
		hash := blake2bHasher.Sum(nil)

		copyLen := numOfRandomBytes - pos
		if copyLen > len(hash) {
			copyLen = len(hash)
		}
		copy(entropy[pos:], hash[:copyLen])
		inputBytes = hash
	}

	// Process bytes more efficiently
	var entropyEncoded strings.Builder
	entropyEncoded.Grow(int(size))

	for i := 0; i < len(entropy); i += 8 {
		end := i + 8
		if end > len(entropy) {
			end = len(entropy)
		}

		var n uint64
		for j, b := range entropy[i:end] {
			n |= uint64(b) << (8 * (end - i - 1 - j))
		}

		entropyEncoded.WriteString(crock32.Encode(n))
	}

	result := strings.ToUpper(entropyEncoded.String())
	if len(result) < int(size) {
		result = strings.Repeat("0", int(size)-len(result)) + result
	}

	return result[:int(size)], nil
}
