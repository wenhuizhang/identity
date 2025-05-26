// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
)

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
		if showCmdIn.MetadataID == "" {
			fmt.Fprintf(os.Stderr, "Metadata ID: ")
			_, err := fmt.Scanln(&showCmdIn.MetadataID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading metadata ID: %v\n", err)
				return
			}
		}
		if showCmdIn.MetadataID == "" {
			fmt.Fprintf(os.Stderr, "No metadata ID provided\n")
			return
		}

		metadata, err := metadataService.GetMetadata(cache.VaultId, cache.IssuerId, showCmdIn.MetadataID)
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
