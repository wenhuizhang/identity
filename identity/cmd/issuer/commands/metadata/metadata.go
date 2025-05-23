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

var (
	// setup the metadata service
	metadataFilesystemRepository = metadataFilesystem.NewMetadataFilesystemRepository()
	issuerFilesystemRepository   = issuerFilesystem.NewIssuerFilesystemRepository()
	metadataService              = metadata.NewMetadataService(metadataFilesystemRepository, issuerFilesystemRepository)

	// setup the command flags
	generateIdpClientId     string
	generateIdpClientSecret string
	generateIdpIssuerUrl    string
	showMetadataId          string
	forgetMetadataId        string
	loadMetadataId          string
)

var MetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Generate important metadata for your Agent and MCP Server identities",
	Long: `
The metadata command is used to generate important metadata for your Agent and MCP Server identities.
`,
}

//nolint:lll // Allow long lines for CLI
var metadataGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate new metadata for your Agent and MCP Server identities",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForMetadata()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the idp client id is not set, prompt the user for it interactively
		if generateIdpClientId == "" {
			fmt.Fprintf(os.Stderr, "IDP Client ID: ")
			_, err := fmt.Scanln(&generateIdpClientId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IDP Client ID: %v\n", err)
				return
			}
		}
		if generateIdpClientId == "" {
			fmt.Fprintf(os.Stderr, "No IDP Client ID provided\n")
			return
		}

		// if the idp client secret is not set, prompt the user for it interactively
		if generateIdpClientSecret == "" {
			fmt.Fprintf(os.Stderr, "IDP Client Secret: ")
			_, err := fmt.Scanln(&generateIdpClientSecret)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IDP Client Secret: %v\n", err)
				return
			}
		}
		if generateIdpClientSecret == "" {
			fmt.Fprintf(os.Stderr, "No IDP Client Secret provided\n")
			return
		}

		// if the idp issuer url is not set, prompt the user for it interactively
		if generateIdpIssuerUrl == "" {
			fmt.Fprintf(os.Stderr, "IDP Issuer URL: ")
			_, err := fmt.Scanln(&generateIdpIssuerUrl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IDP Issuer URL: %v\n", err)
				return
			}
		}
		if generateIdpIssuerUrl == "" {
			fmt.Fprintf(os.Stderr, "No IDP Issuer URL provided\n")
			return
		}

		idpConfig := issuerTypes.IdpConfig{
			ClientId:     generateIdpClientId,
			ClientSecret: generateIdpClientSecret,
			IssuerUrl:    generateIdpIssuerUrl,
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
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
	},
}

//nolint:lll // Allow long lines for CLI
var metadataListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing metadata",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForMetadata()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
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

//nolint:lll // Allow long lines for CLI
var metadataShowCmd = &cobra.Command{
	Use:   "show [metadata_id]",
	Short: "Show the chosen metadata",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForMetadata()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the metadata id is not set, prompt the user for it interactively
		if showMetadataId == "" {
			fmt.Fprintf(os.Stderr, "Metadata ID: ")
			_, err := fmt.Scanln(&showMetadataId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading metadata ID: %v\n", err)
				return
			}
		}
		if showMetadataId == "" {
			fmt.Fprintf(os.Stderr, "No metadata ID provided\n")
			return
		}

		metadata, err := metadataService.GetMetadata(cache.VaultId, cache.IssuerId, showMetadataId)
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

//nolint:lll // Allow long lines for CLI
var metadataForgetCmd = &cobra.Command{
	Use:   "forget [metadata_id]",
	Short: "Forget the chosen metadata",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForMetadata()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the metadata id is not set, prompt the user for it interactively
		if forgetMetadataId == "" {
			fmt.Fprintf(os.Stderr, "Metadata ID: ")
			_, err := fmt.Scanln(&forgetMetadataId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading metadata ID: %v\n", err)
				return
			}
		}
		if forgetMetadataId == "" {
			fmt.Fprintf(os.Stderr, "No metadata ID provided\n")
			return
		}

		err = metadataService.ForgetMetadata(cache.VaultId, cache.IssuerId, forgetMetadataId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting metadata: %v\n", err)
			return
		}

		// If the metadata was the current metadata in the cache, clear the cache of metadata, and badge IDs
		if cache.MetadataId == forgetMetadataId {
			cache.MetadataId = ""
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot metadata with ID: %s\n", forgetMetadataId)

	},
}

//nolint:lll // Allow long lines for CLI
var metadataLoadCmd = &cobra.Command{
	Use:   "load [metadata_id]",
	Short: "Load a metadata configuration",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault and issuer id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForMetadata()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the metadata id is not set, prompt the user for it interactively
		if loadMetadataId == "" {
			fmt.Fprintf(os.Stderr, "Metadata ID: ")
			_, err := fmt.Scanln(&loadMetadataId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading metadata ID: %v\n", err)
				return
			}
		}
		if loadMetadataId == "" {
			fmt.Fprintf(os.Stderr, "No metadata ID provided\n")
			return
		}

		// check the metadata id is valid
		metadata, err := metadataService.GetMetadata(cache.VaultId, cache.IssuerId, loadMetadataId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting metadata: %v\n", err)
			return
		}
		if metadata == nil {
			fmt.Fprintf(os.Stderr, "No metadata found with ID: %s\n", loadMetadataId)
			return
		}

		// save the metadata id to the cache
		cache.MetadataId = metadata.Id
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded metadata with ID: %s\n", loadMetadataId)

	},
}

func init() {
	metadataGenerateCmd.Flags().StringVarP(&generateIdpClientId, "idp-client-id", "i", "", "IDP Client ID")
	metadataGenerateCmd.Flags().StringVarP(&generateIdpClientSecret, "idp-client-secret", "s", "", "IDP Client Secret")
	metadataGenerateCmd.Flags().StringVarP(&generateIdpIssuerUrl, "idp-issuer-url", "u", "", "IDP Issuer URL")
	MetadataCmd.AddCommand(metadataGenerateCmd)

	MetadataCmd.AddCommand(metadataListCmd)

	metadataShowCmd.Flags().StringVarP(&showMetadataId, "metadata-id", "m", "", "The ID of the metadata to show")
	MetadataCmd.AddCommand(metadataShowCmd)

	metadataForgetCmd.Flags().StringVarP(&forgetMetadataId, "metadata-id", "m", "", "The ID of the metadata to forget")
	MetadataCmd.AddCommand(metadataForgetCmd)

	metadataLoadCmd.Flags().StringVarP(&loadMetadataId, "metadata-id", "m", "", "The ID of the metadata to load")
	MetadataCmd.AddCommand(metadataLoadCmd)
}
