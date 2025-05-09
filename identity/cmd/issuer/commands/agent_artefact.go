// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var artefactCmd = &cobra.Command{
	Use:   "artefact",
	Short: "Create and manage your agent artefacts",
}

var artefactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your existing agent artefacts",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Listing all of your existing agent artefacts")
	},
}

var artefactLoadCmd = &cobra.Command{
	Use:   "load [artefact_id]",
	Short: "Load an existing agent artefact <artefact_id>",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Printf("Loading agent artefact %s\n", args[0])
	},
}

var artefactShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the currently loaded agent artefact",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Showing the currently loaded agent artefact")
	},
}

var artefactCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and loads a new agent artefact",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Creating a new agent artefact")
	},
}

var artefactForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the current agent artefact",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Forgetting the current agent artefact")
	},
}

var artefactPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the current agent artefact identity",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Publishing the current agent artefact identity")
	},
}

func init() {
	artefactCmd.AddCommand(artefactListCmd)
	artefactCmd.AddCommand(artefactLoadCmd)
	artefactCmd.AddCommand(artefactShowCmd)
	artefactCmd.AddCommand(artefactCreateCmd)
	artefactCmd.AddCommand(artefactForgetCmd)
	artefactCmd.AddCommand(artefactPublishCmd)

	// Add the agent artefact version commands to the agent artefact command
	artefactCmd.AddCommand(versionCmd)
}
