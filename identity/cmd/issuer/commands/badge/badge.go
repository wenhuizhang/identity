// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/cmd/issuer/commands/badge/issue"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/data/filesystem"
	"github.com/spf13/cobra"
)

var (
	// setup the badge service
	badgeFilesystemRepository = filesystem.NewBadgeFilesystemRepository()
	badgeService              = badge.NewBadgeService(badgeFilesystemRepository)

	// setup the command flags
	publishBadgeId string
	showBadgeId    string
	forgetBadgeId  string
	loadBadgeId    string
)

var BadgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Issue and publish badges for your Agent and MCP Server identities",
	Long: `
The badge command is used to issue and publish badges for your Agent and MCP Server identities.
`,
}

var badgeIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue badges using different data sources",
	Long: `
The issue command is used to create Badges for your Agent and MCP Server identities from various data sources.
`,
}

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
		if publishBadgeId == "" {
			if cache.BadgeId != "" {
				fmt.Fprintf(os.Stderr, "Badge ID to publish (default: %s):\n", cache.BadgeId)
				_, err = fmt.Scanln(&publishBadgeId)

				if err != nil {
					// If the user just presses Enter, publishBadgeId will be "" and err will be an "unexpected newline" error.
					// We should allow this and use the default value.
					if err.Error() != "unexpected newline" {
						fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
						return
					}
				}
				// If the user just presses Enter, publishBadgeId will be "" and we should use the default value from the cache.
				if publishBadgeId == "" {
					publishBadgeId = cache.BadgeId
				}
			} else {
				fmt.Fprintf(os.Stderr, "Badge ID to publish:\n")
				_, err := fmt.Scanln(&publishBadgeId)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
					return
				}
			}
		}

		// if the badge id is still not set, then the cache badge is empty and the user has not provided one
		if publishBadgeId == "" {
			fmt.Fprintf(os.Stderr, "No badge ID provided.\n")
			return
		}

		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, publishBadgeId)
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

var badgeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details of the chosen badge",
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
		if showBadgeId == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to show:\n")
			_, err := fmt.Scanln(&showBadgeId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
				return
			}
		}
		if showBadgeId == "" {
			fmt.Fprintf(os.Stderr, "No badge ID provided.\n")
			return
		}

		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, showBadgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting badge: %v\n", err)
			return
		}
		badgeJSON, err := json.MarshalIndent(badge, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling badge to JSON: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(badgeJSON))
	},
}

var badgeForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the chosen badge",
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
		if forgetBadgeId == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to forget:\n")
			_, err := fmt.Scanln(&forgetBadgeId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
				return
			}
		}
		if forgetBadgeId == "" {
			fmt.Fprintf(os.Stderr, "No badge ID provided.\n")
			return
		}

		err = badgeService.ForgetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, forgetBadgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting badge: %v\n", err)
			return
		}

		// If the badge was the current badge in the cache, clear the cache of badge id
		if cache.BadgeId == forgetBadgeId {
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot badge with ID: %s\n", forgetBadgeId)
	},
}

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
		if loadBadgeId == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to load:\n")
			_, err := fmt.Scanln(&loadBadgeId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading badge ID: %v\n", err)
				return
			}
		}
		if loadBadgeId == "" {
			fmt.Fprintf(os.Stderr, "No badge ID provided.\n")
			return
		}

		// check the badge id is valid
		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, loadBadgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting badge: %v\n", err)
			return
		}
		if badge == nil {
			fmt.Fprintf(os.Stderr, "Badge with ID %s not found\n", loadBadgeId)
			return
		}

		// save the metadata id to the cache
		cache.BadgeId = loadBadgeId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded badge with ID: %s\n", loadBadgeId)

	},
}

func init() {
	badgeIssueCmd.AddCommand(issue.IssueFileCmd)
	badgeIssueCmd.AddCommand(issue.IssueOasfCmd)
	badgeIssueCmd.AddCommand(issue.IssueMcpServerCmd)
	badgeIssueCmd.AddCommand(issue.IssueA2AWellKnownCmd)
	BadgeCmd.AddCommand(badgeIssueCmd)

	badgePublishCmd.Flags().StringVarP(&publishBadgeId, "badge-id", "b", "", "The ID of the badge to publish")
	BadgeCmd.AddCommand(badgePublishCmd)

	BadgeCmd.AddCommand(badgeListCmd)

	badgeShowCmd.Flags().StringVarP(&showBadgeId, "badge-id", "b", "", "The ID of the badge to show")
	BadgeCmd.AddCommand(badgeShowCmd)

	badgeForgetCmd.Flags().StringVarP(&forgetBadgeId, "badge-id", "b", "", "The ID of the badge to forget")
	BadgeCmd.AddCommand(badgeForgetCmd)

	badgeLoadCmd.Flags().StringVarP(&loadBadgeId, "badge-id", "b", "", "The ID of the badge to load")
	BadgeCmd.AddCommand(badgeLoadCmd)
}
