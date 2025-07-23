package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// Global flags
var (
	jsonOutput bool
	quietMode  bool
)

var rootCmd = &cobra.Command{
	Use:   "uuidkey",
	Short: "UUID to Base32-Crockford Key encoder/decoder",
	Long: `uuidkey is a fast, secure tool for UUID generation and Base32-Crockford encoding.

It provides:
  - UUID generation (v4, v6, v7)
  - Base32-Crockford encoding/decoding
  - API key generation with configurable entropy
  - Smart input detection (no input = generate)`,
}

func init() {
	// Global persistent flags
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "Minimal output")
	
	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// output handles formatting output based on flags
func output(data map[string]interface{}, defaultOutput string) {
	// Use the current command if set (for tests), otherwise use rootCmd
	cmd := getCurrentCommand()
	out := cmd.OutOrStdout()
	
	if jsonOutput {
		encoder := json.NewEncoder(out)
		// Use compact output to match test expectations
		_ = encoder.Encode(data)
	} else if quietMode && defaultOutput != "" {
		_, _ = fmt.Fprintln(out, defaultOutput)
	} else if defaultOutput != "" {
		_, _ = fmt.Fprintln(out, defaultOutput)
	}
}

// outputError handles error output based on flags
func outputError(err error) {
	if jsonOutput {
		data := map[string]interface{}{
			"error": err.Error(),
		}
		output(data, "")
	}
}

