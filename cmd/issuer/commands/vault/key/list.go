// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"context"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	vaultsrv "github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	cache        *cliCache.Cache
	vaultService vaultsrv.VaultService
}

func NewCmdList(
	cache *cliCache.Cache,
	vaultService vaultsrv.VaultService,
) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all keys in the vault",
		Run: func(cmd *cobra.Command, args []string) {
			c := ListCommand{
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

func (cmd *ListCommand) Run(ctx context.Context) error {
	err := cmd.cache.ValidateForKey()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %w", err)
	}

	// get the vault configuration
	vault, err := cmd.vaultService.GetVault(cmd.cache.VaultId)
	if err != nil {
		return fmt.Errorf("error getting vault: %w", err)
	}

	service, err := newKeyService(vault)
	if err != nil {
		return fmt.Errorf("error creating key service: %w", err)
	}

	keys, err := service.ListKeys(ctx)
	if err != nil {
		return fmt.Errorf("error listing keys: %w", err)
	}

	if len(keys) == 0 {
		return fmt.Errorf("no keys found in the vault")
	}

	fmt.Fprintf(os.Stdout, "Keys in vault '%s':\n", vault.Name)

	for _, key := range keys {
		fmt.Fprintf(os.Stdout, "- %s\n", key)
	}

	return nil
}
