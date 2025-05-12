// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"fmt"

	"github.com/spf13/cobra"
)

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage your connection to an Identity Network node",
	Long: `
The network command is used to manage your connection to an Identity Network node. With it you can:

- (setup) Setup the connection to an Identity Network node
- (test) Test the connection to an Identity Network node
- (forget) Forget the connection to an Identity Network node
`,
}

var networkConnectCmd = &cobra.Command{
	Use:   "setup [identity_node_address]",
	Short: "Setup the connection to an Identity Network node",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Setting up connection to an Identity Network node")
	},
}

var networkTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the connection to an Identity Network node",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Testing connection to an Identity Network node")
	},
}

var networkForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the connection to an Identity Network node",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Forgetting connection to an Identity Network node")
	},
}

func init() {
	NetworkCmd.AddCommand(networkConnectCmd)
	NetworkCmd.AddCommand(networkTestCmd)
	NetworkCmd.AddCommand(networkForgetCmd)
}
