# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

UUIDKey is a Go library that encodes UUIDs into human-readable keys using Base32-Crockford encoding. It provides secure API key generation with configurable entropy levels and includes CRC32 checksum validation.

## Key Commands

### Development
```bash
# Install the library
make install

# Run tests
make test

# Run a specific test
go test -run TestFunctionName

# Run benchmarks
make bench

# Generate code coverage
make coverage

# Lint the code
make lint

# Format code
make fmt

# Generate documentation
make generate
```

### Development Environment
```bash
# Set up development environment (requires Devbox)
devbox shell
```

## Architecture

### Core Components

1. **Key Encoding/Decoding** (`key.go`):
   - `Key` type represents an encoded UUID
   - `Encode(uuid.UUID)` and `Decode(string)` for basic conversion
   - `EncodeWithHyphens()` and `DecodeWithHyphens()` for hyphenated format
   - Base32-Crockford encoding ensures URL-safe, human-readable keys

2. **API Key Generation** (`apikey.go`):
   - `APIKey` struct with Prefix, Key, Entropy (128/160/256 bits), and CRC32 Checksum
   - `NewAPIKey(prefix, options...)` generates secure keys with configurable entropy
   - `ParseAPIKey(string)` validates and parses API keys
   - Uses BLAKE2b for entropy generation and follows GitHub Secret Scanning format

3. **Crockford Base32 Implementation** (`crockford.go`):
   - Custom implementation using Go's standard library `encoding/base32`
   - `crock32Encode(uint32)` and `crock32Decode(string)` for number-based encoding
   - Maintains backward compatibility with the original external library
   - Implements character normalization (O→0, I/L→1) per Crockford spec

### Key Design Patterns

- **Functional Options**: Configuration uses `Option` type (e.g., `WithEntropy(bits)`)
- **Immutable Types**: Key types are immutable with validation on creation
- **Error Handling**: All functions return descriptive errors for invalid inputs
- **Performance Focus**: Extensive benchmarking ensures optimal performance

### Testing Strategy

- Unit tests in `*_test.go` files cover all public APIs
- Benchmarks in `uuidkey_benchmark_test.go` track performance
- Race condition detection enabled in coverage tests
- CI/CD via GitHub Actions runs tests on every push/PR

### Dependencies

- `golang.org/x/crypto` - BLAKE2b hashing for entropy (v0.35.0+)
- Test dependencies: `gofrs/uuid` and `google/uuid`
- No external dependencies for Base32-Crockford encoding (uses standard library)

## Recent Changes

### v1.1.0 (unreleased)
- Replaced external `github.com/richardlehane/crock32` dependency with standard library implementation
- Custom `crock32Encode`/`crock32Decode` functions maintain backward compatibility
- Updated `golang.org/x/crypto` to v0.35.0 to fix CVE-2025-22869
- Improved test coverage to 95.8%
- Performance optimized with fixed-size buffers and lookup tables
- Performance results: encoding ~136.5ns/op, decoding ~282.4ns/op (comparable to original external library)
- Consolidated test files to reduce duplication