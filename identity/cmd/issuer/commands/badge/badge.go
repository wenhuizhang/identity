// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package badge

import (
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/data/filesystem"
	"github.com/spf13/cobra"
)

var (
	// setup the badge service
	badgeFilesystemRepository = filesystem.NewBadgeFilesystemRepository()
	badgeService              = badge.NewBadgeService(badgeFilesystemRepository)

	// setup the command flags
	issueFilePath  string
	publishBadgeId string
	showBadgeId    string
	forgetBadgeId  string
	loadBadgeId    string
)

var BadgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Issue and publish badges for your Agent and MCP Server identities",
	Long: `
The badge command is used to issue and publish badges for your Agent and MCP Server identities. With it you can:

- (issue) Issue a new badge for the issuer and metadata
- (publish) Publish the chosen badge
- (list) List your existing badges for the current issuer and metadata
- (show) Show details of the chosen badge
- (forget) Forget a specific badge
`,
}

var badgeIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue a new badge for the current metadata",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in the local configuration. Please load and existing metadata or generate a new metadata first.\n")
			return
		}

		// if the file path is not set, prompt the user for it interactively
		if issueFilePath == "" {
			fmt.Fprintf(os.Stderr, "Full file path to the data you want to sign in the badge: \n")
			fmt.Scanln(&issueFilePath)
		}
		if issueFilePath == "" {
			fmt.Fprintf(os.Stderr, "No file path provided\n")
			return
		}

		badgeId, err := badgeService.IssueBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, issueFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error issuing badge: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Issued badge with ID: %s\n", badgeId)

		// Save the badge ID to the cache
		cache.BadgeId = badgeId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
	},
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
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in the local configuration. Please load and existing metadata or generate a new metadata first.\n")
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
				fmt.Scanln(&publishBadgeId)
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
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in the local configuration. Please load and existing metadata or generate a new metadata first.\n")
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
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in the local configuration. Please load and existing metadata or generate a new metadata first.\n")
			return
		}

		// if the badge id is not set, prompt the user for it interactively
		if showBadgeId == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to show:\n")
			fmt.Scanln(&showBadgeId)
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
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in the local configuration. Please load and existing metadata or generate a new metadata first.\n")
			return
		}

		// if the badge id is not set, prompt the user for it interactively
		if forgetBadgeId == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to forget:\n")
			fmt.Scanln(&forgetBadgeId)
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
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in the local configuration. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in the local configuration. Please load and existing metadata or generate a new metadata first.\n")
			return
		}

		// if the badge id is not set, prompt the user for it interactively
		if loadBadgeId == "" {
			fmt.Fprintf(os.Stderr, "Badge ID to load:\n")
			fmt.Scanln(&loadBadgeId)
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

	badgeIssueCmd.Flags().StringVarP(&issueFilePath, "file-path", "f", "", "The file path to the data you want to sign in the badge")
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
