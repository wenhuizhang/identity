// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	issuersrv "github.com/agntcy/identity/internal/issuer/issuer"
	"github.com/agntcy/identity/internal/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type ShowFlags struct {
	IssuerID string
}

type ShowCommand struct {
	cache         *clicache.Cache
	issuerService issuersrv.IssuerService
}

func NewCmdShow(
	cache *clicache.Cache,
	issuerService issuersrv.IssuerService,
) *cobra.Command {
	flags := NewShowFlags()

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show details of an issuer configuration",
		Run: func(cmd *cobra.Command, args []string) {
			c := ShowCommand{
				cache:         cache,
				issuerService: issuerService,
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
	cmd.Flags().StringVarP(&f.IssuerID, "issuer-id", "i", "", "The ID of the issuer to show")
}

func (cmd *ShowCommand) Run(ctx context.Context, flags *ShowFlags) error {
	err := cmd.cache.ValidateVaultId()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %w", err)
	}

	// if the issuer id is not set, prompt the user for it interactively
	if flags.IssuerID == "" {
		err := cmdutil.ScanRequired("Issuer ID to show", &flags.IssuerID)
		if err != nil {
			return fmt.Errorf("error reading issuer ID: %w", err)
		}
	}

	issuer, err := cmd.issuerService.GetIssuer(cmd.cache.VaultId, cmd.cache.KeyID, flags.IssuerID)
	if err != nil {
		return fmt.Errorf("error getting issuer: %w", err)
	}

	if issuer == nil {
		return fmt.Errorf("no issuer found with ID: %s", flags.IssuerID)
	}

	issuerJSON, err := json.MarshalIndent(issuer, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling metadata to JSON: %w", err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(issuerJSON))

	return nil
}
