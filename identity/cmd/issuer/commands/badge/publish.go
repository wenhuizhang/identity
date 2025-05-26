// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/spf13/cobra"
)

var badgePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the chosen badge",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer, metadata and badge ids
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
		// if there is a badge id in the cache, use it as the default when prompting
		if pubCmdIn.BadgeID == "" {
			if cache.BadgeId != "" {
				fmt.Fprintf(os.Stderr, "Badge ID to publish (default: %s):\n", cache.BadgeId)
				_, err = fmt.Scanln(&pubCmdIn.BadgeID)

				if err != nil {
					// If the user just presses Enter, pubCmdIn.BadgeId will be "" and err will be an "unexpected newline" error.
					// We should allow this and use the default value.
					if err.Error() != "unexpected newline" {
						fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
						return
					}
				}
				// If the user just presses Enter, pubCmdIn.BadgeId will be "" and we should use the default value from the cache.
				if pubCmdIn.BadgeID == "" {
					pubCmdIn.BadgeID = cache.BadgeId
				}
			} else {
				fmt.Fprintf(os.Stderr, "Badge ID to publish:\n")
				_, err := fmt.Scanln(&pubCmdIn.BadgeID)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
					return
				}
			}
		}

		// if the badge id is still not set, then the cache badge is empty and the user has not provided one
		if pubCmdIn.BadgeID == "" {
			fmt.Fprintf(os.Stderr, "No badge ID provided.\n")
			return
		}

		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, pubCmdIn.BadgeID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting badge: %v\n", err)
			return
		}

		_, err = badgeService.PublishBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, badge)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error publishing badge: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", "Publishing the current badge")
	},
}
