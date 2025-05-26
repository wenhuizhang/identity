// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

var badgeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing badges for the current metadata",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForBadge()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		badges, err := badgeService.GetAllBadges(cache.VaultId, cache.IssuerId, cache.MetadataId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing badges: %v\n", err)
			return
		}
		if len(badges) == 0 {
			fmt.Fprintf(os.Stdout, "%s\n", "No badges found")
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Existing badge ids:")
		for _, badge := range badges {
			fmt.Fprintf(os.Stdout, "- %s\n", badge.Id)
		}
	},
}
