// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verify

import (
	"encoding/json"
	"fmt"
	"os"

	idtypes "github.com/agntcy/identity/internal/core/id/types"
	"github.com/agntcy/identity/internal/core/vc/jose"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/nodeapi"
	"github.com/lestrrat-go/jwx/v3/jws"

	"github.com/spf13/cobra"
)

type VerifyCmdInput struct {
	IdentityNodeURL string
	BadgeFilePath   string
}

var (
	// setup the verify command flags
	verifyCmdIn = &VerifyCmdInput{}
)

const (
	defaultNodeAddress = "http://localhost:4000"
)

var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify an Agent or MCP Server Badge from a file",
	Run: func(cmd *cobra.Command, args []string) {

		// if the file path is not set, prompt the user for it interactively
		if verifyCmdIn.BadgeFilePath == "" {
			fmt.Fprintf(os.Stderr, "Full file path to the badge file: ")
			_, err := fmt.Scanln(&verifyCmdIn.BadgeFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file path: %v\n", err)
				return
			}
		}
		if verifyCmdIn.BadgeFilePath == "" {
			fmt.Fprintf(os.Stderr, "No file path provided\n")
			return
		}

		// if the identity node address is not set, prompt the user for it interactively
		if verifyCmdIn.IdentityNodeURL == "" {
			fmt.Fprintf(os.Stdout, "Identity node address (default %s): ", defaultNodeAddress)

			_, err := fmt.Scanln(&verifyCmdIn.IdentityNodeURL)
			if err != nil {
				// If the user just presses Enter, registerIdentityNodeAddress will be "" and err will
				// be an "unexpected newline" error. We should allow this and use the default value.
				if err.Error() != "unexpected newline" {
					fmt.Fprintf(os.Stderr, "Error reading identity node address: %v\n", err)
					return
				}
			}
		}
		// If no address was entered (input was empty or only whitespace), use the default.
		if verifyCmdIn.IdentityNodeURL == "" {
			verifyCmdIn.IdentityNodeURL = defaultNodeAddress
		}

		// Check if the badge file exists
		if _, err := os.Stat(verifyCmdIn.BadgeFilePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "File does not exist: %s\n", verifyCmdIn.BadgeFilePath)
			return
		}

		// Read the badge file
		vcData, err := os.ReadFile(verifyCmdIn.BadgeFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			return
		}

		// Unmarshal the badge data
		var vc vctypes.EnvelopedCredential
		if err := json.Unmarshal(vcData, &vc); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshaling badge data: %v\n", err)
			return
		}

		// Create a node client to interact with the identity node
		nodeClientPrv := nodeapi.NewNodeClientProvider()
		client, err := nodeClientPrv.New(verifyCmdIn.IdentityNodeURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating node client: %v\n", err)
			return
		}
		if client == nil {
			fmt.Fprintf(os.Stderr, "Node client is nil, check your identity node address: %s\n", verifyCmdIn.IdentityNodeURL)
			return
		}

		switch vc.EnvelopeType {
		case vctypes.CREDENTIAL_ENVELOPE_TYPE_EMBEDDED_PROOF:
			fmt.Fprintf(os.Stderr, "Badge verification is not supported for embedded proof badges yet.\n")
			return
		case vctypes.CREDENTIAL_ENVELOPE_TYPE_JOSE:

			// Decode the JWT
			raw, err := jws.Parse([]byte(vc.Value))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing JWT: %v\n", err)
				return
			}

			var validatedVC vctypes.VerifiableCredential

			err = json.Unmarshal(raw.Payload(), &validatedVC)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshaling JWT payload: %v\n", err)
				return
			}

			credentialSubject := validatedVC.CredentialSubject
			if credentialSubject == nil {
				fmt.Fprintf(os.Stderr, "Credential subject is missing in the badge\n")
				return
			}

			// extract the Resolver Metadata ID from the credential subject map
			resolverMetadataID, ok := credentialSubject["id"].(string)
			if !ok || resolverMetadataID == "" {
				fmt.Fprintf(os.Stderr, "Resolver Metadata ID is missing or invalid in the badge\n")
				return
			}

			// Resolve the Resolver Metadata ID to get the public key
			resolvedMetadata, err := client.ResolveMetadataByID(cmd.Context(), resolverMetadataID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error resolving Resolver Metadata ID: %v\n", err)
				return
			}
			if resolvedMetadata == nil {
				fmt.Fprintf(os.Stderr, "Resolver Metadata not found for ID: %s\n", resolverMetadataID)
				return
			}

			// convert resolvedMetadata.VerificationMethods to JWKs
			var jwks idtypes.Jwks
			for _, vm := range resolvedMetadata.VerificationMethod {
				jwks.Keys = append(jwks.Keys, vm.PublicKeyJwk)
			}

			// Verify the badge using the Resolver Metadata public key
			parsedVC, err := jose.Verify(&jwks, &vc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error verifying badge: %v\n", err)
				return
			}

			fmt.Fprintf(os.Stdout, "\nBadge verified successfully!\n\n")
			fmt.Fprintf(os.Stdout, "Badge ID: %s\n", parsedVC.ID)
			fmt.Fprintf(os.Stdout, "Badge Type: %s\n", parsedVC.Type)
			fmt.Fprintf(os.Stdout, "Badge Issuer: %s\n", parsedVC.Issuer)
			fmt.Fprintf(os.Stdout, "Badge Issuance Date: %s\n", parsedVC.IssuanceDate)

			// Create a more readable output of credential subject
			credentialSubjectBytes, err := json.MarshalIndent(parsedVC.CredentialSubject, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting credential subject: %v\n", err)
				return
			}

			fmt.Fprintf(os.Stdout, "Badge Credential Subject: %v\n\n", string(credentialSubjectBytes))
		default:
			fmt.Fprintf(os.Stderr, "Unsupported badge envelope type: %s\n", vc.EnvelopeType)
		}
	},
}

func init() {
	VerifyCmd.Flags().StringVarP(&verifyCmdIn.BadgeFilePath, "file", "f", "", "Path to the badge file")
	VerifyCmd.Flags().StringVarP(&verifyCmdIn.IdentityNodeURL, "identity-node-address", "i", "", "Identity node address")
}
