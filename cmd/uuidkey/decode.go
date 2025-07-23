package main

import (
	"fmt"
	
	"github.com/agentstation/uuidkey"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
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

func init() {
	rootCmd.AddCommand(decodeCmd)
}

func runDecode(cmd *cobra.Command, args []string) error {
	key := args[0]
	parsedKey, err := uuidkey.Parse(key)
	if err != nil {
		err = fmt.Errorf("error parsing key: %w", err)
		outputError(err)
		return err
	}
	
	uuid, err := parsedKey.Decode()
	if err != nil {
		err = fmt.Errorf("error decoding key: %w", err)
		outputError(err)
		return err
	}
	
	output(map[string]interface{}{"key": key, "uuid": uuid}, uuid)
	return nil
}