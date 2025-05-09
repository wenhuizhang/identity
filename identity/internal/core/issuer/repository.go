// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"

	"github.com/agntcy/identity/internal/core/issuer/types"
	"github.com/agntcy/identity/internal/pkg/pagination"
)

// IssuerRepository is the interface for the Issuer repository
type IssuerRepository interface {
	GetIssuers(
		ctx context.Context,
		paginationFilter pagination.PaginationFilter,
		query *string,
	) (*pagination.Pageable[types.Issuer], error)
	CreateIssuer(
		ctx context.Context,
		issuerentity *types.Issuer,
	) (*types.Issuer, error)
	GetIssuer(
		ctx context.Context,
		issuer string,
		withFields ...string,
	) (*types.Issuer, error)
	UpdateIssuer(
		ctx context.Context,
		issuerentity *types.Issuer,
	) (*types.Issuer, error)
	DeleteIssuer(ctx context.Context, issuer string) error
}
