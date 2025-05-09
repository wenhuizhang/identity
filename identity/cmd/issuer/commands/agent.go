// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var AgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Create and manage your agent identities",
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your existing agents",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Listing all of your existing agents")
	},
}

var agentShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the currently loaded agent",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Showing the currently loaded agent")
	},
}

var agentLoadCmd = &cobra.Command{
	Use:   "load [agent_id]",
	Short: "Load an existing agent <agent_id>",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Printf("Loading agent %s\n", args[0])
	},
}

var agentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and loads a new agent",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Creating a new agent")
	},
}

var agentForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the current agent",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Forgetting the current agent")
	},
}

var agentPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the current agent identity",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println(
			"Publishing the current agent identity",
		)
	},
}

func init() {
	AgentCmd.AddCommand(agentListCmd)
	AgentCmd.AddCommand(agentLoadCmd)
	AgentCmd.AddCommand(agentShowCmd)
	AgentCmd.AddCommand(agentCreateCmd)
	AgentCmd.AddCommand(agentForgetCmd)
	AgentCmd.AddCommand(agentPublishCmd)

	// Add the agent artefact commands to the agent command
	AgentCmd.AddCommand(artefactCmd)
}
