package main

import (
	"fmt"
	
	"github.com/agentstation/uuidkey"
	"github.com/spf13/cobra"
)

var encodeCmd = &cobra.Command{
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

func init() {
	encodeCmd.Flags().BoolVar(&noHyphens, "no-hyphens", false, "Omit hyphens in output")
	rootCmd.AddCommand(encodeCmd)
}

func runEncode(cmd *cobra.Command, args []string) error {
	uuid := args[0]
	key, err := uuidkey.Encode(uuid, getKeyOptions()...)
	if err != nil {
		err = fmt.Errorf("error encoding UUID: %w", err)
		outputError(err)
		return err
	}
	
	output(map[string]interface{}{"uuid": uuid, "key": key.String()}, key.String())
	return nil
}