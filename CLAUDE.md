# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

UUIDKey is a Go library that encodes UUIDs into human-readable keys using Base32-Crockford encoding. It provides secure API key generation with configurable entropy levels and includes CRC32 checksum validation.

## Key Commands

### Development
```bash
# Download dependencies
make deps

# Install the uuidkey binary
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

# Run go vet for static analysis
make vet

# View all available commands
make help

# Build the CLI binary
make build

# Build binaries for all platforms
make build-all

# Test release process locally
make release-snapshot

# Create a version tag
make tag VERSION=v1.2.3

# Create a release (requires tag)
make release
```

### Development Environment
```bash
# Set up development environment (requires Devbox)
devbox shell

# Or use make command
make devbox
```

### CLI Usage
```bash
# Generate operations (no input = generate)
uuidkey uuid                     # Generate new UUID v4 (default)
uuidkey uuid -v 6                # Generate UUID v6 (k-sortable)
uuidkey uuid -v 7                # Generate UUID v7 (k-sortable, millisecond precision)
uuidkey key                      # Generate new key  
uuidkey key --no-hyphens         # Generate key without hyphens
uuidkey apikey -p MYAPP          # Generate API key with prefix
uuidkey apikey -p MYAPP -e 256   # Generate with 256-bit entropy

# Smart operations (auto-detect input type)
uuidkey uuid <uuid-or-key>       # Encode UUID or decode key
uuidkey key <uuid>               # Encode UUID to key

# Parse operations
uuidkey apikey <apikey>          # Parse and validate API key

# Output formats
uuidkey uuid --json              # JSON output
uuidkey uuid -q                  # Quiet mode (value only)
uuidkey key -v 7 --json          # Generate v7 key with JSON output

# Explicit operations
uuidkey encode <uuid>            # Explicitly encode
uuidkey decode <key>             # Explicitly decode
```

## Architecture

### Core Components

1. **CLI Application** (`cmd/uuidkey/`):
   - Built with Cobra framework for better command organization
   - Smart command-line interface with intuitive defaults (no input = generate)
   - Commands: `uuid`, `key`, `apikey`, `encode`, `decode`, `version`
   - UUID versions: v4 (random), v6 (k-sortable), v7 (k-sortable, millisecond precision)
   - Flags: `--prefix/-p`, `--entropy/-e`, `--version/-v`, `--no-hyphens`, `--json`, `--quiet/-q`
   - Enhanced version command with binary checksum verification
   - Build information embedded via ldflags (version, commit, date, built by)
   - Cross-platform binary distribution via GoReleaser

2. **Key Encoding/Decoding** (`key.go`):
   - `Key` type represents an encoded UUID
   - `Encode(uuid.UUID)` and `Decode(string)` for basic conversion
   - `EncodeWithHyphens()` and `DecodeWithHyphens()` for hyphenated format
   - Base32-Crockford encoding ensures URL-safe, human-readable keys

3. **API Key Generation** (`apikey.go`):
   - `APIKey` struct with Prefix, Key, Entropy (128/160/256 bits), and CRC32 Checksum
   - Format: `PREFIX_KEY+ENTROPY_CHECKSUM` (e.g., `AGNTSTNP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJBJLSO7WBNWY3F_A1B2C3D8`)
   - `NewAPIKey(prefix, options...)` generates secure keys with configurable entropy
   - `ParseAPIKey(string)` validates and parses API keys
   - Uses BLAKE2b for entropy generation and follows GitHub Secret Scanning format

4. **Crockford Base32 Implementation** (`crockford.go`):
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
- Race condition detection enabled in coverage tests (`make coverage`)
- CI/CD via GitHub Actions runs tests on every push/PR
- Target coverage: 95.8%

### Dependencies

- `golang.org/x/crypto` - BLAKE2b hashing for entropy (v0.35.0+)
- Test dependencies: `gofrs/uuid` and `google/uuid`
- No external dependencies for Base32-Crockford encoding (uses standard library)

### Release Process

1. **Automated Releases**: GitHub Actions workflow triggers on version tags (v*)
2. **Cross-Platform Binaries**: Builds for Linux, macOS, Windows, FreeBSD (amd64, arm64, arm)
3. **Release Artifacts**: Binaries, checksums.txt, changelog, archives (tar.gz/zip)
4. **Version Injection**: Build metadata embedded via ldflags
5. **Local Testing**: Use `make release-snapshot` before creating actual release

### CI/CD Configuration

- **GoReleaser** (`.goreleaser.yaml`): Handles building, archiving, and publishing
- **GitHub Actions** (`.github/workflows/release.yml`): Automates the release pipeline
- **Makefile**: Provides convenient commands for the release workflow

## Recent Changes

### v1.2.0 (unreleased)
- Migrated CLI to Cobra framework for better command organization
- Added smart CLI with intuitive command structure (no input = generate)
- Commands: `uuid`, `key`, `apikey` with automatic input type detection
- Added UUID version support: v4 (random), v6 (k-sortable), v7 (k-sortable, millisecond)
- Enhanced version command with binary checksum verification
- Added flags: `--version/-v`, `--no-hyphens`, `--json`, `--quiet` for flexible output
- Uses `github.com/gofrs/uuid` for proper UUID generation
- Build information embedded via ldflags with VCS info from runtime/debug
- Implemented automated release process with GoReleaser
- Added cross-platform binary distribution via GitHub releases
- Integrated goreleaser into devbox environment

### v1.1.0 (unreleased)
- Replaced external `github.com/richardlehane/crock32` dependency with standard library implementation
- Custom `crock32Encode`/`crock32Decode` functions maintain backward compatibility
- Updated `golang.org/x/crypto` to v0.35.0 to fix CVE-2025-22869
- Improved test coverage to 95.8%
- Performance optimized with fixed-size buffers and lookup tables
- Performance results: encoding ~136.5ns/op, decoding ~282.4ns/op (comparable to original external library)
- Consolidated test files to reduce duplication