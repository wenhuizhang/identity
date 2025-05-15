// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var BadgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Issue and publish badges for your Agent and MCP Server identities",
	Long: `
The badge command is used to issue and publish badges for your Agent and MCP Server identities. With it you can:

- (issue) Issue and load a new badge
- (publish) Publish the current badge
- (list) List all of your existing badges
- (load) Load an existing badge
- (show) Show the currently loaded badge
- (forget) Forget the current badge
`,
}

var badgeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your existing badges",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Listing all of your existing badges")
	},
}

var badgeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the currently loaded badge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Showing the currently loaded badge")
	},
}

var badgeLoadCmd = &cobra.Command{
	Use:   "load [badge_id]",
	Short: "Load an existing badge <badge_id>",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "Loading badge %s\n", args[0])
	},
}

var badgeIssueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issues and loads a new badge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Issuing and loading a new badge")
	},
}

var badgeForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the current badge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n", "Forgetting the current badge")
	},
}

var badgePublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the current badge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stdout, "%s\n",
			"Publishing the current badge",
		)
	},
}

func init() {
	BadgeCmd.AddCommand(badgeListCmd)
	BadgeCmd.AddCommand(badgeLoadCmd)
	BadgeCmd.AddCommand(badgeShowCmd)
	BadgeCmd.AddCommand(badgeIssueCmd)
	BadgeCmd.AddCommand(badgeForgetCmd)
	BadgeCmd.AddCommand(badgePublishCmd)
}
