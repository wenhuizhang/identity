// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

var badgeLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a badge configuration",
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

		// if the badge id is not set, prompt the user for it interactively
		if loadCmdIn.BadgeID == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to load:\n")
			_, err := fmt.Scanln(&loadCmdIn.BadgeID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
				return
			}
		}
		if loadCmdIn.BadgeID == "" {
			fmt.Fprintf(os.Stderr, "No badge ID provided.\n")
			return
		}

		// check the badge id is valid
		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, loadCmdIn.BadgeID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting badge: %v\n", err)
			return
		}
		if badge == nil {
			fmt.Fprintf(os.Stderr, "Badge with ID %s not found\n", loadCmdIn.BadgeID)
			return
		}

		// save the metadata id to the cache
		cache.BadgeId = loadCmdIn.BadgeID
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded badge with ID: %s\n", loadCmdIn.BadgeID)

	},
}
