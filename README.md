<!-- [![Sourcegraph](https://sourcegraph.com/github.com/agentstation/uuidkey/-/badge.svg?style=flat-square)](https://sourcegraph.com/github.com/agentstation/uuidkey?badge) -->
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/agentstation/uuidkey)
[![Go Report Card](https://goreportcard.com/badge/github.com/agentstation/uuidkey?style=flat-square)](https://goreportcard.com/report/github.com/agentstation/uuidkey)
<!-- [![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/agentstation/uuidkey/echo.yml?style=flat-square)](https://github.com/agentstation/uuidkey/actions)
[![Codecov](https://img.shields.io/codecov/c/github/agentstation/uuidkey.svg?style=flat-square)](https://codecov.io/gh/agentstation/uuidkey) -->
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/agentstation/uuidkey/master/LICENSE)
<!-- [![Forum](https://img.shields.io/badge/community-forum-00afd1.svg?style=flat-square)](https://github.com/agentstation/uuidkey/discussions) -->
<!-- [![Twitter](https://img.shields.io/badge/twitter-@agentstationHQ-55acee.svg?style=flat-square)](https://twitter.com/agentstationHQ) -->

```sh
                         _   _  _   _  ___ ____     _  __          
                        | | | || | | ||_ _|  _ \   | |/ /___ _   _ 
                        | | | || | | | | || | | |  | ' // _ \ | | |
                        | |_| || |_| | | || |_| |  | . \  __/ |_| |
                         \___/  \___/ |___|____/   |_|\_\___|\__, |
                                                             |___/ 
```

The `uuidkey` package encodes UUIDs to a readable `Key` format via the Base32-Crockford codec.

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# uuidkey

```go
import "github.com/agentstation/uuidkey"
```

Package uuidkey encodes UUIDs to a readable Key format via the Base32\-Crockford codec.

## Index

- [type Key](<#Key>)
  - [func Encode\(uuid string\) Key](<#Encode>)
  - [func FromString\(key string\) \(Key, error\)](<#FromString>)
  - [func \(k Key\) Decode\(\) string](<#Key.Decode>)
  - [func \(k Key\) String\(\) string](<#Key.String>)
  - [func \(k Key\) UUIDString\(\) \(string, error\)](<#Key.UUIDString>)
  - [func \(k Key\) Valid\(\) bool](<#Key.Valid>)


<a name="Key"></a>
## type Key

Key is a UUID Key string.

```go
type Key string
```

<a name="Encode"></a>
### func Encode

```go
func Encode(uuid string) Key
```

Encode will encode a given UUID string into a Key without validation.

<a name="FromString"></a>
### func FromString

```go
func FromString(key string) (Key, error)
```

FromString will convert a Key formatted string type into a Key type.

<a name="Key.Decode"></a>
### func \(Key\) Decode

```go
func (k Key) Decode() string
```

Decode will decode a given Key into a UUID string without validation.

<a name="Key.String"></a>
### func \(Key\) String

```go
func (k Key) String() string
```

String will convert your Key into a string.

<a name="Key.UUIDString"></a>
### func \(Key\) UUIDString

```go
func (k Key) UUIDString() (string, error)
```

UUIDString will validate and convert a given Key into a UUID string.

<a name="Key.Valid"></a>
### func \(Key\) Valid

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
make help

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

Testing & Benchmarking
  test                  Run Go tests
  bench                 Run Go benchmarks
```

## Benchmarks

> **Note:** These benchmarks were run on an Apple M2 Max CPU with 12 cores (8 performance and 4 efficiency) and 32 GB of memory, running macOS 14.6.1.

*Your mileage may vary.*

```sh
go test -bench=.
goos: darwin
goarch: arm64
pkg: github.com/agentstation/uuidkey
BenchmarkValidate-12                 33994471            35.02 ns/op
BenchmarkFromString-12               32470240            35.94 ns/op
BenchmarkFromKey-12                   4773018           253.2 ns/op
BenchmarkEncode-12                    3167922           371.5 ns/op
BenchmarkDecode-12                    5677419           211.7 ns/op
BenchmarkValidateInvalid-12          1000000000             0.2881 ns/op
BenchmarkFromStringValid-12          32319241            35.99 ns/op
BenchmarkFromStringInvalid-12        69830540            16.41 ns/op
BenchmarkUUIDStringValid-12           4940355           246.7 ns/op
BenchmarkUUIDStringInvalid-12        70641040            16.33 ns/op
PASS
ok      github.com/agentstation/uuidkey    13.168s
```