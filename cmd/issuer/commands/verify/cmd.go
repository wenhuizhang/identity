// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	v1alphaclient "github.com/agntcy/identity/api/client/models"
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
		Short: "Verify badges from a file",
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

	// Unmarshal the badge data into the expected structure
	var vcs v1alphaclient.V1alpha1GetVcWellKnownResponse
	if err := json.Unmarshal(vcData, &vcs); err != nil {
		return fmt.Errorf("error unmarshalling badge data: %w", err)
	}

	if len(vcs.Vcs) == 0 {
		return fmt.Errorf("no verifiable credentials found in the file: %s", flags.BadgeFilePath)
	}

	// for each Verifiable Credential in the response, verify it
	for _, vc := range vcs.Vcs {
		var envelopedCredential vctypes.EnvelopedCredential
		envelopedCredential.Value = vc.Value

		// Set the envelope type based on the provided type
		switch *vc.EnvelopeType {
		case v1alphaclient.V1alpha1CredentialEnvelopeTypeCREDENTIALENVELOPETYPEJOSE:
			envelopedCredential.EnvelopeType = vctypes.CREDENTIAL_ENVELOPE_TYPE_JOSE
		case v1alphaclient.V1alpha1CredentialEnvelopeTypeCREDENTIALENVELOPETYPEEMBEDDEDPROOF:
			fmt.Fprintf(os.Stdout, "Embedded proof envelope type is not supported yet, skipping: %s\n", *vc.EnvelopeType)
			continue
		default:
			fmt.Fprintf(os.Stdout, "Skipping unsupported envelope type: %s\n", *vc.EnvelopeType)
			continue
		}

		verifiedVC, err := cmd.verifyService.VerifyCredential(ctx, &envelopedCredential, flags.IdentityNodeURL)
		if err != nil {
			return err
		}

		// Print verification results
		if err := printVerifiedBadgeInfo(verifiedVC); err != nil {
			return fmt.Errorf("error printing badge info: %w", err)
		}
	}

	return nil
}

func printVerifiedBadgeInfo(verifiedVC *vctypes.VerifiableCredential) error {
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
