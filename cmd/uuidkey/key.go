package main

import (
	"fmt"

	"github.com/agentstation/uuidkey"
	"github.com/spf13/cobra"
)

var (
	noHyphens bool
)

var keyCmd = &cobra.Command{
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

func init() {
	keyCmd.Flags().BoolVar(&noHyphens, "no-hyphens", false, "Omit hyphens in output")
	keyCmd.Flags().IntVarP(&uuidVersion, "version", "v", 4, "UUID version: 4, 6, or 7")
	rootCmd.AddCommand(keyCmd)
}

func runKey(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Generate new key
		u, err := generateUUID()
		if err != nil {
			outputError(err)
			return err
		}
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
		
		// Output formatting
		if quietMode {
			output(data, key.String())
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
		// Output both UUID and key for consistency
		if quietMode {
			output(map[string]interface{}{"uuid": input, "key": key.String()}, key.String())
		} else if jsonOutput {
			output(map[string]interface{}{"uuid": input, "key": key.String()}, "")
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
		// Output both key and UUID for consistency
		if quietMode {
			output(map[string]interface{}{"key": input, "uuid": u}, u)
		} else if jsonOutput {
			output(map[string]interface{}{"key": input, "uuid": u}, "")
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Key:  %s\n", input)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "UUID: %s\n", u)
		}
	default:
		err := fmt.Errorf("invalid input: input doesn't match UUID or key pattern")
		outputError(err)
		return err
	}
	
	return nil
}

// getKeyOptions returns options for key encoding based on flags
func getKeyOptions() []uuidkey.Option {
	var opts []uuidkey.Option
	if noHyphens {
		opts = append(opts, uuidkey.WithoutHyphens)
	}
	return opts
}