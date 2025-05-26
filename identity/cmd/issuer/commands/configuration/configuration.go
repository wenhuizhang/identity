// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package configuration

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	bsvc "github.com/agntcy/identity/internal/issuer/badge"
	bfr "github.com/agntcy/identity/internal/issuer/badge/data/filesystem"
	isvc "github.com/agntcy/identity/internal/issuer/issuer"
	ifr "github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	msvc "github.com/agntcy/identity/internal/issuer/metadata"
	mfr "github.com/agntcy/identity/internal/issuer/metadata/data/filesystem"
	vsvc "github.com/agntcy/identity/internal/issuer/vault"
	vfr "github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/spf13/cobra"
)

var (
	// setup the vault service
	vaultFilesystemRepository = vfr.NewVaultFilesystemRepository()
	vaultService              = vsvc.NewVaultService(vaultFilesystemRepository)

	// setup the issuer service
	issuerFilesystemRepository = ifr.NewIssuerFilesystemRepository()
	oidcAuth                   = oidc.NewAuthenticator()
	nodeClientPrv              = nodeapi.NewNodeClientProvider()
	issuerService              = isvc.NewIssuerService(issuerFilesystemRepository, oidcAuth, nodeClientPrv)

	// setup the metadata service
	metadataFilesystemRepository = mfr.NewMetadataFilesystemRepository()
	metadataService              = msvc.NewMetadataService(
		metadataFilesystemRepository,
		issuerFilesystemRepository,
		oidcAuth,
		nodeClientPrv,
	)

	// setup the badge service
	badgeFilesystemRepository = bfr.NewBadgeFilesystemRepository()
	badgeService              = bsvc.NewBadgeService(
		badgeFilesystemRepository,
		metadataFilesystemRepository,
		issuerFilesystemRepository,
		oidcAuth,
		nodeClientPrv,
	)
)

var ConfigurationCmd = &cobra.Command{
	Use:   "config",
	Short: "Display the local configuration context",
	Run: func(cmd *cobra.Command, args []string) {

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		if cache == nil || cache.VaultId == "" {
			fmt.Fprintf(
				os.Stderr,
				"No vault found in the local configuration. Please load an existing vault or connect to a new vault first.\n")
			return
		} else {

			vault, err := vaultService.GetVault(cache.VaultId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading vault: %v\n", err)
				return
			}
			if vault == nil {
				fmt.Fprintf(os.Stderr, "No vault found with ID: %s\n", cache.VaultId)
				return
			}
			fmt.Fprintf(os.Stdout, "\nCurrent Identity CLI configuration context:\n")
			fmt.Fprintf(os.Stdout, "- Vault: %s (%s vault), id: %s\n", vault.Name, vault.Type, vault.Id)
		}

		if cache.IssuerId != "" {

			issuer, err := issuerService.GetIssuer(cache.VaultId, cache.IssuerId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading issuer: %v\n", err)
				return
			}
			if issuer == nil {
				fmt.Fprintf(os.Stderr, "No issuer found with ID: %s\n", cache.IssuerId)
				return
			}
			fmt.Fprintf(os.Stdout, "- Issuer: %s\n", issuer.ID)

		}

		if cache.MetadataId != "" {
			metadata, err := metadataService.GetMetadata(cache.VaultId, cache.IssuerId, cache.MetadataId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading metadata: %v\n", err)
				return
			}
			if metadata == nil {
				fmt.Fprintf(os.Stderr, "No metadata found with ID: %s\n", cache.MetadataId)
				return
			}
			fmt.Fprintf(os.Stdout, "- Metadata: %s\n", metadata.ID)

		}

		if cache.BadgeId != "" {
			badge, err := badgeService.GetBadge(cache.VaultId, cache.IssuerId, cache.MetadataId, cache.BadgeId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading badge: %v\n", err)
				return
			}
			if badge == nil {
				fmt.Fprintf(os.Stderr, "No badge found with ID: %s\n", cache.BadgeId)
				return
			}
			fmt.Fprintf(os.Stdout, "- Badge: %s\n", badge.Id)

		}

		fmt.Fprintf(os.Stdout, "\n")
	},
}
