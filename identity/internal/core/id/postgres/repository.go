// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"errors"

	idcore "github.com/agntcy/identity/internal/core/id"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/db"
	"gorm.io/gorm"
)

type idPostgresRepository struct {
	dbContext db.Context
}

func NewIdRepository(dbContext db.Context) idcore.IdRepository {
	return &idPostgresRepository{
		dbContext: dbContext,
	}
}

func (r *idPostgresRepository) CreateID(
	ctx context.Context,
	metadata *idtypes.ResolverMetadata,
) (*idtypes.ResolverMetadata, error) {
	model := newResolverMetadataModel(metadata)

	result := r.dbContext.Client().Create(model)
	if result.Error != nil {
		return nil, errutil.Err(
			result.Error, "there was an error creating the resolver metadata",
		)
	}

	return metadata, nil
}

func (r *idPostgresRepository) ResolveID(
	ctx context.Context,
	id string,
) (*idtypes.ResolverMetadata, error) {
	var metadata ResolverMetadata
	result := r.dbContext.Client().First(&metadata, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errutil.Err(
			result.Error, "there was an error fetching the resolver metadata",
		)
	}

	return metadata.ToCoreType(), nil
}
