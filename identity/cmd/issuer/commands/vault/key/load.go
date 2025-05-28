// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"context"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

type LoadFlags struct {
	KeyID string
}

type LoadCommand struct {
	cache        *cliCache.Cache
	vaultService vault.VaultService
}

func NewCmdLoad(
	cache *cliCache.Cache,
	vaultService vault.VaultService,
) *cobra.Command {
	flags := NewLoadFlags()

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load a key from the vault",
		Run: func(cmd *cobra.Command, args []string) {
			c := LoadCommand{
				cache:        cache,
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
	cmd.Flags().StringVarP(&f.KeyID, "key-id", "k", "", "The ID of the key to load")
}

func (cmd *LoadCommand) Run(ctx context.Context, flags *LoadFlags) error {
	err := cmd.cache.ValidateForKey()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	// if the key id is not set, prompt the user for it interactively
	if flags.KeyID == "" {
		fmt.Fprintf(os.Stdout, "Key ID: ")

		_, err := fmt.Scanln(&flags.KeyID)
		if err != nil {
			return fmt.Errorf("error reading key ID: %v", err)
		}
	}

	if flags.KeyID == "" {
		return fmt.Errorf("no key ID provided")
	}

	// get the vault configuration
	vault, err := cmd.vaultService.GetVault(cmd.cache.VaultId)
	if err != nil {
		return fmt.Errorf("error getting vault: %v", err)
	}

	service, err := newKeyService(vault)
	if err != nil {
		return fmt.Errorf("error creating key service: %v", err)
	}

	// check if the key exists
	_, err = service.RetrievePrivKey(ctx, flags.KeyID)
	if err != nil {
		return fmt.Errorf("error retrieving private key: %v", err)
	}

	_, err = service.RetrievePubKey(ctx, flags.KeyID)
	if err != nil {
		return fmt.Errorf("error retrieving public key: %v", err)
	}

	// save the key id to the cache
	cmd.cache.KeyID = flags.KeyID

	err = cliCache.SaveCache(cmd.cache)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Loaded Key with ID: %s\n", flags.KeyID)

	return nil
}
