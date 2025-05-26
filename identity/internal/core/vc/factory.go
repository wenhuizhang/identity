// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vc

import (
	"encoding/json"
	"time"

	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	"github.com/agntcy/identity/internal/core/vc/types"
	"github.com/google/uuid"
)

type VerifiableCredentialOption func(*types.VerifiableCredential) error

func New(options ...VerifiableCredentialOption) (*types.VerifiableCredential, error) {
	vc := &types.VerifiableCredential{
		Context: []string{
			"https://www.w3.org/ns/credentials/v2",
			"https://www.w3.org/ns/credentials/examples/v2",
		},
		ID:           uuid.NewString(),
		IssuanceDate: time.Now().UTC().Format(time.RFC3339),
	}
	for _, opt := range options {
		err := opt(vc)
		if err != nil {
			return nil, err
		}
	}

	return vc, nil
}

func WithIssuer(issuer *issuertypes.Issuer) VerifiableCredentialOption {
	return func(vc *types.VerifiableCredential) error {
		vc.Issuer = issuer.Organization
		return nil
	}
}

func WithCredentialContent[T types.CredentialSubject](
	content *types.CredentialContent[T],
) VerifiableCredentialOption {
	return func(vc *types.VerifiableCredential) error {
		sub, err := json.Marshal(content.Content)
		if err != nil {
			return err
		}

		vc.Type = append(vc.Type, content.Type.String())
		vc.CredentialSubject = string(sub)
		return nil
	}
}

// Schemas can be used to include JSON Schemas within the Verifiable Credential created by [Create]
// more information can be found [here]
//
// [here]: https://www.w3.org/TR/vc-data-model-2.0/#data-schemas
func WithCredentialSchema(schemas ...string) VerifiableCredentialOption {
	return func(vc *types.VerifiableCredential) error {
		if vc.CredentialSchema != nil {
			vc.CredentialSchema = make([]*types.CredentialSchema, 0, len(schemas))
		}

		for _, schema := range schemas {
			vc.CredentialSchema = append(vc.CredentialSchema, &types.CredentialSchema{
				Type: "JsonSchema",
				ID:   schema,
			})
		}

		return nil
	}
}
