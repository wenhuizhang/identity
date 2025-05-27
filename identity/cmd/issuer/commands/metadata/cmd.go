// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"github.com/spf13/cobra"

	issuerFilesystem "github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	"github.com/agntcy/identity/internal/issuer/metadata"
	metadataFilesystem "github.com/agntcy/identity/internal/issuer/metadata/data/filesystem"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/agntcy/identity/internal/pkg/oidc"
)

type GenerateCmdInput struct {
	IdpClientID     string
	IdpClientSecret string
	IdpIssuerURL    string
}

type ShowCmdInput struct {
	MetadataID string
}

type ForgetCmdInput struct {
	MetadataID string
}

type LoadCmdInput struct {
	MetadataID string
}

var (
	// setup the metadata service
	metadataFilesystemRepository = metadataFilesystem.NewMetadataFilesystemRepository()
	issuerFilesystemRepository   = issuerFilesystem.NewIssuerFilesystemRepository()
	oidcAuth                     = oidc.NewAuthenticator()
	nodeClientPrv                = nodeapi.NewNodeClientProvider()
	metadataService              = metadata.NewMetadataService(
		metadataFilesystemRepository,
		issuerFilesystemRepository,
		oidcAuth,
		nodeClientPrv,
	)

	// setup the command flags

	genCmdIn  = &GenerateCmdInput{}
	showCmdIn = &ShowCmdInput{}
	forgCmdIn = &ForgetCmdInput{}
	loadCmdIn = &LoadCmdInput{}
)

var MetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Generate important metadata for your Agent and MCP Server identities",
	Long: `
The metadata command is used to generate important metadata for your Agent and MCP Server identities.
`,
}

func init() {
	metadataGenerateCmd.Flags().StringVarP(&genCmdIn.IdpClientID, "idp-client-id", "c", "", "IDP Client ID")
	metadataGenerateCmd.Flags().StringVarP(&genCmdIn.IdpClientSecret, "idp-client-secret", "s", "", "IDP Client Secret")
	metadataGenerateCmd.Flags().StringVarP(&genCmdIn.IdpIssuerURL, "idp-issuer-url", "u", "", "IDP Issuer URL")
	MetadataCmd.AddCommand(metadataGenerateCmd)

	MetadataCmd.AddCommand(metadataListCmd)

	metadataShowCmd.Flags().StringVarP(&showCmdIn.MetadataID, "metadata-id", "m", "", "The ID of the metadata to show")
	MetadataCmd.AddCommand(metadataShowCmd)

	metadataForgetCmd.Flags().StringVarP(&forgCmdIn.MetadataID, "metadata-id", "m", "", "The ID of the metadata to forget")
	MetadataCmd.AddCommand(metadataForgetCmd)

	metadataLoadCmd.Flags().StringVarP(&loadCmdIn.MetadataID, "metadata-id", "m", "", "The ID of the metadata to load")
	MetadataCmd.AddCommand(metadataLoadCmd)
}
