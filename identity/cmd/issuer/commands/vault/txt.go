// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/core/keystore"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	"github.com/agntcy/identity/internal/pkg/jwkutil"
	"github.com/spf13/cobra"
)

var TxtCmd = &cobra.Command{
	Use:   "txt [output_path]",
	Short: "Connect to .txt file",
	Long:  "Connect to .txt file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultFilesystemRepository := filesystem.NewVaultFilesystemRepository()
		vaultService := vault.NewVaultService(vaultFilesystemRepository)

		filePath := args[0]

		fileStorageConfig := keystore.FileStorageConfig{
			FilePath: filePath,
		}

		service, err := keystore.NewKeyService(keystore.FileStorage, fileStorageConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating key service: %v\n", err)
			return
		}

		priv, err := jwkutil.GenerateJWK("RS256", "sig", "test-rsa")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JWK: %v\n", err)
			return
		}

		ctx := context.Background()
		err = service.SaveKey(ctx, priv.KID, priv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving key: %v\n", err)
			return
		}

		txtConfig := internalIssuerTypes.VaultTxt{
			FilePath: filePath,
		}

		var config internalIssuerTypes.VaultConfig = &txtConfig

		vault, err := vaultService.ConnectVault(internalIssuerTypes.VaultTypeTxt, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error connecting to vault: %v\n", err)
			return
		}

		cmd.Printf("Successfully connected to vault: %s\n", vault.Id)

		err = cliCache.SaveCache(
			&cliCache.Cache{
				VaultId: vault.Id,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving cache: %v\n", err)
			return
		}
	},
}

func init() {}
