// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"fmt"
	"net/url"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	coreissuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	issuertypes "github.com/agntcy/identity/internal/issuer/issuer/types"
	idptypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/spf13/cobra"
)

//nolint:lll // Allow long lines for CLI
var issuerRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register as an Issuer",
	Long:  "Register with an identity provider, such as DUO or Okta, to manage your Agent and MCP identities",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault id
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}

		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(os.Stderr, "No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		}

		// if the identity node address is not set, prompt the user for it interactively
		if registerCmdIn.IdentityNodeURL == "" {
			fmt.Fprintf(os.Stdout, "Identity node address (default %s): ", defaultNodeAddress)

			_, err = fmt.Scanln(&registerCmdIn.IdentityNodeURL)
			if err != nil {
				// If the user just presses Enter, registerIdentityNodeAddress will be "" and err will be an "unexpected newline" error.
				// We should allow this and use the default value.
				if err.Error() != "unexpected newline" {
					fmt.Fprintf(os.Stderr, "Error reading identity node address: %v\n", err)
					return
				}
			}
		}
		// If no address was entered (input was empty or only whitespace), use the default.
		if registerCmdIn.IdentityNodeURL == "" {
			registerCmdIn.IdentityNodeURL = defaultNodeAddress
		}

		// if the client ID is not set, prompt the user for it interactively
		if registerCmdIn.ClientID == "" {
			fmt.Fprintf(os.Stdout, "IdP client ID: ")
			_, err = fmt.Scanln(&registerCmdIn.ClientID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IdP client ID: %v\n", err)
				return
			}
			if registerCmdIn.ClientID == "" {
				fmt.Fprintf(os.Stderr, "IdP Client ID cannot be empty.\n")
				return
			}
		}

		// if the client secret is not set, prompt the user for it interactively
		if registerCmdIn.ClientSecret == "" {
			fmt.Fprintf(os.Stdout, "IdP client secret: ")
			_, err = fmt.Scanln(&registerCmdIn.ClientSecret)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IdP client secret: %v\n", err)
				return
			}
			if registerCmdIn.ClientSecret == "" {
				fmt.Fprintf(os.Stderr, "IdP Client secret cannot be empty.\n")
				return
			}
		}

		// if the issuer URL is not set, prompt the user for it interactively
		if registerCmdIn.IssuerURL == "" {
			fmt.Fprintf(os.Stdout, "IdP issuer URL: ")
			_, err = fmt.Scanln(&registerCmdIn.IssuerURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading IdP issuer URL: %v\n", err)
				return
			}
			if registerCmdIn.IssuerURL == "" {
				fmt.Fprintf(os.Stderr, "IdP issuer URL cannot be empty.\n")
				return
			}
		}

		// if the organization is not set, prompt the user for it interactively
		if registerCmdIn.Organization == "" {
			fmt.Fprintf(os.Stdout, "Organization name: ")
			_, err = fmt.Scanln(&registerCmdIn.Organization)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading organization name: %v\n", err)
				return
			}
			if registerCmdIn.Organization == "" {
				fmt.Fprintf(os.Stderr, "Organization name cannot be empty.\n")
				return
			}
		}

		// if the sub-organization is not set, prompt the user for it interactively
		if registerCmdIn.SubOrganization == "" {
			fmt.Fprintf(os.Stdout, "Sub-organization name (default %s): ", registerCmdIn.Organization)
			_, err = fmt.Scanln(&registerCmdIn.SubOrganization)
			if err != nil {
				// If the user just presses Enter, registerSubOrganization will be "" and err will be an "unexpected newline" error.
				// We should allow this and use the default value.
				if err.Error() != "unexpected newline" {
					fmt.Fprintf(os.Stderr, "Error reading sub-organization name: %v\n", err)
					return
				}
				if registerCmdIn.SubOrganization == "" {
					registerCmdIn.SubOrganization = registerCmdIn.Organization
				}
			}
		}

		idpConfig := idptypes.IdpConfig{
			ClientId:     registerCmdIn.ClientID,
			ClientSecret: registerCmdIn.ClientSecret,
			IssuerUrl:    registerCmdIn.IssuerURL,
		}

		// extract the root url from the issuer URL as the common name
		issuerUrl, err := url.Parse(idpConfig.IssuerUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing issuer URL: %v\n", err)
			return
		}
		commonName := issuerUrl.Hostname()
		if commonName == "" {
			fmt.Fprintf(os.Stderr, "Error extracting common name from issuer URL: %v\n", err)
			return
		}

		pubKey, err := vaultSrv.RetrievePubKey(cmd.Context(), cache.VaultId, cache.KeyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retreiving public key: %v\n", err)
			return
		}

		coreIssuer := coreissuertypes.Issuer{
			Organization:    registerCmdIn.Organization,
			SubOrganization: registerCmdIn.SubOrganization,
			CommonName:      commonName,
			PublicKey:       pubKey,
		}

		issuer := issuertypes.Issuer{
			Issuer:          coreIssuer,
			ID:              idpConfig.ClientId,
			IdentityNodeURL: registerCmdIn.IdentityNodeURL,
			IdpConfig:       &idpConfig,
		}

		issuerId, err := issuerService.RegisterIssuer(cmd.Context(), cache.VaultId, &issuer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error registering as an Issuer: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "\nSuccessfully registered as an Issuer with ID: %s\n", issuerId)

		// Update the cache with the new issuer ID
		cache.IssuerId = issuerId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
	},
}
