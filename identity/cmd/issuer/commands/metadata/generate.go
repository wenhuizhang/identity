// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

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
		if genCmdIn.IdpClientID == "" {
			fmt.Fprintf(os.Stderr, "IDP Client ID: ")
			_, err := fmt.Scanln(&genCmdIn.IdpClientID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IDP Client ID: %v\n", err)
				return
			}
		}
		if genCmdIn.IdpClientID == "" {
			fmt.Fprintf(os.Stderr, "No IDP Client ID provided\n")
			return
		}

		// if the idp client secret is not set, prompt the user for it interactively
		if genCmdIn.IdpClientSecret == "" {
			fmt.Fprintf(os.Stderr, "IDP Client Secret: ")
			_, err := fmt.Scanln(&genCmdIn.IdpClientSecret)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IDP Client Secret: %v\n", err)
				return
			}
		}
		if genCmdIn.IdpClientSecret == "" {
			fmt.Fprintf(os.Stderr, "No IDP Client Secret provided\n")
			return
		}

		// if the idp issuer url is not set, prompt the user for it interactively
		if genCmdIn.IdpIssuerURL == "" {
			fmt.Fprintf(os.Stderr, "IDP Issuer URL: ")
			_, err := fmt.Scanln(&genCmdIn.IdpIssuerURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IDP Issuer URL: %v\n", err)
				return
			}
		}
		if genCmdIn.IdpIssuerURL == "" {
			fmt.Fprintf(os.Stderr, "No IDP Issuer URL provided\n")
			return
		}

		idpConfig := issuerTypes.IdpConfig{
			ClientId:     genCmdIn.IdpClientID,
			ClientSecret: genCmdIn.IdpClientSecret,
			IssuerUrl:    genCmdIn.IdpIssuerURL,
		}

		metadataId, err := metadataService.GenerateMetadata(cmd.Context(), cache.VaultId, cache.IssuerId, &idpConfig)
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
