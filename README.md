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

## UUID Library Compatibility

This package is designed to work with any UUID that follows the official UUID specification ([RFC 4122](https://tools.ietf.org/html/rfc4122)). If your UUID implementation adheres to this standard, it should be compatible with this package.

We specifically test and maintain compatibility with the following UUID libraries:

- [github.com/gofrs/uuid](https://github.com/gofrs/uuid) (v4.4.0+)
- [github.com/google/uuid](https://github.com/google/uuid) (v1.6.0+)

While we officially support and test against these specific versions, any UUID library that follows the RFC 4122 specification should work with this package.

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

2. Encode a UUID string to a Key type:

```go
key, _ := uuidkey.Encode("d1756360-5da0-40df-9926-a76abff5601d")
fmt.Println(key) // Output: 38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X
```

3. Decode a Key type to a UUID string with Key format validation:

```go
key, _ := uuidkey.Parse("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
uuid, err := key.UUID()
if err != nil {
    log.Fatal("Error:", err)
}
fmt.Println(uuid) // Output: d1756360-5da0-40df-9926-a76abff5601d
```

4. Decode a Key type to a UUID string with only basic Key length validation:

```go
key, _ := uuidkey.Parse("38QARV0-1ET0G6Z-2CJD9VA-2ZZAR0X")
uuid, err := key.Decode()
if err != nil {
    log.Fatal("Error:", err)
}
fmt.Println(uuid) // Output: d1756360-5da0-40df-9926-a76abff5601d
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
- [type Key](<#Key>)
  - [func Encode\(uuid string\) \(Key, error\)](<#Encode>)
  - [func EncodeBytes\(uuid \[16\]byte\) \(Key, error\)](<#EncodeBytes>)
  - [func Parse\(key string\) \(Key, error\)](<#Parse>)
  - [func \(k Key\) Bytes\(\) \(\[16\]byte, error\)](<#Key.Bytes>)
  - [func \(k Key\) Decode\(\) \(string, error\)](<#Key.Decode>)
  - [func \(k Key\) String\(\) string](<#Key.String>)
  - [func \(k Key\) UUID\(\) \(string, error\)](<#Key.UUID>)
  - [func \(k Key\) Valid\(\) bool](<#Key.Valid>)


## Constants

<a name="KeyLength"></a>Key validation constraint constants

```go
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
```

<a name="Key"></a>
## type [Key](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L34>)

Key is a UUID Key string.

```go
type Key string
```

<a name="Encode"></a>
### func [Encode](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L118>)

```go
func Encode(uuid string) (Key, error)
```

Encode will encode a given UUID string into a Key with basic length validation.

<a name="EncodeBytes"></a>
### func [EncodeBytes](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L147>)

```go
func EncodeBytes(uuid [16]byte) (Key, error)
```

EncodeBytes encodes a \[16\]byte UUID into a Key.

<a name="Parse"></a>
### func [Parse](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L42>)

```go
func Parse(key string) (Key, error)
```

Parse converts a Key formatted string into a Key type.

<a name="Key.Bytes"></a>
### func \(Key\) [Bytes](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L196>)

```go
func (k Key) Bytes() ([16]byte, error)
```

Bytes converts a Key to a \[16\]byte UUID.

<a name="Key.Decode"></a>
### func \(Key\) [Decode](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L167>)

```go
func (k Key) Decode() (string, error)
```

Decode will decode a given Key into a UUID string with basic length validation.

<a name="Key.String"></a>
### func \(Key\) [String](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L37>)

```go
func (k Key) String() string
```

String will convert your Key into a string.

<a name="Key.UUID"></a>
### func \(Key\) [UUID](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L95>)

```go
func (k Key) UUID() (string, error)
```

UUID will validate and convert a given Key into a UUID string.

<a name="Key.Valid"></a>
### func \(Key\) [Valid](<https://github.com/agentstation/uuidkey/blob/master/uuidkey.go#L68>)

```go
func (k Key) Valid() bool
```

Valid verifies if a given Key follows the correct format. The format should be:

- 31 characters long
- Uppercase
- Contains only alphanumeric characters
- Contains 3 hyphens
- Each part is 7 characters long

Examples of valid keys:

- 38QARV0\-1ET0G6Z\-2CJD9VA\-2ZZAR0X
- ABCDEFG\-1234567\-HIJKLMN\-OPQRSTU

Examples of invalid keys:

- 38QARV0\-1ET0G6Z\-2CJD9VA\-2ZZAR0 \(too short\)
- 38qarv0\-1ET0G6Z\-2CJD9VA\-2ZZAR0X \(contains lowercase\)
- 38QARV0\-1ET0G6Z\-2CJD9VA\-2ZZAR0X\- \(extra hyphen\)
- 38QARV0\-1ET0G6Z\-2CJD9VA2ZZAR0X \(missing hyphen\)
- 38QARV0\-1ET0G6\-2CJD9VA\-2ZZAR0X \(part too short\)

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

> **Note:** These benchmarks were run on an Apple M2 Max CPU with 12 cores (8 performance and 4 efficiency) and 32 GB of memory, running macOS 14.6.1.

*Your mileage may vary.*

```sh
make bench
Running go benchmarks...
go test ./... -tags=bench -bench=.
goos: darwin
goarch: arm64
pkg: github.com/agentstation/uuidkey
cpu: Apple M2 Max
BenchmarkValidate-12           	33527211	        35.72 ns/op
BenchmarkParse-12              	32329798	        36.96 ns/op
BenchmarkFromKey-12            	 4886846	       250.6 ns/op
BenchmarkEncode-12             	 3151844	       377.0 ns/op
BenchmarkDecode-12             	 5587066	       216.7 ns/op
BenchmarkValidateInvalid-12    	1000000000	         0.2953 ns/op
BenchmarkParseValid-12         	32424325	        36.89 ns/op
BenchmarkParseInvalid-12       	70131522	        17.01 ns/op
BenchmarkUUIDValid-12          	 4693452	       247.2 ns/op
BenchmarkUUIDInvalid-12        	70141429	        16.92 ns/op
PASS
ok  	github.com/agentstation/uuidkey	13.365s
```
