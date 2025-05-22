// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	issuerFilesystem "github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	"github.com/agntcy/identity/internal/issuer/metadata"
	metadataFilesystem "github.com/agntcy/identity/internal/issuer/metadata/data/filesystem"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

var metadataFilesystemRepository = metadataFilesystem.NewMetadataFilesystemRepository()
var issuerFilesystemRepository = issuerFilesystem.NewIssuerFilesystemRepository()
var metadataService = metadata.NewMetadataService(metadataFilesystemRepository, issuerFilesystemRepository)

//nolint:lll // Allow long lines for CLI
var MetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Issue and publish important metadata for your Agent and MCP Server identities",
	Long: `
The metadata command is used to issue and publish important metadata for your Agent and MCP Server identities. With it you can:

- (generate) Generate new metadata
- (list) List all of the existing metadata for the issuer
- (show) Show details of a specific metadata
- (load) Load a metadata configuration
- (forget) Forget a specific metadata
`,
}

//nolint:mnd // Allow magic number for args
var metadataGenerateCmd = &cobra.Command{
	Use:   "generate [idp_client_id] [idp_client_secret] [idp_issuer_url]",
	Short: "Generate new metadata for your Agent and MCP Server identities",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}

		clientID := args[0]
		clientSecret := args[1]
		issuerURL := args[2]
		idpConfig := issuerTypes.IdpConfig{
			ClientId:     clientID,
			ClientSecret: clientSecret,
			IssuerUrl:    issuerURL,
		}

		metadataId, err := metadataService.GenerateMetadata(cache.VaultId, cache.IssuerId, &idpConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating metadata: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Generated metadata with ID: %s\n", metadataId)

		// Update the cache with the new metadata ID
		cache.MetadataId = metadataId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
			return
		}
	},
}

var metadataListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing metadata",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}

		allMetadata, err := metadataService.GetAllMetadata(cache.VaultId, cache.IssuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing metadata: %v\n", err)
			return
		}
		if len(allMetadata) == 0 {
			fmt.Fprintf(os.Stdout, "%s\n", "No metadata found")
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Existing metadata ids:")
		for _, metadata := range allMetadata {
			fmt.Fprintf(os.Stdout, "- %s\n", metadata.Id)
		}
	},
}

var metadataShowCmd = &cobra.Command{
	Use:   "show [metadata_id]",
	Short: "Show the chosen metadata",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}

		metadataId := args[0]

		metadata, err := metadataService.GetMetadata(cache.VaultId, cache.IssuerId, metadataId)
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

var metadataForgetCmd = &cobra.Command{
	Use:   "forget [metadata_id]",
	Short: "Forget the chosen metadata",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}

		metadataId := args[0]

		err = metadataService.ForgetMetadata(cache.VaultId, cache.IssuerId, metadataId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting metadata: %v\n", err)
			return
		}

		// If the metadata was the current metadata in the cache, clear the cache of metadata, and badge IDs
		if cache.MetadataId == metadataId {
			cache.MetadataId = ""
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot metadata with ID: %s\n", metadataId)

	},
}

var metadataLoadCmd = &cobra.Command{
	Use:   "load [metadata_id]",
	Short: "Load a metadata configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(os.Stderr, "No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}

		metadata_id := args[0]

		// check the metadata id is valid
		metadata, err := metadataService.GetMetadata(cache.VaultId, cache.IssuerId, metadata_id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting metadata: %v\n", err)
			return
		}
		if metadata == nil {
			fmt.Fprintf(os.Stderr, "No metadata found with ID: %s\n", metadata_id)
			return
		}

		// save the metadata id to the cache
		cache.MetadataId = metadata.Id
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded metadata with ID: %s\n", metadata_id)

	},
}

func init() {
	MetadataCmd.AddCommand(metadataGenerateCmd)
	MetadataCmd.AddCommand(metadataListCmd)
	MetadataCmd.AddCommand(metadataShowCmd)
	MetadataCmd.AddCommand(metadataForgetCmd)
	MetadataCmd.AddCommand(metadataLoadCmd)
}
