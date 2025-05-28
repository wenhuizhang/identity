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

type ListCommand struct {
	cache         *clicache.Cache
	issuerService issuer.IssuerService
}

func NewCmdList(
	cache *clicache.Cache,
	issuerService issuer.IssuerService,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List your existing issuer configurations",
		Long:  "List your existing issuer configurations",
		Run: func(cmd *cobra.Command, args []string) {
			c := ListCommand{
				cache:         cache,
				issuerService: issuerService,
			}

			err := c.Run(cmd.Context())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}

func (cmd *ListCommand) Run(ctx context.Context) error {
	err := cmd.cache.ValidateForIssuer()
	if err != nil {
		return fmt.Errorf("error validating local configuration: %v", err)
	}

	issuers, err := cmd.issuerService.GetAllIssuers(cmd.cache.VaultId, cmd.cache.KeyID)
	if err != nil {
		return fmt.Errorf("error listing issuers: %v", err)
	}

	if len(issuers) == 0 {
		return fmt.Errorf("no issuers found")
	}

	fmt.Fprintf(os.Stdout, "Existing issuers:\n")

	for _, issuer := range issuers {
		fmt.Fprintf(os.Stdout, "- %s, %s\n", issuer.ID, issuer.CommonName)
	}

	return nil
}
