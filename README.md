```sh
                         _   _  _   _  ___ ____     _  __          
                        | | | || | | ||_ _|  _ \   | |/ /___ _   _ 
                        | | | || | | | | || | | |  | ' // _ \ | | |
                        | |_| || |_| | | || |_| |  | . \  __/ |_| |
                         \___/  \___/ |___|____/   |_|\_\___|\__, |
                                                             |___/ 
```
<!-- [![Sourcegraph](https://sourcegraph.com/github.com/agentstation/uuidkey/-/badge.svg?style=flat-square)](https://sourcegraph.com/github.com/agentstation/uuidkey?badge) -->
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/agentstation/uuidkey)
[![Go Report Card](https://goreportcard.com/badge/github.com/agentstation/uuidkey?style=flat-square)](https://goreportcard.com/report/github.com/agentstation/uuidkey)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/agentstation/uuidkey/ci.yaml?style=flat-square)](https://github.com/agentstation/uuidkey/actions)
[![codecov](https://codecov.io/gh/agentstation/uuidkey/branch/master/graph/badge.svg?token=35UM5QX1Q3)](https://codecov.io/gh/agentstation/uuidkey)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/agentstation/uuidkey/master/LICENSE)
<!-- [![Forum](https://img.shields.io/badge/community-forum-00afd1.svg?style=flat-square)](https://github.com/agentstation/uuidkey/discussions) -->
<!-- [![Twitter](https://img.shields.io/badge/twitter-@agentstationHQ-55acee.svg?style=flat-square)](https://twitter.com/agentstationHQ) -->

The `uuidkey` package encodes UUIDs to a readable `Key` format via the Base32-Crockford codec.

<div align="center"><h3><a href="https://docs.agentstation.ai/blog/beautiful-api-keys/">📚 Read the original article on why we made this!</a></h3></div>

> **Note:** Thanks to everyone for the feedback and suggestions from the original article, we learned a lot and made improvements to follow the GitHub Secret Scanning format (with checksum) and added additional entropy options to ensure UUIDv7 encoding can be used in a wide variety of use cases. You can still use the `Encode` function to generate a `Key` type without the GitHub prefix format and additional entropy - but we recommend using the `NewAPIKey` function with an 8 character prefix to stay symetrical ;P

## Overview

The `uuidkey` package generates secure, readable API keys by encoding UUIDs using Base32-Crockford with additional security features.

You can use the `uuidkey` package to generate API keys for your application using the `NewAPIKey` function (recommended to guarantee at least 128 bits of entropy and follow the GitHub Secret Scanning format) or the `Encode` function (to generate just a `Key` type).

## API Key Format

```
AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJJLSO7WBNWY3F_A1B2C3D8
└─────┘ └──────────────────────────┘└────────────┘ └──────┘
Prefix        Key (crock32 UUID)        Entropy      Checksum
```

### Components
1. **Prefix** - Company/application identifier (e.g., "AGNTSTNP")
2. **UUID Key** - Base32-Crockford encoded UUID
3. **Entropy** - Additional random data (128, 160, or 256 bits)
4. **Checksum** - CRC32 checksum (8 characters) for validation

### Security Features
1. **Secret Scanning** - Formatted for GitHub Secret Scanning detection
2. **Validation** - CRC32 checksum for error detection and validation
3. **Entropy Options** - Configurable entropy levels that ensure UUIDv7 security (128, 160, or 256 bits)

## Compatibility

Compatible with any UUID library following RFC 4122 specification. Officially tested with:
- [github.com/gofrs/uuid](https://github.com/gofrs/uuid) (v4.4.0+)
- [github.com/google/uuid](https://github.com/google/uuid) (v1.6.0+)

## Installation

To install the `uuidkey` package, use the following command:

```sh
go get github.com/agentstation/uuidkey
```

## Usage

To use the `uuidkey` package in your Go code, follow these steps:

1. Import the package:

```go
import "github.com/agentstation/uuidkey"
```

2. Create and parse API Keys:

```go
// Create a new API Key with default settings (160-bit entropy)
apiKey := uuidkey.NewAPIKey("AGNTSTNP", "d1756360-5da0-40df-9926-a76abff5601d")
fmt.Println(apiKey) // Output: AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJJLSO7WBNWY3F_A1B2C3D8

// Create an API Key with 128-bit entropy
apiKey = uuidkey.NewAPIKey("AGNTSTNP", "d1756360-5da0-40df-9926-a76abff5601d", uuidkey.With128BitEntropy)
fmt.Println(apiKey) // Output: AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJJLSO7WBNWY3F_A1B2C3D8

// Parse an existing API Key
apiKey, err := uuidkey.ParseAPIKey("AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJJLSO7WBNWY3F_A1B2C3D8")
if err != nil {
    log.Fatal("Error:", err)
}
fmt.Printf("Prefix: %s, Key: %s, Entropy: %s\n", apiKey.Prefix, apiKey.Key, apiKey.Entropy)
```

3. Work with UUID Keys directly:

```go
// With hyphens (default)
key, _ := uuidkey.Encode("d1756360-5da0-40df-9926-a76abff5601d")
fmt.Println(key) // Output: 38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X

// Without hyphens
key, _ := uuidkey.Encode("d1756360-5da0-40df-9926-a76abff5601d", uuidkey.WithoutHyphens)
fmt.Println(key) // Output: 38QARV01ET0G6Z2CJD9VA2ZZAR0X
```

4. Decode a Key type to a UUID string with Key format validation:

```go
// With hyphens
key, _ := uuidkey.Parse("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
uuid, err := key.UUID()
if err != nil {
    log.Fatal("Error:", err)
}
fmt.Println(uuid) // Output: d1756360-5da0-40df-9926-a76abff5601d

// Without hyphens
key, _ := uuidkey.Parse("38QARV01ET0G6Z2CJD9VA2ZZAR0X")
uuid, err := key.UUID()
if err != nil {
    log.Fatal("Error:", err)
}
fmt.Println(uuid) // Output: d1756360-5da0-40df-9926-a76abff5601d
```

5. Decode a Key type to a UUID string with only basic Key length validation:

```go
key, _ := uuidkey.Parse("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
uuid, err := key.Decode()
if err != nil {
    log.Fatal("Error:", err)
}
fmt.Println(uuid) // Output: d1756360-5da0-40df-9926-a76abff5601d
```

6. Work directly with UUID bytes:

```go
// Convert UUID string to bytes
uuidStr := "d1756360-5da0-40df-9926-a76abff5601d"
uuidBytes, err := hex.DecodeString(strings.ReplaceAll(uuidStr, "-", ""))
if err != nil {
    log.Fatal("Error:", err)
}

// Convert to [16]byte array
var uuid [16]byte
copy(uuid[:], uuidBytes)

// Now use the bytes with EncodeBytes (with hyphens)
key, _ := uuidkey.EncodeBytes(uuid)
fmt.Println(key) // Output: 38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X

// Without hyphens
key, _ := uuidkey.EncodeBytes(uuid, uuidkey.WithoutHyphens)
fmt.Println(key) // Output: 38QARV01ET0G6Z2CJD9VA2ZZAR0X

// Convert Key back to UUID bytes
key, _ := uuidkey.Parse("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
bytes, err := key.Bytes()
if err != nil {
    log.Fatal("Error:", err)
}
fmt.Printf("%x", bytes) // Output: d17563605da040df9926a76abff5601d
```

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# uuidkey

```go
import "github.com/agentstation/uuidkey"
```

Package uuidkey encodes UUIDs to a readable Key format via the Base32\-Crockford codec.

## Index

- [Constants](<#constants>)
- [type APIKey](<#APIKey>)
  - [func NewAPIKey\(prefix, uuid string, opts ...Option\) \(APIKey, error\)](<#NewAPIKey>)
  - [func NewAPIKeyFromBytes\(prefix string, uuid \[16\]byte, opts ...Option\) \(APIKey, error\)](<#NewAPIKeyFromBytes>)
  - [func ParseAPIKey\(apikey string\) \(APIKey, error\)](<#ParseAPIKey>)
  - [func \(a APIKey\) String\(\) string](<#APIKey.String>)
- [type Key](<#Key>)
  - [func Encode\(uuid string, opts ...Option\) \(Key, error\)](<#Encode>)
  - [func EncodeBytes\(uuid \[16\]byte, opts ...Option\) \(Key, error\)](<#EncodeBytes>)
  - [func Parse\(key string\) \(Key, error\)](<#Parse>)
  - [func \(k Key\) Bytes\(\) \(\[16\]byte, error\)](<#Key.Bytes>)
  - [func \(k Key\) Decode\(\) \(string, error\)](<#Key.Decode>)
  - [func \(k Key\) IsValid\(\) bool](<#Key.IsValid>)
  - [func \(k Key\) String\(\) string](<#Key.String>)
  - [func \(k Key\) UUID\(\) \(string, error\)](<#Key.UUID>)
- [type Option](<#Option>)


## Constants

<a name="KeyLengthWithHyphens"></a>Key validation constraint constants

```go
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
```

<a name="APIKey"></a>
## type [APIKey](<https://github.com/agentstation/uuidkey/blob/master/apikey.go#L115-L120>)

APIKey represents a compound key consisting of four parts or segments: \- Prefix: A company or application identifier \(e.g., "AGNTSTNP"\) \- Key: A UUID\-based identifier encoded in Base32\-Crockford \- Entropy: Additional segment of random data for increased uniqueness \- Checksum: CRC32 checksum of the previous components \(8 characters\)

Format:

```
[Prefix]_[UUID Key][Entropy]_[Checksum]

AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJJLSO7WBNWY3F_A1B2C3D8
└─────┘ └──────────────────────────┘└────────────┘ └──────┘
Prefix        Key (crock32 UUID)        Entropy      Checksum
```

```go
type APIKey struct {
    Prefix   string
    Key      Key
    Entropy  string
    Checksum string
}
```

<a name="NewAPIKey"></a>
### func [NewAPIKey](<https://github.com/agentstation/uuidkey/blob/master/apikey.go#L150>)

```go
func NewAPIKey(prefix, uuid string, opts ...Option) (APIKey, error)
```

NewAPIKey creates a new APIKey from a string prefix, string UUID, and options.

<a name="NewAPIKeyFromBytes"></a>
### func [NewAPIKeyFromBytes](<https://github.com/agentstation/uuidkey/blob/master/apikey.go#L185>)

```go
func NewAPIKeyFromBytes(prefix string, uuid [16]byte, opts ...Option) (APIKey, error)
```

NewAPIKeyFromBytes creates a new APIKey from a string prefix, \[16\]byte UUID, and options.

<a name="ParseAPIKey"></a>
### func [ParseAPIKey](<https://github.com/agentstation/uuidkey/blob/master/apikey.go#L216>)

```go
func ParseAPIKey(apikey string) (APIKey, error)
```

ParseAPIKey will parse a given APIKey string into an APIKey type.

<a name="APIKey.String"></a>
### func \(APIKey\) [String](<https://github.com/agentstation/uuidkey/blob/master/apikey.go#L133>)

```go
func (a APIKey) String() string
```

String returns the complete API key as a string with all components joined

<a name="Key"></a>
## type [Key](<https://github.com/agentstation/uuidkey/blob/master/key.go#L37>)

Key is a UUID Key string.

```go
type Key string
```

<a name="Encode"></a>
### func [Encode](<https://github.com/agentstation/uuidkey/blob/master/key.go#L156>)

```go
func Encode(uuid string, opts ...Option) (Key, error)
```

Encode will encode a given UUID string into a Key. It pre\-allocates the exact string capacity needed for better performance.

<a name="EncodeBytes"></a>
### func [EncodeBytes](<https://github.com/agentstation/uuidkey/blob/master/key.go#L214>)

```go
func EncodeBytes(uuid [16]byte, opts ...Option) (Key, error)
```

EncodeBytes encodes a \[16\]byte UUID into a Key.

<a name="Parse"></a>
### func [Parse](<https://github.com/agentstation/uuidkey/blob/master/key.go#L74>)

```go
func Parse(key string) (Key, error)
```

Parse converts a Key formatted string into a Key type.

<a name="Key.Bytes"></a>
### func \(Key\) [Bytes](<https://github.com/agentstation/uuidkey/blob/master/key.go#L316>)

```go
func (k Key) Bytes() ([16]byte, error)
```

Bytes converts a Key to a \[16\]byte UUID.

<a name="Key.Decode"></a>
### func \(Key\) [Decode](<https://github.com/agentstation/uuidkey/blob/master/key.go#L258>)

```go
func (k Key) Decode() (string, error)
```

Decode will decode a given Key into a UUID string with basic length validation.

<a name="Key.IsValid"></a>
### func \(Key\) [IsValid](<https://github.com/agentstation/uuidkey/blob/master/key.go#L90>)

```go
func (k Key) IsValid() bool
```

IsValid verifies if a given Key follows the correct format. The format should be:

- 31 characters long \(with hyphens\) or 28 characters \(without hyphens\)
- Uppercase
- Contains only alphanumeric characters
- Contains 3 hyphens \(if hyphenated\)
- Each part is 7 characters long
- Each part contains only valid crockford base32 characters \(I, L, O, U are not allowed\)

<a name="Key.String"></a>
### func \(Key\) [String](<https://github.com/agentstation/uuidkey/blob/master/key.go#L40>)

```go
func (k Key) String() string
```

String will convert your Key into a string.

<a name="Key.UUID"></a>
### func \(Key\) [UUID](<https://github.com/agentstation/uuidkey/blob/master/key.go#L129>)

```go
func (k Key) UUID() (string, error)
```

UUID will validate and convert a given Key into a UUID string.

<a name="Option"></a>
## type [Option](<https://github.com/agentstation/uuidkey/blob/master/key.go#L57>)

Option is a function that configures options

```go
type Option func(c *config)
```

<a name="With128BitEntropy"></a>With128BitEntropy expects 128 bits of entropy in the APIKey

```go
var With128BitEntropy Option = func(c *config) {
    c.entropySize = EntropyBits128
}
```

<a name="With160BitEntropy"></a>With160BitEntropy expects 160 bits of entropy in the APIKey

```go
var With160BitEntropy Option = func(c *config) {
    c.entropySize = EntropyBits160
}
```

<a name="With256BitEntropy"></a>With256BitEntropy expects 256 bits of entropy in the APIKey

```go
var With256BitEntropy Option = func(c *config) {
    c.entropySize = EntropyBits256
}
```

<a name="WithoutHyphens"></a>WithoutHyphens expects no hyphens in the Key

```go
var WithoutHyphens Option = func(c *config) {
    c.hyphens = false
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

## Makefile

```sh
jack@devbox ➜ make help

Usage:
  make <target>

General
  help                  Display the list of targets and their descriptions

Tooling
  install-devbox        Install Devbox
  devbox-update         Update Devbox
  devbox                Run Devbox shell

Installation
  install               Download go modules

Development
  fmt                   Run go fmt
  generate              Generate and embed go documentation into README.md
  vet                   Run go vet
  lint                  Run golangci-lint

Benchmarking, Testing, & Coverage
  bench                 Run Go benchmarks
  test                  Run Go tests
  coverage              Run tests and generate coverage report
```

## Benchmarks

> **Note:** These benchmarks were run on an Apple M2 Max CPU with 12 cores (8 performance and 4 efficiency) and 32 GB of memory, running macOS 14.6.1. The results are not representative of all systems, but should give you a general idea of the performance of the package. I was running other processes on the machine while running the benchmarks.

*Your mileage may vary.*

```sh
devbox ➜ make bench
Running go benchmarks...
go test ./... -tags=bench -bench=.
goos: darwin
goarch: arm64
pkg: github.com/agentstation/uuidkey
BenchmarkValidate-12                     	40536168	        29.66 ns/op
BenchmarkValidateInvalid-12              	820111176	         1.439 ns/op
BenchmarkParse-12                        	38579367	        29.78 ns/op
BenchmarkParseInvalid-12                 	668715536	         1.777 ns/op
BenchmarkUUID-12                         	 4767662	       250.6 ns/op
BenchmarkUUIDInvalid-12                  	66099007	        17.17 ns/op
BenchmarkEncode-12                       	 7224882	       162.9 ns/op
BenchmarkDecode-12                       	 5363344	       220.4 ns/op
BenchmarkBytes-12                        	 5438636	       217.0 ns/op
BenchmarkEncodeBytes-12                  	12546094	        94.13 ns/op
BenchmarkValidateWithHyphens-12          	38416134	        29.51 ns/op
BenchmarkValidateWithoutHyphens-12       	39153255	        29.23 ns/op
BenchmarkParseWithHyphens-12             	38402560	        30.19 ns/op
BenchmarkParseWithoutHyphens-12          	38653306	        29.85 ns/op
BenchmarkEncodeWithHyphens-12            	 7146574	       164.0 ns/op
BenchmarkEncodeWithoutHyphens-12         	 7255610	       163.3 ns/op
BenchmarkDecodeWithHyphens-12            	 5368426	       221.0 ns/op
BenchmarkDecodeWithoutHyphens-12         	 5370716	       221.0 ns/op
BenchmarkBytesWithHyphens-12             	 5430710	       220.6 ns/op
BenchmarkBytesWithoutHyphens-12          	 5032964	       217.1 ns/op
BenchmarkEncodeBytesWithHyphens-12       	12419739	        92.85 ns/op
BenchmarkEncodeBytesWithoutHyphens-12    	12544892	        92.87 ns/op
BenchmarkString-12                       	1000000000	         0.2875 ns/op
BenchmarkValidateInvalidFormat-12        	824398530	         1.434 ns/op
BenchmarkParseInvalidFormat-12           	668982553	         1.777 ns/op
BenchmarkDecodeInvalidFormat-12          	10187534	       114.5 ns/op
BenchmarkEncodeInvalidUUID-12            	11438924	       102.3 ns/op
BenchmarkBytesInvalidFormat-12           	10280540	       113.4 ns/op
PASS
ok  	github.com/agentstation/uuidkey	36.679s
```

