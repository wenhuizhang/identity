// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package connect

import (
	"context"
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	vaulttypes "github.com/agntcy/identity/internal/issuer/vault/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type HashicorpFlags struct {
	Address   string
	Token     string
	Namespace string
	VaultName string
}

type HashicorpCommand struct {
	vaultService vault.VaultService
}

func NewCmdHashicorp(vaultService vault.VaultService) *cobra.Command {
	flags := NewHashicorpFlags()

	cmd := &cobra.Command{
		Use:   "hashicorp",
		Short: "Connect to a HashiCorp Vault instance",
		Run: func(cmd *cobra.Command, args []string) {
			c := HashicorpCommand{
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

func NewHashicorpFlags() *HashicorpFlags {
	return &HashicorpFlags{}
}

func (f *HashicorpFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&f.Address,
		"address",
		"a",
		"",
		"The address of the HashiCorp Vault instance",
	)
	cmd.Flags().StringVarP(
		&f.Token,
		"token",
		"t",
		"",
		"The token to authenticate with the HashiCorp Vault instance",
	)
	cmd.Flags().StringVarP(
		&f.Namespace,
		"namespace",
		"n",
		"",
		"The namespace to use in the HashiCorp Vault instance")
	cmd.Flags().StringVarP(
		&f.VaultName,
		"vault-name",
		"v",
		"",
		"Name of the vault",
	)
}

func (cmd *HashicorpCommand) Run(ctx context.Context, flags *HashicorpFlags) error {
	// if the vault address is not set, prompt the user for it interactively
	if flags.Address == "" {
		fmt.Fprintf(os.Stdout, "Address of the HashiCorp Vault instance: ")

		_, err := fmt.Scanln(&flags.Address)
		if err != nil {
			return fmt.Errorf("error reading vault address: %v", err)
		}
	}
	if flags.Address == "" {
		return fmt.Errorf("no vault address provided")
	}

	// if the vault token is not set, prompt the user for it interactively
	if flags.Token == "" {
		fmt.Fprintf(os.Stdout, "Token to authenticate with the HashiCorp Vault instance: ")
		_, err := fmt.Scanln(&flags.Token)
		if err != nil {
			return fmt.Errorf("error reading vault token: %v", err)
		}
	}
	if flags.Token == "" {
		return fmt.Errorf("no vault token provided")
	}

	// if the vault namespace is not set, prompt the user for it interactively
	if flags.Namespace == "" {
		fmt.Fprintf(os.Stdout, "(Optional) Namespace to use in the HashiCorp Vault instance: ")
		_, err := fmt.Scanln(&flags.Namespace)
		// If the user just presses Enter, Namespace will be "" and err will be an "unexpected newline" error.
		// We should allow this and use the empty value.
		if err != nil {
			if err.Error() != "unexpected newline" {
				return fmt.Errorf("error reading vault namespace: %v", err)
			}
		}
	}

	// if the vault name is not set, prompt the user for it interactively
	if flags.VaultName == "" {
		fmt.Fprintf(os.Stdout, "Name of the vault: ")
		_, err := fmt.Scanln(&flags.VaultName)
		if err != nil {
			return fmt.Errorf("error reading vault name: %v", err)
		}
	}
	if flags.VaultName == "" {
		return fmt.Errorf("no vault name provided")
	}

	hashicorpConfig := vaulttypes.VaultHashicorp{
		Address:   flags.Address,
		Token:     flags.Token,
		Namespace: flags.Namespace,
	}

	var config vaulttypes.VaultConfig = &hashicorpConfig

	vault := vaulttypes.Vault{
		Id:     uuid.NewString(),
		Name:   flags.VaultName,
		Type:   vaulttypes.VaultTypeHashicorp,
		Config: config,
	}

	vaultId, err := cmd.vaultService.ConnectVault(&vault)
	if err != nil {
		return fmt.Errorf("error configuring Hashicorp vault: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Successfully configured Hashicorp vault with ID: %s\n", vaultId)

	err = cliCache.SaveCache(
		&cliCache.Cache{
			VaultId: vaultId,
		},
	)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %v", err)
	}

	return nil
}
