// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	issuerMetadata "github.com/agntcy/identity/internal/issuer/metadata"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

//nolint:lll // Allow long lines for CLI
var MetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Issue and publish important metadata for your Agent and MCP Server identities",
	Long: `
The metadata command is used to issue and publish important metadata for your Agent and MCP Server identities. With it you can:

- (generate) Generate new metadata
- (list) List all of the existing metadata for the issuer
- (show) Show details of a specific metadata
- (forget) Forget a specific metadata
`,
}

//nolint:mnd // Allow magic number for args
var metadataGenerateCmd = &cobra.Command{
	Use:   "generate [vault_id] [issuer_id] [idp_client_id] [idp_client_secret] [idp_issuer_url]",
	Short: "Generate new metadata for your Agent and MCP Server identities",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		vaultId := args[0]
		issuerId := args[1]
		clientID := args[2]
		clientSecret := args[3]
		issuerURL := args[4]
		idpConfig := issuerTypes.IdpConfig{
			ClientId:     clientID,
			ClientSecret: clientSecret,
			IssuerUrl:    issuerURL,
		}

		_, err := issuerMetadata.GenerateMetadata(vaultId, issuerId, &idpConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating metadata: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Creating a new metadata")
	},
}

var metadataListCmd = &cobra.Command{
	Use:   "list [vault_id] [issuer_id]",
	Short: "List your existing metadata",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		vaultId := args[0]
		issuerId := args[1]

		metadataIds, err := issuerMetadata.ListMetadataIds(vaultId, issuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing metadata: %v\n", err)
			return
		}
		if len(metadataIds) == 0 {
			fmt.Fprintf(os.Stdout, "%s\n", "No metadata found")
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Existing metadata ids:")
		for _, metadataId := range metadataIds {
			fmt.Fprintf(os.Stdout, "- %s\n", metadataId)
		}
	},
}

//nolint:mnd // Allow magic number for args
var metadataShowCmd = &cobra.Command{
	Use:   "show [vault_id] [issuer_id] [metadata_id]",
	Short: "Show the chosen metadata",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		vaultId := args[0]
		issuerId := args[1]
		metadataId := args[2]

		metadata, err := issuerMetadata.GetMetadata(vaultId, issuerId, metadataId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting metadata: %v\n", err)
			return
		}
		metadataJSON, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling metadata to JSON: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(metadataJSON))
	},
}

//nolint:mnd // Allow magic number for args
var metadataForgetCmd = &cobra.Command{
	Use:   "forget [vault_id] [issuer_id] [metadata_id]",
	Short: "Forget the chosen metadata",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		vaultId := args[0]
		issuerId := args[1]
		metadataId := args[2]

		err := issuerMetadata.ForgetMetadata(vaultId, issuerId, metadataId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting metadata: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Forgot metadata with ID: %s\n", metadataId)
	},
}

func init() {
	MetadataCmd.AddCommand(metadataGenerateCmd)
	MetadataCmd.AddCommand(metadataListCmd)
	MetadataCmd.AddCommand(metadataShowCmd)
	MetadataCmd.AddCommand(metadataForgetCmd)
}
