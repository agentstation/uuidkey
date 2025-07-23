package main

import (
	"github.com/spf13/cobra"
)

// newRootCommand creates a fresh root command instance
// This is used in tests to ensure complete isolation between test runs
func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuidkey",
		Short: "UUID to Base32-Crockford Key encoder/decoder",
		Long: `uuidkey is a fast, secure tool for UUID generation and Base32-Crockford encoding.

It provides:
  - UUID generation (v4, v6, v7)
  - Base32-Crockford encoding/decoding
  - API key generation with configurable entropy
  - Smart input detection (no input = generate)`,
	}

	// Add persistent flags
	cmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	cmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Minimal output")
	
	// Disable completion command in tests
	cmd.CompletionOptions.DisableDefaultCmd = true

	// Add commands
	cmd.AddCommand(newUUIDCommand())
	cmd.AddCommand(newKeyCommand())
	cmd.AddCommand(newAPIKeyCommand())
	cmd.AddCommand(newEncodeCommand())
	cmd.AddCommand(newDecodeCommand())
	cmd.AddCommand(newVersionCommand())

	return cmd
}

// newUUIDCommand creates a fresh uuid command instance
func newUUIDCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuid [input]",
		Short: "Generate or process UUIDs",
		Long: `Generate new UUIDs or encode/decode based on input.

With no input, generates a new UUID.
With a UUID as input, encodes it to a key.
With a key as input, decodes it to a UUID.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runUUID,
		Example: `  # Generate new UUID v4
  uuidkey uuid
  
  # Generate UUID v6 (k-sortable)
  uuidkey uuid -v 6
  
  # Generate UUID v7 (k-sortable, millisecond precision)
  uuidkey uuid -v 7
  
  # Encode UUID to key
  uuidkey uuid 550e8400-e29b-41d4-a716-446655440000
  
  # Decode key to UUID
  uuidkey uuid 1AGX100-3H9PGEM-2KHCH36-1AM8000`,
	}

	cmd.Flags().IntVarP(&uuidVersion, "version", "v", 4, "UUID version: 4, 6, or 7")
	cmd.Flags().StringVarP(&timeFlag, "time", "t", "", "Custom timestamp for v6/v7 (RFC3339 format)")

	return cmd
}

// newKeyCommand creates a fresh key command instance
func newKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key [input]",
		Short: "Generate or encode keys",
		Long: `Generate new Base32-Crockford keys or encode UUIDs to keys.

With no input, generates a new key from a random UUID.
With a UUID as input, encodes it to a key.
With a key as input, validates the key format.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runKey,
		Example: `  # Generate new key
  uuidkey key
  
  # Generate key without hyphens
  uuidkey key --no-hyphens
  
  # Encode UUID to key
  uuidkey key 550e8400-e29b-41d4-a716-446655440000
  
  # Validate existing key
  uuidkey key 1AGX100-3H9PGEM-2KHCH36-1AM8000`,
	}

	cmd.Flags().BoolVar(&noHyphens, "no-hyphens", false, "Omit hyphens in output")
	cmd.Flags().IntVarP(&uuidVersion, "version", "v", 4, "UUID version: 4, 6, or 7")

	return cmd
}

// newAPIKeyCommand creates a fresh apikey command instance
func newAPIKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apikey [apikey]",
		Short: "Generate or parse API keys",
		Long: `Generate new API keys with configurable entropy or parse existing ones.

With no input, generates a new API key (requires --prefix).
With an API key as input, parses and validates it.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runAPIKey,
		Example: `  # Generate API key
  uuidkey apikey -p MYAPP
  
  # Generate with 256-bit entropy
  uuidkey apikey -p MYAPP -e 256
  
  # Parse existing API key
  uuidkey apikey MYAPP_38QARV01ET0G6Z2CJD9VA2ZZAR0XJBJLSO7WBNWY3F_A1B2C3D8`,
	}

	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "API key prefix (required for generation)")
	cmd.Flags().IntVarP(&entropy, "entropy", "e", 128, "Entropy bits: 128, 160, or 256")
	cmd.Flags().IntVarP(&uuidVersion, "version", "v", 4, "UUID version: 4, 6, or 7")

	return cmd
}

// newEncodeCommand creates a fresh encode command instance
func newEncodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encode <uuid>",
		Short: "Encode UUID to Base32-Crockford key",
		Long:  `Explicitly encode a UUID to Base32-Crockford key format.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runEncode,
		Example: `  # Encode UUID to key
  uuidkey encode 550e8400-e29b-41d4-a716-446655440000
  
  # Encode without hyphens
  uuidkey encode 550e8400-e29b-41d4-a716-446655440000 --no-hyphens`,
	}

	cmd.Flags().BoolVar(&noHyphens, "no-hyphens", false, "Omit hyphens in output")

	return cmd
}

// newDecodeCommand creates a fresh decode command instance
func newDecodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode <key>",
		Short: "Decode Base32-Crockford key to UUID",
		Long:  `Explicitly decode a Base32-Crockford key to UUID format.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runDecode,
		Example: `  # Decode key to UUID
  uuidkey decode 1AGX100-3H9PGEM-2KHCH36-1AM8000
  
  # Decode key without hyphens
  uuidkey decode 1AGX1003H9PGEM2KHCH361AM8000`,
	}

	return cmd
}

// newVersionCommand creates a fresh version command instance
func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version and build information",
		Long: `Display detailed version, build information, and optionally verify binary integrity.

The version command shows:
  - Version number and build metadata
  - Git commit and build date
  - Go version and platform
  - Binary checksum verification (optional)`,
		RunE: runVersion,
		Example: `  # Show version
  uuidkey version
  
  # Verify binary checksum
  uuidkey version --verify
  
  # Verify against online checksums
  uuidkey version --verify-online`,
	}

	cmd.Flags().BoolVar(&verifyChecksum, "verify", false, "Verify binary checksum")
	cmd.Flags().BoolVar(&verifyOnline, "verify-online", false, "Verify against online checksums")

	return cmd
}

// setCurrentCommand sets a command instance as the current command for output function
// This is needed because the output function references rootCmd
var currentCmd *cobra.Command

func setCurrentCommand(cmd *cobra.Command) {
	currentCmd = cmd
}

// getCurrentCommand returns the current command or rootCmd as fallback
func getCurrentCommand() *cobra.Command {
	if currentCmd != nil {
		return currentCmd
	}
	return rootCmd
}