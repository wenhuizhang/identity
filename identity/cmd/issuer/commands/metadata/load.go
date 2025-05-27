// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
)

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
		if loadCmdIn.MetadataID == "" {
			fmt.Fprintf(os.Stderr, "Metadata ID: ")
			_, err := fmt.Scanln(&loadCmdIn.MetadataID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading metadata ID: %v\n", err)
				return
			}
		}
		if loadCmdIn.MetadataID == "" {
			fmt.Fprintf(os.Stderr, "No metadata ID provided\n")
			return
		}

		// check the metadata id is valid
		metadata, err := metadataService.GetMetadata(cache.VaultId, cache.KeyID, cache.IssuerId, loadCmdIn.MetadataID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting metadata: %v\n", err)
			return
		}
		if metadata == nil {
			fmt.Fprintf(os.Stderr, "No metadata found with ID: %s\n", loadCmdIn.MetadataID)
			return
		}

		// save the metadata id to the cache
		cache.MetadataId = metadata.ID
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded metadata with ID: %s\n", loadCmdIn.MetadataID)

	},
}
