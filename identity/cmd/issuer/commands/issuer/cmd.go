// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	issuer "github.com/agntcy/identity/internal/issuer/issuer"
	"github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	"github.com/agntcy/identity/internal/issuer/vault"
	vfr "github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/spf13/cobra"
)

const (
	defaultNodeAddress = "http://localhost:4000"
)

type RegisterCmdInput struct {
	IdentityNodeURL string
	ClientID        string
	ClientSecret    string
	IssuerURL       string
	Organization    string
	SubOrganization string
}

// TODO: remove globals
//
//nolint:godox // To be fixed in the next PR
var (
	// setup the issuer service
	issuerFilesystemRepository = filesystem.NewIssuerFilesystemRepository()
	oidcAuth                   = oidc.NewAuthenticator()
	nodeClientPrv              = nodeapi.NewNodeClientProvider()
	vaultRepository            = vfr.NewVaultFilesystemRepository()
	vaultSrv                   = vault.NewVaultService(vaultRepository)
	issuerService              = issuer.NewIssuerService(issuerFilesystemRepository, oidcAuth, nodeClientPrv)

	// setup the command flags
	registerCmdIn  = &RegisterCmdInput{}
	showIssuerId   string
	forgetIssuerId string
	loadIssuerId   string
)

var IssuerCmd = &cobra.Command{
	Use:   "issuer",
	Short: "Register as an Issuer and manage issuer configurations",
	Long: `
The issuer command is used to register as an Issuer and manage issuer configurations.
`,
}

//nolint:lll // Allow long lines for CLI
func init() {
	issuerRegisterCmd.Flags().StringVarP(&registerCmdIn.IdentityNodeURL, "identity-node-address", "i", "", "Identity node address")
	issuerRegisterCmd.Flags().StringVarP(&registerCmdIn.ClientID, "client-id", "c", "", "IdP client ID")
	issuerRegisterCmd.Flags().StringVarP(&registerCmdIn.ClientSecret, "client-secret", "s", "", "IdP client secret")
	issuerRegisterCmd.Flags().StringVarP(&registerCmdIn.IssuerURL, "issuer-url", "u", "", "IdP issuer URL")
	issuerRegisterCmd.Flags().StringVarP(&registerCmdIn.Organization, "organization", "o", "", "Organization name")
	issuerRegisterCmd.Flags().StringVarP(&registerCmdIn.SubOrganization, "sub-organization", "b", "", "Sub-organization name")
	IssuerCmd.AddCommand(issuerRegisterCmd)

	IssuerCmd.AddCommand(issuerListCmd)

	issuerShowCmd.Flags().StringVarP(&showIssuerId, "issuer-id", "i", "", "The ID of the issuer to show")
	IssuerCmd.AddCommand(issuerShowCmd)

	issuerForgetCmd.Flags().StringVarP(&forgetIssuerId, "issuer-id", "i", "", "The ID of the issuer to forget")
	IssuerCmd.AddCommand(issuerForgetCmd)

	issuerLoadCmd.Flags().StringVarP(&loadIssuerId, "issuer-id", "i", "", "The ID of the issuer to load")
	IssuerCmd.AddCommand(issuerLoadCmd)
}
