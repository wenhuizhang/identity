// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	issuerSetup "github.com/agntcy/identity/internal/issuer/setup"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage your connection to an Identity Network node",
	Long: `
The network command is used to manage your connection to an Identity Network node. With it you can:

- (config) Configure the connection to an Identity Network node
- (test) Test the connection to an Identity Network node
- (forget) Forget the connection to an Identity Network node
`,
}

var networkConfigCmd = &cobra.Command{
	Use:   "config [identity_node_address]",
	Short: "Configure the connection to an Identity Network node",
	Long:  "Configure the connection to an Identity Network node using the provided identity node address.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		identityNodeAddress := args[0]

		config := issuerTypes.IdentityNodeConfig{
			IdentityNodeAddress: identityNodeAddress,
		}

		configPath, err := issuerSetup.ConfigureNetwork(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error configuring Identity Network node: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSaved Identity Node configuration to %s\n\n", configPath)
	},
}

var networkTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the connection to an Identity Network node",
	Run: func(cmd *cobra.Command, args []string) {

		err := issuerSetup.TestNetworkConnection()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError testing Identity Network node connection: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "\nSuccessfully connected to Identity Network node\n\n")
	},
}

var networkForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the connection to an Identity Network node",
	Run: func(cmd *cobra.Command, args []string) {

		configPath, err := issuerSetup.ForgetNetworkConnection()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError forgetting Identity Network node connection: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSuccessfully removed Identity Network node configuration from %s\n\n", configPath)
	},
}

func init() {
	NetworkCmd.AddCommand(networkConfigCmd)
	NetworkCmd.AddCommand(networkTestCmd)
	NetworkCmd.AddCommand(networkForgetCmd)
}
