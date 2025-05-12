// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"
	"net/http"

	"github.com/agntcy/identity/internal/core/issuer/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/go-kivik/kivik/v4"
)

// Repository is the interface for the Issuer repository
type Repository interface {
	CreateIssuer(
		ctx context.Context,
		issuer *types.Issuer,
	) (*types.Issuer, error)
}

type repository struct {
	dbContext *kivik.DB
}

// NewIssuerRepository creates a new instance of the IssuerRepository
func NewRepository(dbContext *kivik.DB) Repository {
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
	if _, err := r.dbContext.Put(ctx, issuer.CommonName, issuer); err != nil {
		if kivik.HTTPStatus(err) == http.StatusConflict {
			return nil, errutil.Err(err, "issuer exists")
		}

		return nil, errutil.Err(err, "failed to create issuer")
	}

	return issuer, nil
}
