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

type FileFlags struct {
	FilePath  string
	VaultName string
}

type FileCommand struct {
	vaultService vault.VaultService
}

func NewCmdFile(vaultService vault.VaultService) *cobra.Command {
	flags := NewFileFlags()

	cmd := &cobra.Command{
		Use:   "file",
		Short: "Create a local vault file to store your cryptographic keys",
		Run: func(cmd *cobra.Command, args []string) {
			c := FileCommand{
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

func NewFileFlags() *FileFlags {
	return &FileFlags{}
}

func (f *FileFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.FilePath, "file-path", "f", "", "Path to the file")
	cmd.Flags().StringVarP(&f.VaultName, "vault-name", "v", "", "Name of the vault")
}

func (cmd *FileCommand) Run(ctx context.Context, flags *FileFlags) error {
	// if the file path is not set, prompt the user for it interactively
	if flags.FilePath == "" {
		fmt.Fprintf(os.Stdout, "File path to store the vault: ")

		_, err := fmt.Scanln(&flags.FilePath)
		if err != nil {
			return fmt.Errorf("error reading file path: %v", err)
		}
	}

	if flags.FilePath == "" {
		return fmt.Errorf("no file path provided")
	}

	// if the vault name is not set, prompt the user for it interactively
	if flags.VaultName == "" {
		fmt.Fprintf(os.Stdout, "Vault name: ")

		_, err := fmt.Scanln(&flags.VaultName)
		if err != nil {
			return fmt.Errorf("error reading vault name: %v", err)
		}
	}

	if flags.VaultName == "" {
		return fmt.Errorf("no vault name provided")
	}

	fileConfig := vaulttypes.VaultFile{
		FilePath: flags.FilePath,
	}

	var config vaulttypes.VaultConfig = &fileConfig

	vault := vaulttypes.Vault{
		Id:     uuid.NewString(),
		Name:   flags.VaultName,
		Type:   vaulttypes.VaultTypeFile,
		Config: config,
	}

	vaultId, err := cmd.vaultService.ConnectVault(&vault)
	if err != nil {
		return fmt.Errorf("error configuring file vault: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Successfully configured file vault with ID: %s\n", vaultId)

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
