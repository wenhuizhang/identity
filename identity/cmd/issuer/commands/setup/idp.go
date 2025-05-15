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

var IdpCmd = &cobra.Command{
	Use:   "idp",
	Short: "Manage your connection to an Identity Provider, such as DUO or Okta",
	Long: `
The idp command is used to manage your connection to an Identity Provider. With it you can:

- (config) Setup the connection to an Identity Provider
- (test) Test the connection to an Identity Provider
- (forget) Forget the connection to an Identity Provider
`,
}

//nolint:mnd // Allow magic number 3 for args
var idpConfigCmd = &cobra.Command{
	Use:   "config [client_id] [client_secret] [issuer_url]",
	Short: "Setup the connection to an Identity Provider",
	Long:  "Setup the connection to an Identity Provider using the provided client ID, client secret, and issuer URL.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		clientID := args[0]
		clientSecret := args[1]
		issuerURL := args[2]

		config := issuerTypes.IdpConfig{
			ClientId:     clientID,
			ClientSecret: clientSecret,
			IssuerUrl:    issuerURL,
		}

		configPath, err := issuerSetup.ConfigureIdp(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error configuring Identity Provider: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSaved Identity Provider configuration to %s\n\n", configPath)
	},
}

var idpTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the connection to an Identity Provider",
	Run: func(cmd *cobra.Command, args []string) {

		token, err := issuerSetup.TestIdpConnection()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError testing Identity Provider connection: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSuccessfully connected to Identity Provider and received token: %s\n\n", token.AccessToken)
	},
}

var idpForgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Forget the connection to an Identity Provider",
	Run: func(cmd *cobra.Command, args []string) {

		configPath, err := issuerSetup.ForgetIdpConnection()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError forgetting Identity Provider configuration: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nDeleted Identity Provider configuration at %s\n\n", configPath)
	},
}

func init() {
	IdpCmd.AddCommand(idpConfigCmd)
	IdpCmd.AddCommand(idpTestCmd)
	IdpCmd.AddCommand(idpForgetCmd)
}
