// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

type LoadFlags struct {
	VaultID string
}

type LoadCommand struct {
	vaultService vault.VaultService
}

func NewCmdLoad(vaultService vault.VaultService) *cobra.Command {
	flags := NewLoadFlags()

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load a vault configuration",
		Run: func(cmd *cobra.Command, args []string) {
			c := LoadCommand{
				vaultService: vaultService,
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

func NewLoadFlags() *LoadFlags {
	return &LoadFlags{}
}

func (f *LoadFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.VaultID, "vault-id", "v", "", "The ID of the vault to load")
}

func (cmd *LoadCommand) Run(ctx context.Context, flags *LoadFlags) error {
	// if the vault id is not set, prompt the user for it interactively
	if flags.VaultID == "" {
		fmt.Fprintf(os.Stdout, "Vault ID to load:\n")

		_, err := fmt.Scanln(&flags.VaultID)
		if err != nil {
			return fmt.Errorf("error reading vault ID: %v", err)
		}
	}

	if flags.VaultID == "" {
		return fmt.Errorf("no vault ID provided")
	}

	// check the vault id is valid
	vault, err := cmd.vaultService.GetVault(flags.VaultID)
	if err != nil {
		return fmt.Errorf("error getting vault: %v", err)
	}

	if vault == nil {
		fmt.Fprintf(os.Stdout, "No vault found with ID: %s\n", flags.VaultID)
		return nil
	}

	// save the vault id to the cache
	err = cliCache.SaveCache(
		&cliCache.Cache{
			VaultId: vault.Id,
		},
	)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Loaded vault with ID: %s\n", flags.VaultID)

	return nil
}
