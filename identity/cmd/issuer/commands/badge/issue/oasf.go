// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issue

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	badge "github.com/agntcy/identity/internal/issuer/badge"
	"github.com/agntcy/identity/internal/issuer/badge/data/filesystem"
	issfs "github.com/agntcy/identity/internal/issuer/issuer/data/filesystem"
	mdfs "github.com/agntcy/identity/internal/issuer/metadata/data/filesystem"
	"github.com/agntcy/identity/internal/issuer/vault"
	vfs "github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/spf13/cobra"
)

var (
	// setup the command flags
	issueOasfPath string
)

var IssueOasfCmd = &cobra.Command{
	Use:   "oasf",
	Short: "Issue a badge based on a local OASF file",
	Run: func(cmd *cobra.Command, args []string) {
		// setup the badge service
		badgeFilesystemRepository := filesystem.NewBadgeFilesystemRepository()
		issuerRepository := issfs.NewIssuerFilesystemRepository()
		mdRepository := mdfs.NewMetadataFilesystemRepository()
		oidcAuth := oidc.NewAuthenticator()
		nodeClientPrv := nodeapi.NewNodeClientProvider()
		badgeService := badge.NewBadgeService(
			badgeFilesystemRepository,
			mdRepository,
			issuerRepository,
			oidcAuth,
			nodeClientPrv,
		)
		vaultRepository := vfs.NewVaultFilesystemRepository()
		vaultSrv := vault.NewVaultService(vaultRepository)

		// load the cache to get the vault, issuer and metadata ids
		cache, err := cliCache.LoadCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading local configuration: %v\n", err)
			return
		}
		err = cache.ValidateForBadge()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating local configuration: %v\n", err)
			return
		}

		// if the file path is not set, prompt the user for it interactively
		if issueOasfPath == "" {
			fmt.Fprintf(os.Stderr, "Full file path to the OASF you want to sign in the badge: \n")
			_, err := fmt.Scanln(&issueOasfPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading OASF path: %v\n", err)
				return
			}
		}
		if issueOasfPath == "" {
			fmt.Fprintf(os.Stderr, "No OASF path provided\n")
			return
		}

		badgeContentData, err := os.ReadFile(issueOasfPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			return
		}

		prvKey, err := vaultSrv.RetrievePrivKey(cmd.Context(), cache.VaultId, cache.KeyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error retreiving public key: %v\n", err)
			return
		}

		badgeId, err := badgeService.IssueBadge(
			cache.VaultId,
			cache.IssuerId,
			cache.MetadataId,
			&vctypes.CredentialContent[vctypes.BadgeClaims]{
				Type:    vctypes.CREDENTIAL_CONTENT_TYPE_AGENT_BADGE,
				Content: vctypes.BadgeClaims{Badge: string(badgeContentData)},
			},
			prvKey,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error issuing badge: %v\n", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Issued badge with ID: %s\n", badgeId)

		// Save the badge ID to the cache
		cache.BadgeId = badgeId
		err = cliCache.SaveCache(cache)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}
	},
}

func init() {
	IssueOasfCmd.Flags().StringVarP(
		&issueOasfPath,
		"oasf-path",
		"o",
		"",
		"The file path to the OASF you want to sign in the badge",
	)
}
