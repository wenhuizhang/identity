// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"encoding/json"
	"fmt"
	"os"

	issuer "github.com/agntcy/identity/internal/issuer/issuer"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/spf13/cobra"
)

var IssuerCmd = &cobra.Command{
	Use:   "issuer",
	Short: "Setup your issuer environment, including your vault, identity provider, and identity network",
	Long: `
The setup command is used to configure your local environment for the Identity CLI tool. With it you can:

- (register) Register with an identity provider, such as DUO or Okta, to manage your Agent and MCP identities
- (list) List your existing issuer configurations
- (show) Show details of an issuer configuration
- (forget) Forget an issuer configuration
`,
}

//nolint:mnd // Allow magic number for args
var issuerRegisterCmd = &cobra.Command{
	Use:   "register [identity_node_address] [idp_client_id] [idp_client_secret] [idp_issuer_url]",
	Short: "Register as an Issuer",
	Long:  "Register as an Issuer with an Identity Network using the provided client ID, client secret, and issuer URL.",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		identityNodeAddress := args[0]
		clientID := args[1]
		clientSecret := args[2]
		issuerURL := args[3]

		config := issuerTypes.IdpConfig{
			ClientId:     clientID,
			ClientSecret: clientSecret,
			IssuerUrl:    issuerURL,
		}

		_, err := issuer.RegisterIssuer(identityNodeAddress, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error registering as an Issuer: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nRegistered as an Issuer with Identity Network node at %s\n", identityNodeAddress)
	},
}

var issuerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your existing issuer configurations",
	Run: func(cmd *cobra.Command, args []string) {

		issuers, err := issuer.ListIssuerIds()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing issuers: %v\n", err)
			return
		}
		if len(issuers) == 0 {
			fmt.Fprintf(os.Stdout, "No issuers found.\n")
			return
		}
		fmt.Fprintf(os.Stdout, "Existing issuers:\n")
		for _, issuer := range issuers {
			fmt.Fprintf(os.Stdout, "- %s\n", issuer)
		}
	},
}
var issuerShowCmd = &cobra.Command{
	Use:   "show [issuer_id]",
	Short: "Show details of an issuer configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		issuerId := args[0]
		issuer, err := issuer.GetIssuer(issuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting issuer: %v\n", err)
			return
		}
		if issuer == nil {
			fmt.Fprintf(os.Stdout, "No issuer found with ID: %s\n", issuerId)
			return
		}

		issuerJSON, err := json.MarshalIndent(issuer, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling metadata to JSON: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "%s\n", string(issuerJSON))
	},
}

var issuerForgetCmd = &cobra.Command{
	Use:   "forget [issuer_id]",
	Short: "Forget an issuer configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		issuerId := args[0]
		err := issuer.ForgetIssuer(issuerId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error forgetting issuer: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stdout, "Forgot issuer with ID: %s\n", issuerId)
	},
}

func init() {
	IssuerCmd.AddCommand(issuerRegisterCmd)
	IssuerCmd.AddCommand(issuerListCmd)
	IssuerCmd.AddCommand(issuerShowCmd)
	IssuerCmd.AddCommand(issuerForgetCmd)
}
