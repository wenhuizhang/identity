// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package key

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

type ShowFlags struct {
	KeyID string
}

type ShowCommand struct {
	cache        *cliCache.Cache
	vaultService vault.VaultService
}

func NewCmdShow(
	cache *cliCache.Cache,
	vaultService vault.VaultService,
) *cobra.Command {
	flags := NewShowFlags()

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show details of a specific key in the vault",
		Run: func(cmd *cobra.Command, args []string) {
			c := ShowCommand{
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

func NewShowFlags() *ShowFlags {
	return &ShowFlags{}
}

func (f *ShowFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.KeyID, "key-id", "k", "", "The ID of the key to show")
}

func (cmd *ShowCommand) Run(ctx context.Context, flags *ShowFlags) error {
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

	publicKey, err := service.RetrievePubKey(ctx, flags.KeyID)
	if err != nil {
		return fmt.Errorf("error retrieving public key: %v", err)
	}

	if publicKey == nil {
		return fmt.Errorf("no public key found for Key ID: %s", flags.KeyID)
	}

	// convert the public key to a string representation
	publicKeyStr, err := json.MarshalIndent(publicKey, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling public key: %v", err)
	}

	privateKey, err := service.RetrievePrivKey(ctx, flags.KeyID)
	if err != nil {
		return fmt.Errorf("error retrieving private key: %v", err)
	}

	// convert the private key to a string representation
	if privateKey == nil {
		return fmt.Errorf("no private key found for Key ID: %s", flags.KeyID)
	}

	privateKeyStr, err := json.MarshalIndent(privateKey, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling private key: %v", err)
	}

	fmt.Fprintf(os.Stdout, "\nKey ID: %s\n", flags.KeyID)
	fmt.Fprintf(os.Stdout, "\nPublic Key: %s\n", publicKeyStr)
	fmt.Fprintf(os.Stdout, "\nPrivate Key: %s\n", privateKeyStr)

	return nil
}
