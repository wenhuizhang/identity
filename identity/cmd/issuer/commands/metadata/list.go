// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
)

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

		allMetadata, err := metadataService.GetAllMetadata(cache.VaultId, cache.KeyID, cache.IssuerId)
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
			fmt.Fprintf(os.Stdout, "- %s\n", metadata.ID)
		}
	},
}
