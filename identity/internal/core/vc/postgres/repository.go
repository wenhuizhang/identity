// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"

	vccore "github.com/agntcy/identity/internal/core/vc"
	"github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/db"
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
) (*types.VerifiableCredential, error) {
	model := newVerifiableCredentialModel(credential)

	result := r.dbContext.Client().Create(model)
	if result.Error != nil {
		return nil, errutil.Err(
			result.Error, "there was an error creating the verifiable credential",
		)
	}

	return credential, nil
}
