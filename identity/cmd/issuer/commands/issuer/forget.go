// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"
	"fmt"
	"os"

	clicache "github.com/agntcy/identity/cmd/issuer/cache"
	issuer "github.com/agntcy/identity/internal/issuer/issuer"
	"github.com/spf13/cobra"
)

type ForgetFlags struct {
	IssuerID string
}

type ForgetCommand struct {
	cache         *clicache.Cache
	issuerService issuer.IssuerService
}

func NewCmdForget(
	cache *clicache.Cache,
	issuerService issuer.IssuerService,
) *cobra.Command {
	flags := NewForgetFlags()

	cmd := &cobra.Command{
		Use:   "forget",
		Short: "Forget an issuer configuration",
		Run: func(cmd *cobra.Command, args []string) {
			c := ForgetCommand{
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

func NewForgetFlags() *ForgetFlags {
	return &ForgetFlags{}
}

func (f *ForgetFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.IssuerID, "issuer-id", "i", "", "The ID of the issuer to forget")
}

func (cmd *ForgetCommand) Run(ctx context.Context, flags *ForgetFlags) error {
	err := cmd.cache.ValidateForIssuer()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	// if the issuer id is not set, prompt the user for it interactively
	if flags.IssuerID == "" {
		fmt.Fprintf(os.Stdout, "Issuer ID to forget:\n")

		_, err := fmt.Scanln(&flags.IssuerID)
		if err != nil {
			return fmt.Errorf("error reading issuer ID: %v", err)
		}
	}

	if flags.IssuerID == "" {
		return fmt.Errorf("no issuer ID provided")
	}

	err = cmd.issuerService.ForgetIssuer(
		cmd.cache.VaultId,
		cmd.cache.KeyID,
		flags.IssuerID,
	)
	if err != nil {
		return fmt.Errorf("error forgetting issuer: %v", err)
	}

	// If the issuer was the current issuer in the cache, clear the cache of issuer, metadata, and badge IDs
	if cmd.cache.IssuerId == flags.IssuerID {
		cmd.cache.IssuerId = ""
		cmd.cache.MetadataId = ""
		cmd.cache.BadgeId = ""
		err = clicache.SaveCache(cmd.cache)
		if err != nil {
			return fmt.Errorf("error saving local configuration: %v", err)
		}
	}

	fmt.Fprintf(os.Stdout, "Forgot issuer with ID: %s\n", flags.IssuerID)

	return nil
}
