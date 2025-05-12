// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package setup

import (
	"fmt"

	"github.com/spf13/cobra"
)

var IdpCmd = &cobra.Command{
	Use:   "idp",
	Short: "Manage your connection to an Identity Provider, such as DUO or Okta",
	Long: `
The idp command is used to manage your connection to an Identity Provider. With it you can:

- (setup) Setup the connection to an Identity Provider
- (test) Test the connection to an Identity Provider
- (forget) Forget the connection to an Identity Provider
`,
}

var idpConnectCmd = &cobra.Command{
	Use:   "setup [client_id] [client_secret] [issuer_url]",
	Short: "Setup the connection to an Identity Provider",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Setting up connection to an Identity Provider")
	},
}

var idpTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the connection to an Identity Provider",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Testing connection to an Identity Provider")
	},
}

var idpForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the connection to an Identity Provider",
	Run: func(cmd *cobra.Command, args []string) {
		//nolint:forbidigo // Allow print for CLI
		fmt.Println("Forgetting connection to an Identity Provider")
	},
}

func init() {
	IdpCmd.AddCommand(idpConnectCmd)
	IdpCmd.AddCommand(idpTestCmd)
	IdpCmd.AddCommand(idpForgetCmd)
}
