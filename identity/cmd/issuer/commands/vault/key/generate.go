// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"context"
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/pkg/joseutil"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type GenerateCommand struct {
	cache        *clicache.Cache
	vaultService vault.VaultService
}

func NewCmdGenerate(
	cache *clicache.Cache,
	vaultService vault.VaultService,
) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate a new cryptographic key for the vault",
		Run: func(cmd *cobra.Command, args []string) {
			c := GenerateCommand{
				cache:        cache,
				vaultService: vaultService,
			}

			err := c.Run(cmd.Context())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
}

func (cmd *GenerateCommand) Run(ctx context.Context) error {
	err := cmd.cache.ValidateForKey()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	vault, err := cmd.vaultService.GetVault(cmd.cache.VaultId)
	if err != nil {
		return fmt.Errorf("error getting vault: %v", err)
	}

	service, err := newKeyService(vault)
	if err != nil {
		return fmt.Errorf("error creating key service: %v", err)
	}

	keyId := uuid.NewString()

	priv, err := joseutil.GenerateJWK("RS256", "sig", keyId)
	if err != nil {
		return fmt.Errorf("error generating JWK: %v", err)
	}

	err = service.SaveKey(ctx, priv.KID, priv)
	if err != nil {
		return fmt.Errorf("error saving key: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Successfully generated key with ID: %s\n", keyId)

	cmd.cache.KeyID = keyId

	err = clicache.SaveCache(cmd.cache)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %v", err)
	}

	return nil
}
