// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

//nolint:lll // Allow long lines for CLI
var MetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Issue and publish important metadata for your Agent and MCP Server identities",
	Long: `
The metadata command is used to issue and publish important metadata for your Agent and MCP Server identities. With it you can:

- (issue) Issue and load a new metadata
- (publish) Publish the current metadata
- (list) List all of your existing metadata
- (load) Load an existing metadata
- (show) Show the currently loaded metadata
- (forget) Forget the current metadata
`,
}

var metadataListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your existing metadata",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Listing all of your existing metadata")
	},
}

var metadataShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the currently loaded metadata",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Showing the currently loaded metadata")
	},
}

var metadataLoadCmd = &cobra.Command{
	Use:   "load [metadata_id]",
	Short: "Load an existing metadata <metadata_id>",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Printf("Loading metadata %s\n", args[0])
	},
}

var metadataIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issues and loads a new metadata",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Creating a new metadata")
	},
}

var metadataForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the current metadata",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Forgetting the current metadata")
	},
}

var metadataPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the current metadata",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println(
			"Publishing the current metadata",
		)
	},
}

func init() {
	MetadataCmd.AddCommand(metadataListCmd)
	MetadataCmd.AddCommand(metadataLoadCmd)
	MetadataCmd.AddCommand(metadataShowCmd)
	MetadataCmd.AddCommand(metadataIssueCmd)
	MetadataCmd.AddCommand(metadataForgetCmd)
	MetadataCmd.AddCommand(metadataPublishCmd)
}
