// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"

	"github.com/agntcy/identity/internal/core/issuer/types"
	"github.com/agntcy/identity/pkg/db"
)

// Repository is the interface for the Issuer repository
type Repository interface {
	CreateIssuer(
		ctx context.Context,
		issuer *types.Issuer,
	) (*types.Issuer, error)
}

type repository struct {
	dbContext *db.Context
}

// NewIssuerRepository creates a new instance of the IssuerRepository
func NewRepository(dbContext *db.Context) Repository {
	return &repository{
		dbContext,
	}
}

// CreateIssuer creates a new Issuer
func (r *repository) CreateIssuer(
	ctx context.Context,
	issuer *types.Issuer,
) (*types.Issuer, error) {
	// Create the issuer

	return issuer, nil
}
