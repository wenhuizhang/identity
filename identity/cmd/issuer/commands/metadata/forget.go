// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
)

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
		if forgCmdIn.MetadataID == "" {
			fmt.Fprintf(os.Stderr, "Metadata ID: ")
			_, err := fmt.Scanln(&forgCmdIn.MetadataID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading metadata ID: %v\n", err)
				return
			}
		}
		if forgCmdIn.MetadataID == "" {
			fmt.Fprintf(os.Stderr, "No metadata ID provided\n")
			return
		}

		err = metadataService.ForgetMetadata(cache.VaultId, cache.KeyID, cache.IssuerId, forgCmdIn.MetadataID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting metadata: %v\n", err)
			return
		}

		// If the metadata was the current metadata in the cache, clear the cache of metadata, and badge IDs
		if cache.MetadataId == forgCmdIn.MetadataID {
			cache.MetadataId = ""
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot metadata with ID: %s\n", forgCmdIn.MetadataID)

	},
}
