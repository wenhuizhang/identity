// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	verifysrv "github.com/agntcy/identity/internal/issuer/verify"
	"github.com/agntcy/identity/internal/pkg/cmdutil"

	"github.com/spf13/cobra"
)

const (
	defaultNodeAddress = "http://localhost:4000"
)

type VerifyFlags struct {
	IdentityNodeURL string
	BadgeFilePath   string
}

type VerifyCommand struct {
	verifyService verifysrv.VerifyService
}

func NewCmd(verifyService verifysrv.VerifyService) *cobra.Command {
	flags := NewVerifyFlags()

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify an Agent or MCP Server Badge from a file",
		Run: func(cmd *cobra.Command, args []string) {
			c := VerifyCommand{
				verifyService: verifyService,
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

func NewVerifyFlags() *VerifyFlags {
	return &VerifyFlags{}
}

func (f *VerifyFlags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.BadgeFilePath, "file", "f", "", "Path to the badge file")
	cmd.Flags().
		StringVarP(&f.IdentityNodeURL, "identity-node-address", "i", "", "Identity node address")
}

func (cmd *VerifyCommand) Run(ctx context.Context, flags *VerifyFlags) error {
	// if the file path is not set, prompt the user for it interactively
	if flags.BadgeFilePath == "" {
		err := cmdutil.ScanRequired("Full file path to the badge file", &flags.BadgeFilePath)
		if err != nil {
			return fmt.Errorf("error reading file path: %w", err)
		}
	}

	// if the identity node address is not set, prompt the user for it interactively
	if flags.IdentityNodeURL == "" {
		err := cmdutil.ScanWithDefault(
			"Identity node address",
			defaultNodeAddress,
			&flags.IdentityNodeURL,
		)
		if err != nil {
			return fmt.Errorf("error reading identity node address: %w", err)
		}
	}

	// Check if the badge file exists
	if _, err := os.Stat(flags.BadgeFilePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", flags.BadgeFilePath)
	}

	// Read the badge file
	vcData, err := os.ReadFile(flags.BadgeFilePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal the badge data
	var vc vctypes.EnvelopedCredential
	if err := json.Unmarshal(vcData, &vc); err != nil {
		return fmt.Errorf("error unmarshaling badge data: %w", err)
	}

	verifiedVC, err := cmd.verifyService.VerifyCredential(ctx, &vc, flags.IdentityNodeURL)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "\nBadge verified successfully!\n\n")
	fmt.Fprintf(os.Stdout, "Badge ID: %s\n", verifiedVC.ID)
	fmt.Fprintf(os.Stdout, "Badge Type: %s\n", verifiedVC.Type)
	fmt.Fprintf(os.Stdout, "Badge Issuer: %s\n", verifiedVC.Issuer)
	fmt.Fprintf(os.Stdout, "Badge Issuance Date: %s\n", verifiedVC.IssuanceDate)

	// Create a more readable output of credential subject
	credentialSubjectBytes, err := json.MarshalIndent(verifiedVC.CredentialSubject, "", "  ")
	if err != nil {
		return fmt.Errorf("error formatting credential subject: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Badge Credential Subject: %v\n\n", string(credentialSubjectBytes))

	return nil
}
