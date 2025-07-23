package main

import (
	"fmt"

	"github.com/agentstation/uuidkey"
	"github.com/spf13/cobra"
)

var (
	prefix  string
	entropy int
)

var apikeyCmd = &cobra.Command{
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

func init() {
	apikeyCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "API key prefix (required for generation)")
	apikeyCmd.Flags().IntVarP(&entropy, "entropy", "e", 128, "Entropy bits: 128, 160, or 256")
	apikeyCmd.Flags().IntVarP(&uuidVersion, "version", "v", 4, "UUID version: 4, 6, or 7")
	rootCmd.AddCommand(apikeyCmd)
}

func runAPIKey(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// Generate new API key
		if prefix == "" {
			err := fmt.Errorf("either provide an API key to parse or use --prefix to generate a new one")
			outputError(err)
			return err
		}
		
		u, err := generateUUID()
		if err != nil {
			outputError(err)
			return err
		}
		// Validate entropy value
		if entropy != 128 && entropy != 160 && entropy != 256 {
			err := fmt.Errorf("invalid entropy value: %d (must be 128, 160, or 256)", entropy)
			outputError(err)
			return err
		}
		
		apiKey, err := uuidkey.NewAPIKey(prefix, u, getEntropyOption())
		if err != nil {
			err = fmt.Errorf("error generating API key: %w", err)
			outputError(err)
			return err
		}
		
		data := map[string]interface{}{
			"apikey":       apiKey.String(),
			"prefix":       apiKey.Prefix,
			"uuid":         u,
			"uuid_version": uuidVersion,
			"entropy_bits": entropy,
			"key":          apiKey.Key.String(),
			"entropy":      apiKey.Entropy,
			"checksum":     apiKey.Checksum,
			"full_key":     apiKey.String(),
		}
		output(data, apiKey.String())
		return nil
	}
	
	// Parse existing API key
	apiKeyStr := args[0]
	apiKey, err := uuidkey.ParseAPIKey(apiKeyStr)
	if err != nil {
		err = fmt.Errorf("error parsing API key: %w", err)
		outputError(err)
		return err
	}
	
	// Decode the key to get the UUID
	u, _ := apiKey.Key.Decode()
	
	// Calculate entropy bits based on entropy string length
	entropyBits := 0
	switch len(apiKey.Entropy) {
	case 14:
		entropyBits = 128
	case 21:
		entropyBits = 160
	case 42:
		entropyBits = 256
	}
	
	data := map[string]interface{}{
		"valid":        true,
		"prefix":       apiKey.Prefix,
		"uuid":         u,
		"key":          apiKey.Key.String(),
		"entropy_bits": entropyBits,
		"checksum":     apiKey.Checksum,
	}
	
	out := cmd.OutOrStdout()
	
	if jsonOutput {
		output(data, "")
	} else if quietMode {
		_, _ = fmt.Fprintln(out, "Valid")
	} else {
		_, _ = fmt.Fprintln(out, "Valid API Key")
		_, _ = fmt.Fprintf(out, "Prefix:   %s\n", apiKey.Prefix)
		_, _ = fmt.Fprintf(out, "UUID:     %s\n", u)
		_, _ = fmt.Fprintf(out, "Key:      %s\n", apiKey.Key)
		_, _ = fmt.Fprintf(out, "Entropy:  %d bits\n", entropyBits)
		_, _ = fmt.Fprintf(out, "Checksum: %s\n", apiKey.Checksum)
	}
	
	return nil
}

// getEntropyOption returns the entropy option based on flag
func getEntropyOption() uuidkey.Option {
	switch entropy {
	case 160:
		return uuidkey.With160BitEntropy
	case 256:
		return uuidkey.With256BitEntropy
	default:
		return uuidkey.With128BitEntropy
	}
}