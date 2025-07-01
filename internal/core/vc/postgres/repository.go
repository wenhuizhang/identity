// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"errors"

	errcore "github.com/agntcy/identity/internal/core/errors"
	vccore "github.com/agntcy/identity/internal/core/vc"
	"github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/db"
	"gorm.io/gorm"
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

func (r *vcPostgresRepository) Update(
	ctx context.Context,
	credential *types.VerifiableCredential,
	resolverMetadataID string,
) (*types.VerifiableCredential, error) {
	model := newVerifiableCredentialModel(credential, resolverMetadataID)

	err := r.dbContext.Client().Save(model).Error
	if err != nil {
		return nil, errutil.Err(
			err, "there was an error updating the verifiable credential",
		)
	}

	return credential, nil
}

func (r *vcPostgresRepository) GetByResolverMetadata(
	ctx context.Context,
	resolverMetadataID string,
) ([]*types.VerifiableCredential, error) {
	var storedVCs []*VerifiableCredential

	result := r.dbContext.Client().
		Model(&VerifiableCredential{}).
		Where("resolver_metadata_id = ?", resolverMetadataID).
		Find(&storedVCs)
	if result.Error != nil {
		return nil, errutil.Err(
			result.Error, "there was an error fetching the verifiable credentials",
		)
	}

	vcs := make([]*types.VerifiableCredential, 0, len(storedVCs))
	for _, vc := range storedVCs {
		vcs = append(vcs, vc.ToCoreType())
	}

	return vcs, nil
}

func (r *vcPostgresRepository) GetByID(
	ctx context.Context,
	id string,
) (*types.VerifiableCredential, error) {
	var vc VerifiableCredential

	err := r.dbContext.Client().Where("id = ?", id).First(&vc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcore.ErrResourceNotFound
		}

		return nil, errutil.Err(
			err, "there was an error fetching the verifiable credential",
		)
	}

	return vc.ToCoreType(), nil
}
