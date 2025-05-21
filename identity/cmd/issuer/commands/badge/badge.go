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

var badgeFilesystemRepository = filesystem.NewBadgeFilesystemRepository()
var badgeService = badge.NewBadgeService(badgeFilesystemRepository)

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
	Use:   "issue [badge_file_path]",
	Short: "Issue a new badge for the current metadata",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in cache. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in cache. Please load and existing metadata or generate a new metadata first.\n")
			return
		}

		badgeFilePath := args[0]
		badgeId, err := badgeService.IssueBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, badgeFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error issuing badge: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Issued badge with ID: %s\n", badgeId)

		// Save the badge ID to the cache
		cache.BadgeId = badgeId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
			return
		}
	},
}

//nolint:mnd // Allow magic number for args
var badgePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the chosen badge",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer, metadata and badge ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in cache. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in cache. Please load and existing metadata or generate a new metadata first.\n")
			return
		}
		if cache.BadgeId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No badge found in cache. Please load and existing badge or issue a new badge first.\n")
			return
		}

		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, cache.BadgeId)
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
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in cache. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in cache. Please load and existing metadata or generate a new metadata first.\n")
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
	Use:   "show [badge_id]",
	Short: "Show details of the chosen badge",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in cache. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in cache. Please load and existing metadata or generate a new metadata first.\n")
			return
		}

		badgeId := args[0]

		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, badgeId)
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
	Use:   "forget [badge_id]",
	Short: "Forget the chosen badge",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in cache. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in cache. Please load and existing metadata or generate a new metadata first.\n")
			return
		}
		badgeId := args[1]

		err = badgeService.ForgetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, badgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting badge: %v\n", err)
			return
		}

		// If the badge was the current badge in the cache, clear the cache of badge id
		if cache.BadgeId == badgeId {
			cache.BadgeId = ""
			err = cliCache.SaveCache(cache)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
				return
			}
		}

		fmt.Fprintf(os.Stdout, "Forgot badge with ID: %s\n", badgeId)
	},
}

var badgeLoadCmd = &cobra.Command{
	Use:   "load [badge_id]",
	Short: "Load a badge configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading cache: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in cache. Please load an existing vault or connect to a new vault first.\n")
			return
		}
		if cache.IssuerId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No issuer found in cache. Please load and existing issuer or register a new issuer first.\n")
			return
		}
		if cache.MetadataId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No metadata found in cache. Please load and existing metadata or generate a new metadata first.\n")
			return
		}
		badgeId := args[1]

		// check the badge id is valid
		badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, badgeId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting badge: %v\n", err)
			return
		}
		if badge == nil {
			fmt.Fprintf(os.Stderr, "Badge with ID %s not found\n", badgeId)
			return
		}

		// save the metadata id to the cache
		cache.BadgeId = badgeId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Loaded badge with ID: %s\n", badgeId)

	},
}

func init() {
	BadgeCmd.AddCommand(badgeIssueCmd)
	BadgeCmd.AddCommand(badgePublishCmd)
	BadgeCmd.AddCommand(badgeListCmd)
	BadgeCmd.AddCommand(badgeShowCmd)
	BadgeCmd.AddCommand(badgeForgetCmd)
	BadgeCmd.AddCommand(badgeLoadCmd)
}
