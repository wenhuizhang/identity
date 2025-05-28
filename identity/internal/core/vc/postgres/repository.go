// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"

	vccore "github.com/agntcy/identity/internal/core/vc"
	"github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/db"
	"github.com/agntcy/identity/pkg/log"
)

type vcPostgresRepository struct {
	dbContext db.Context
}

func NewRepository(dbContext db.Context) vccore.Repository {
	return &vcPostgresRepository{
		dbContext: dbContext,
	}
}

func (r *vcPostgresRepository) Create(
	ctx context.Context,
	credential *types.VerifiableCredential,
	resolverMetadataID string,
) (*types.VerifiableCredential, error) {
	model := newVerifiableCredentialModel(credential, resolverMetadataID)

	result := r.dbContext.Client().Create(model)
	if result.Error != nil {
		return nil, errutil.Err(
			result.Error, "there was an error creating the verifiable credential",
		)
	}

	return credential, nil
}

func (r *vcPostgresRepository) GetWellKnown(
	ctx context.Context,
	resolverMetadataID string,
) (*[]*types.EnvelopedCredential, error) {
	var credentials []*VerifiableCredential

	log.Debug("Retrieving well-known verifiable credentials for resolver metadata ID: ", resolverMetadataID)

	result := r.dbContext.Client().
		Model(&VerifiableCredential{}).
		Where("resolver_metadata_id = ?", resolverMetadataID).
		Find(&credentials)
	if result.Error != nil {
		return nil, errutil.Err(
			result.Error, "there was an error fetching the verifiable credentials",
		)
	}

	var envelopedCredentials []*types.EnvelopedCredential

	for _, cred := range credentials {
		if cred.Proof == nil {
			log.Debug("Skipping credential with empty proof for ID: ", cred.ID)
		}

		if cred.Proof.Type == "" {
			log.Debug("Skipping credential with empty proof type for ID: ", cred.ID)
		}

		if cred.Proof.ProofValue == "" {
			log.Debug("Skipping credential with empty proof value for ID: ", cred.ID)
		}

		switch cred.Proof.Type {
		case "JWT":
			envelopedCredentials = append(envelopedCredentials, &types.EnvelopedCredential{
				EnvelopeType: types.CREDENTIAL_ENVELOPE_TYPE_JOSE,
				Value:        cred.Proof.ProofValue,
			})
		default:
			log.Debug("Skipping credential with unsupported proof type: ", cred.Proof.Type, " for ID: ", cred.ID)
		}
	}

	return &envelopedCredentials, nil
}
