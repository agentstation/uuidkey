package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/agentstation/uuidkey"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"
)

var (
	uuidVersion int
	timeFlag    string
)

var uuidCmd = &cobra.Command{
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

func init() {
	uuidCmd.Flags().IntVarP(&uuidVersion, "version", "v", 4, "UUID version: 4, 6, or 7")
	uuidCmd.Flags().StringVarP(&timeFlag, "time", "t", "", "Custom timestamp for v6/v7 (RFC3339 format)")
	rootCmd.AddCommand(uuidCmd)
}

func runUUID(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Generate new UUID
		u, err := generateUUID()
		if err != nil {
			outputError(err)
			return err
		}
		
		// Also encode to key
		key, err := uuidkey.Encode(u, getKeyOptions()...)
		if err != nil {
			err = fmt.Errorf("error encoding UUID: %w", err)
			outputError(err)
			return err
		}
		
		data := map[string]interface{}{
			"uuid":    u,
			"key":     key.String(),
			"version": uuidVersion,
		}
		
		// Add timestamp if provided
		if timeFlag != "" {
			data["timestamp"] = timeFlag
		}
		
		if quietMode {
			output(data, u)
		} else if jsonOutput {
			output(data, "")
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "UUID: %s\n", u)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Key:  %s\n", key.String())
		}
		return nil
	}
	
	input := args[0]
	inputType := detectInputType(input)
	
	switch inputType {
	case "uuid":
		// Encode UUID to key
		key, err := uuidkey.Encode(input, getKeyOptions()...)
		if err != nil {
			err = fmt.Errorf("error encoding UUID: %w", err)
			outputError(err)
			return err
		}
		data := map[string]interface{}{"uuid": input, "key": key.String()}
		if quietMode {
			output(data, key.String())
		} else if jsonOutput {
			output(data, "")
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "UUID: %s\n", input)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Key:  %s\n", key.String())
		}
	case "key":
		// Decode key to UUID
		parsedKey, err := uuidkey.Parse(input)
		if err != nil {
			err = fmt.Errorf("error parsing key: %w", err)
			outputError(err)
			return err
		}
		u, err := parsedKey.Decode()
		if err != nil {
			err = fmt.Errorf("error decoding key: %w", err)
			outputError(err)
			return err
		}
		data := map[string]interface{}{"key": input, "uuid": u}
		if quietMode {
			output(data, u)
		} else if jsonOutput {
			output(data, "")
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Key:  %s\n", input)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "UUID: %s\n", u)
		}
	case "apikey":
		// Parse API key to extract UUID
		apiKey, err := uuidkey.ParseAPIKey(input)
		if err != nil {
			err = fmt.Errorf("error parsing API key: %w", err)
			outputError(err)
			return err
		}
		u, err := apiKey.Key.Decode()
		if err != nil {
			err = fmt.Errorf("error decoding API key: %w", err)
			outputError(err)
			return err
		}
		data := map[string]interface{}{"apikey": input, "uuid": u}
		if quietMode {
			output(data, u)
		} else if jsonOutput {
			output(data, "")
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "API Key: %s\n", input)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "UUID:    %s\n", u)
		}
	default:
		err := fmt.Errorf("invalid input: input doesn't match UUID or key pattern")
		outputError(err)
		return err
	}
	
	return nil
}

// generateUUID generates a new UUID using gofrs/uuid based on the version flag
func generateUUID() (string, error) {
	var u uuid.UUID
	var err error
	
	// Check if custom timestamp was requested
	if timeFlag != "" {
		return "", fmt.Errorf("custom timestamp is not supported in gofrs/uuid v4.4.0 (requires v5+)")
	}
	
	switch uuidVersion {
	case 4:
		u, err = uuid.NewV4()
	case 6:
		u, err = uuid.NewV6()
	case 7:
		u, err = uuid.NewV7()
	default:
		return "", fmt.Errorf("unsupported UUID version: %d (supported: 4, 6, 7)", uuidVersion)
	}
	
	if err != nil {
		return "", fmt.Errorf("error generating UUID: %w", err)
	}
	return u.String(), nil
}

// detectInputType determines if input is a UUID, key, or API key
func detectInputType(input string) string {
	// Remove any whitespace
	input = strings.TrimSpace(input)
	
	// UUID pattern (with or without hyphens)
	uuidPattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{4}-?[0-9a-fA-F]{12}$`)
	if uuidPattern.MatchString(input) {
		return "uuid"
	}
	
	// API key pattern: PREFIX_KEY_CHECKSUM
	if strings.Count(input, "_") == 2 {
		parts := strings.Split(input, "_")
		if len(parts) == 3 && len(parts[2]) == 8 {
			return "apikey"
		}
	}
	
	// Key pattern (26-31 uppercase alphanumeric chars, possibly with hyphens)
	keyPattern := regexp.MustCompile(`^[0-9A-Z\-]{26,31}$`)
	if keyPattern.MatchString(strings.ToUpper(input)) {
		return "key"
	}
	
	return "unknown"
}