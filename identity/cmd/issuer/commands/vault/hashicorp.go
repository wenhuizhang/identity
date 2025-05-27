// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"fmt"
	"os"

	cliCache "github.com/agntcy/identity/cmd/issuer/cache"
	"github.com/agntcy/identity/internal/issuer/vault"
	"github.com/agntcy/identity/internal/issuer/vault/data/filesystem"
	vaulttypes "github.com/agntcy/identity/internal/issuer/vault/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var HashicorpCmd = &cobra.Command{
	Use:   "hashicorp",
	Short: "Connect to a HashiCorp Vault instance",
	Run: func(cmd *cobra.Command, args []string) {

		// if the vault address is not set, prompt the user for it interactively
		if hashicorpCmdIn.Address == "" {
			fmt.Fprintf(os.Stderr, "Address of the HashiCorp Vault instance: ")
			_, err := fmt.Scanln(&hashicorpCmdIn.Address)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault address: %v\n", err)
				return
			}
		}
		if hashicorpCmdIn.Address == "" {
			fmt.Fprintf(os.Stderr, "No vault address provided\n")
			return
		}

		// if the vault token is not set, prompt the user for it interactively
		if hashicorpCmdIn.Token == "" {
			fmt.Fprintf(os.Stderr, "Token to authenticate with the HashiCorp Vault instance: ")
			_, err := fmt.Scanln(&hashicorpCmdIn.Token)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault token: %v\n", err)
				return
			}
		}
		if hashicorpCmdIn.Token == "" {
			fmt.Fprintf(os.Stderr, "No vault token provided\n")
			return
		}

		// if the vault namespace is not set, prompt the user for it interactively
		if hashicorpCmdIn.Namespace == "" {
			fmt.Fprintf(os.Stderr, "(Optional) Namespace to use in the HashiCorp Vault instance: ")
			_, err := fmt.Scanln(&hashicorpCmdIn.Namespace)
			// If the user just presses Enter, Namespace will be "" and err will be an "unexpected newline" error.
			// We should allow this and use the empty value.
			if err != nil {
				if err.Error() != "unexpected newline" {
					fmt.Fprintf(os.Stderr, "Error reading vault namespace: %v\n", err)
					return
				}
			}
		}

		// if the vault name is not set, prompt the user for it interactively
		if hashicorpCmdIn.VaultName == "" {
			fmt.Fprintf(os.Stderr, "Name of the vault: ")
			_, err := fmt.Scanln(&hashicorpCmdIn.VaultName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading vault name: %v\n", err)
				return
			}
		}
		if hashicorpCmdIn.VaultName == "" {
			fmt.Fprintf(os.Stderr, "No vault name provided\n")
			return
		}

		vaultFilesystemRepository := filesystem.NewVaultFilesystemRepository()
		vaultService := vault.NewVaultService(vaultFilesystemRepository)

		hashicorpConfig := vaulttypes.VaultHashicorp{
			Address:   hashicorpCmdIn.Address,
			Token:     hashicorpCmdIn.Token,
			Namespace: hashicorpCmdIn.Namespace,
		}

		var config vaulttypes.VaultConfig = &hashicorpConfig

		vault := vaulttypes.Vault{
			Id:     uuid.NewString(),
			Name:   hashicorpCmdIn.VaultName,
			Type:   vaulttypes.VaultTypeHashicorp,
			Config: config,
		}

		vaultId, err := vaultService.ConnectVault(&vault)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error configuring Hashicorp vault: %v\n", err)
			return
		}

		cmd.Printf("Successfully configured Hashicorp vault with ID: %s\n", vaultId)

		err = cliCache.SaveCache(
			&cliCache.Cache{
				VaultId: vaultId,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving local configuration: %v\n", err)
			return
		}

	},
}

//nolint:lll // Allow long lines for CLI
func init() {
	HashicorpCmd.Flags().StringVarP(&hashicorpCmdIn.Address, "address", "a", "", "The address of the HashiCorp Vault instance")
	HashicorpCmd.Flags().StringVarP(&hashicorpCmdIn.Token, "token", "t", "", "The token to authenticate with the HashiCorp Vault instance")
	HashicorpCmd.Flags().StringVarP(&hashicorpCmdIn.Namespace, "namespace", "n", "", "The namespace to use in the HashiCorp Vault instance")
	HashicorpCmd.Flags().StringVarP(&hashicorpCmdIn.VaultName, "vault-name", "v", "", "Name of the vault")
}
