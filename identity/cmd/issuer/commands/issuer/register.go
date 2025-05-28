// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"
	"fmt"
	"net/url"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	coreissuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	issuer "github.com/agntcy/identity/internal/issuer/issuer"
	issuertypes "github.com/agntcy/identity/internal/issuer/issuer/types"
	idptypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

const (
	defaultNodeAddress = "http://localhost:4000"
)

type RegisterFlags struct {
	IdentityNodeURL string
	ClientID        string
	ClientSecret    string
	IssuerURL       string
	Organization    string
	SubOrganization string
}

type RegisterCommand struct {
	cache         *clicache.Cache
	issuerService issuer.IssuerService
	vaultSrv      vault.VaultService
}

func NewCmdRegister(
	cache *clicache.Cache,
	issuerService issuer.IssuerService,
	vaultSrv vault.VaultService,
) *cobra.Command {
	flags := NewRegisterFlags()

	//nolint:lll // Allow long lines for CLI
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register as an Issuer",
		Long:  "Register with an identity provider, such as DUO or Okta, to manage your Agent and MCP identities",
		Run: func(cmd *cobra.Command, args []string) {
			c := RegisterCommand{
				cache:         cache,
				issuerService: issuerService,
				vaultSrv:      vaultSrv,
			}

			err := c.Run(cmd.Context(), flags)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	flags.AddFlags(cmd)

	return cmd
}

func NewRegisterFlags() *RegisterFlags {
	return &RegisterFlags{}
}

//nolint:lll // Allow long lines for CLI
func (f *RegisterFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.IdentityNodeURL, "identity-node-address", "i", "", "Identity node address")
	cmd.Flags().StringVarP(&f.ClientID, "idp-client-id", "c", "", "IdP client ID")
	cmd.Flags().StringVarP(&f.ClientSecret, "idp-client-secret", "s", "", "IdP client secret")
	cmd.Flags().StringVarP(&f.IssuerURL, "idp-issuer-url", "u", "", "IdP issuer URL")
	cmd.Flags().StringVarP(&f.Organization, "organization", "o", "", "Organization name")
	cmd.Flags().StringVarP(&f.SubOrganization, "sub-organization", "b", "", "Sub-organization name")
}

func (cmd *RegisterCommand) Run(ctx context.Context, flags *RegisterFlags) error {
	err := cmd.cache.ValidateForIssuer()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	// if the identity node address is not set, prompt the user for it interactively
	if flags.IdentityNodeURL == "" {
		fmt.Fprintf(os.Stdout, "Identity node address (default %s): ", defaultNodeAddress)

		_, err = fmt.Scanln(&flags.IdentityNodeURL)
		if err != nil {
			// If the user just presses Enter, registerIdentityNodeAddress will be "" and err will be an "unexpected newline" error.
			// We should allow this and use the default value.
			if err.Error() != "unexpected newline" {
				return fmt.Errorf("error reading identity node address: %v", err)
			}
		}
	}
	// If no address was entered (input was empty or only whitespace), use the default.
	if flags.IdentityNodeURL == "" {
		flags.IdentityNodeURL = defaultNodeAddress
	}

	// if the client ID is not set, prompt the user for it interactively
	if flags.ClientID == "" {
		fmt.Fprintf(os.Stdout, "IdP client ID: ")

		_, err = fmt.Scanln(&flags.ClientID)
		if err != nil {
			return fmt.Errorf("error reading IdP client ID: %v", err)
		}

		if flags.ClientID == "" {
			return fmt.Errorf("idp Client ID cannot be empty")
		}
	}

	// if the client secret is not set, prompt the user for it interactively
	if flags.ClientSecret == "" {
		fmt.Fprintf(os.Stdout, "IdP client secret: ")
		_, err = fmt.Scanln(&flags.ClientSecret)
		if err != nil {
			return fmt.Errorf("error reading IdP client secret: %v", err)
		}
		if flags.ClientSecret == "" {
			return fmt.Errorf("idp client secret cannot be empty")
		}
	}

	// if the issuer URL is not set, prompt the user for it interactively
	if flags.IssuerURL == "" {
		fmt.Fprintf(os.Stdout, "IdP issuer URL: ")
		_, err = fmt.Scanln(&flags.IssuerURL)

		if err != nil {
			return fmt.Errorf("error reading IdP issuer URL: %v", err)
		}

		if flags.IssuerURL == "" {
			return fmt.Errorf("idp issuer URL cannot be empty")
		}
	}

	// if the organization is not set, prompt the user for it interactively
	if flags.Organization == "" {
		fmt.Fprintf(os.Stdout, "Organization name: ")
		_, err = fmt.Scanln(&flags.Organization)
		if err != nil {
			return fmt.Errorf("error reading organization name: %v", err)
		}

		if flags.Organization == "" {
			return fmt.Errorf("organization name cannot be empty")
		}
	}

	// if the sub-organization is not set, prompt the user for it interactively
	if flags.SubOrganization == "" {
		fmt.Fprintf(os.Stdout, "Sub-organization name (default %s): ", flags.Organization)

		_, err = fmt.Scanln(&flags.SubOrganization)
		if err != nil {
			// If the user just presses Enter, registerSubOrganization will be "" and err will be an "unexpected newline" error.
			// We should allow this and use the default value.
			if err.Error() != "unexpected newline" {
				return fmt.Errorf("error reading sub-organization name: %v", err)
			}

			if flags.SubOrganization == "" {
				flags.SubOrganization = flags.Organization
			}
		}
	}

	idpConfig := idptypes.IdpConfig{
		ClientId:     flags.ClientID,
		ClientSecret: flags.ClientSecret,
		IssuerUrl:    flags.IssuerURL,
	}

	// extract the root url from the issuer URL as the common name
	issuerUrl, err := url.Parse(idpConfig.IssuerUrl)
	if err != nil {
		return fmt.Errorf("error parsing issuer URL: %v", err)
	}

	commonName := issuerUrl.Hostname()
	if commonName == "" {
		return fmt.Errorf("error extracting common name from issuer URL: %v", err)
	}

	pubKey, err := cmd.vaultSrv.RetrievePubKey(ctx, cmd.cache.VaultId, cmd.cache.KeyID)
	if err != nil {
		return fmt.Errorf("error retreiving public key: %v", err)
	}

	coreIssuer := coreissuertypes.Issuer{
		Organization:    flags.Organization,
		SubOrganization: flags.SubOrganization,
		CommonName:      commonName,
		PublicKey:       pubKey,
	}

	issuer := issuertypes.Issuer{
		Issuer:          coreIssuer,
		ID:              idpConfig.ClientId,
		IdentityNodeURL: flags.IdentityNodeURL,
		IdpConfig:       &idpConfig,
	}

	issuerId, err := cmd.issuerService.RegisterIssuer(
		ctx,
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		&issuer,
	)
	if err != nil {
		return fmt.Errorf("error registering as an Issuer: %v", err)
	}

	fmt.Fprintf(
		os.Stdout,
		"\nSuccessfully registered as an Issuer with:\n- ID: %s\n- Common Name: %s\n",
		issuerId,
		commonName,
	)

	// Update the cache with the new issuer ID
	cmd.cache.IssuerId = issuerId

	err = clicache.SaveCache(cmd.cache)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %v", err)
	}

	return nil
}
