// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package configuration

import (
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	bsvc "github.com/agntcy/identity/internal/issuer/badge"
	isvc "github.com/agntcy/identity/internal/issuer/issuer"
	msvc "github.com/agntcy/identity/internal/issuer/metadata"
	vsvc "github.com/agntcy/identity/internal/issuer/vault"
	"github.com/spf13/cobra"
)

type ConfigurationCommand struct {
	cache           *clicache.Cache
	vaultService    vsvc.VaultService
	issuerService   isvc.IssuerService
	metadataService msvc.MetadataService
	badgeService    bsvc.BadgeService
}

func NewCmd(
	cache *clicache.Cache,
	vaultService vsvc.VaultService,
	issuerService isvc.IssuerService,
	metadataService msvc.MetadataService,
	badgeService bsvc.BadgeService,
) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Display the local configuration context",
		Run: func(cmd *cobra.Command, args []string) {
			c := ConfigurationCommand{
				cache:           cache,
				vaultService:    vaultService,
				issuerService:   issuerService,
				metadataService: metadataService,
				badgeService:    badgeService,
			}

			err := c.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
}

func (cmd *ConfigurationCommand) Run() error {
	err := cmd.cache.ValidateVaultId()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %w", err)
	}

	vault, err := cmd.vaultService.GetVault(cmd.cache.VaultId)
	if err != nil {
		return fmt.Errorf("error loading vault: %w", err)
	}

	if vault == nil {
		return fmt.Errorf("no vault found with ID: %s", cmd.cache.VaultId)
	}

	fmt.Fprintf(os.Stdout, "\nCurrent Identity CLI configuration context:\n")
	fmt.Fprintf(os.Stdout, "- Vault: %s (%s vault), id: %s\n", vault.Name, vault.Type, vault.Id)

	if cmd.cache.KeyID != "" {
		fmt.Fprintf(os.Stdout, "- Key ID: %s\n", cmd.cache.KeyID)
	} else {
		fmt.Fprintf(os.Stdout, "- Key ID: Not set\n")
	}

	err = cmd.verifyIssuer()
	if err != nil {
		return err
	}

	err = cmd.verifyResolverMetadata()
	if err != nil {
		return err
	}

	err = cmd.verifyBadge()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\n")

	return nil
}

func (cmd *ConfigurationCommand) verifyIssuer() error {
	if cmd.cache.IssuerId == "" {
		return nil
	}

	issuer, err := cmd.issuerService.GetIssuer(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
	)
	if err != nil {
		return fmt.Errorf("error loading issuer: %w", err)
	}

	if issuer == nil {
		return fmt.Errorf("no issuer found with ID: %s", cmd.cache.IssuerId)
	}

	fmt.Fprintf(os.Stdout, "- Issuer: %s\n", issuer.ID)

	return nil
}

func (cmd *ConfigurationCommand) verifyResolverMetadata() error {
	if cmd.cache.MetadataId == "" {
		return nil
	}

	metadata, err := cmd.metadataService.GetMetadata(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		cmd.cache.MetadataId,
	)
	if err != nil {
		return fmt.Errorf("error loading metadata: %w", err)
	}

	if metadata == nil {
		return fmt.Errorf("no metadata found with ID: %s", cmd.cache.MetadataId)
	}

	fmt.Fprintf(os.Stdout, "- Metadata: %s\n", metadata.ID)

	return nil
}

func (cmd *ConfigurationCommand) verifyBadge() error {
	if cmd.cache.BadgeId == "" {
		return nil
	}

	badge, err := cmd.badgeService.GetBadge(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		cmd.cache.IssuerId,
		cmd.cache.MetadataId,
		cmd.cache.BadgeId,
	)
	if err != nil {
		return fmt.Errorf("error loading badge: %w", err)
	}

	if badge == nil {
		return fmt.Errorf("no badge found with ID: %s", cmd.cache.BadgeId)
	}

	fmt.Fprintf(os.Stdout, "- Badge: %s\n", badge.Id)

	return nil
}
