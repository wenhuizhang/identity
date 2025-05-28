// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package metadata

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/metadata"
	issuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type GenerateFlags struct {
	IdpClientID     string
	IdpClientSecret string
	IdpIssuerURL    string
}

type GenerateCommand struct {
	cache           *clicache.Cache
	metadataService metadata.MetadataService
}

func NewCmdGenerate(
	cache *clicache.Cache,
	metadataService metadata.MetadataService,
) *cobra.Command {
	flags := NewGenerateFlags()

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate new metadata for your Agent and MCP Server identities",
		Run: func(cmd *cobra.Command, args []string) {
			c := GenerateCommand{
				cache:           cache,
				metadataService: metadataService,
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

func NewGenerateFlags() *GenerateFlags {
	return &GenerateFlags{}
}

func (f *GenerateFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.IdpClientID, "idp-client-id", "c", "", "IDP Client ID")
	cmd.Flags().StringVarP(&f.IdpClientSecret, "idp-client-secret", "s", "", "IDP Client Secret")
	cmd.Flags().StringVarP(&f.IdpIssuerURL, "idp-issuer-url", "u", "", "IDP Issuer URL")
}

func (cmd *GenerateCommand) Run(ctx context.Context, flags *GenerateFlags) error {
	err := cmd.cache.ValidateForMetadata()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	// if the idp client id is not set, prompt the user for it interactively
	if flags.IdpClientID == "" {
		fmt.Fprintf(os.Stdout, "IDP Client ID: ")

		_, err := fmt.Scanln(&flags.IdpClientID)
		if err != nil {
			return fmt.Errorf("error reading IDP Client ID: %v", err)
		}
	}

	if flags.IdpClientID == "" {
		return fmt.Errorf("no IDP Client ID provided")
	}

	// if the idp client secret is not set, prompt the user for it interactively
	if flags.IdpClientSecret == "" {
		fmt.Fprintf(os.Stdout, "IDP Client Secret: ")

		_, err := fmt.Scanln(&flags.IdpClientSecret)
		if err != nil {
			return fmt.Errorf("error reading IDP Client Secret: %v", err)
		}
	}

	if flags.IdpClientSecret == "" {
		return fmt.Errorf("no IDP Client Secret provided")
	}

	// if the idp issuer url is not set, prompt the user for it interactively
	if flags.IdpIssuerURL == "" {
		fmt.Fprintf(os.Stdout, "IDP Issuer URL: ")

		_, err := fmt.Scanln(&flags.IdpIssuerURL)
		if err != nil {
			return fmt.Errorf("error reading IDP Issuer URL: %v", err)
		}
	}

	if flags.IdpIssuerURL == "" {
		return fmt.Errorf("no IDP Issuer URL provided")
	}

	idpConfig := issuerTypes.IdpConfig{
		ClientId:     flags.IdpClientID,
		ClientSecret: flags.IdpClientSecret,
		IssuerUrl:    flags.IdpIssuerURL,
	}

	metadataId, err := cmd.metadataService.GenerateMetadata(
		ctx,
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		&idpConfig,
	)
	if err != nil {
		return fmt.Errorf("error generating metadata: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Generated metadata with ID: %s\n", metadataId)

	// Update the cache with the new metadata ID
	cmd.cache.MetadataId = metadataId

	err = clicache.SaveCache(cmd.cache)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %v", err)
	}

	return nil
}
