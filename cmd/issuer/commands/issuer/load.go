// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	issuersrv "github.com/agntcy/identity/internal/issuer/issuer"
	"github.com/agntcy/identity/internal/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type LoadFlags struct {
	IssuerID string
}

type LoadCommand struct {
	cache         *clicache.Cache
	issuerService issuersrv.IssuerService
}

func NewCmdLoad(
	cache *clicache.Cache,
	issuerService issuersrv.IssuerService,
) *cobra.Command {
	flags := NewLoadFlags()

	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load an issuer configuration",
		Run: func(cmd *cobra.Command, args []string) {
			c := LoadCommand{
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

func NewLoadFlags() *LoadFlags {
	return &LoadFlags{}
}

func (f *LoadFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.IssuerID, "issuer-id", "i", "", "The ID of the issuer to load")
}

func (cmd *LoadCommand) Run(ctx context.Context, flags *LoadFlags) error {
	err := cmd.cache.ValidateForIssuer()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %w", err)
	}

	// if the issuer id is not set, prompt the user for it interactively
	err = cmdutil.ScanRequiredIfNotSet("Issuer ID to load", &flags.IssuerID)
	if err != nil {
		return fmt.Errorf("error reading issuer ID: %w", err)
	}

	// check the issuer id is valid
	issuer, err := cmd.issuerService.GetIssuer(cmd.cache.VaultId, cmd.cache.KeyID, flags.IssuerID)
	if err != nil {
		return fmt.Errorf("error getting issuer: %w", err)
	}

	if issuer == nil {
		return fmt.Errorf("no issuer found with ID: %s", flags.IssuerID)
	}

	// save the issuer id to the cache
	cmd.cache.IssuerId = flags.IssuerID

	err = clicache.SaveCache(cmd.cache)
	if err != nil {
		return fmt.Errorf("error saving local configuration: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Loaded issuer with ID: %s\n", flags.IssuerID)

	return nil
}
