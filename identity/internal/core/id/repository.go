package id

import (
	"context"

	"github.com/agntcy/identity/internal/core/id/types"
	"github.com/agntcy/identity/internal/pkg/pagination"
)

type IdRepository interface {
	GetIds(
		ctx context.Context,
		paginationFilter pagination.PaginationFilter,
		query *string,
	) (*pagination.Pageable[types.ResolverMetadata], error)
	CreateId(
		ctx context.Context,
		identity *types.ResolverMetadata,
	) (*types.ResolverMetadata, error)
	GetId(
		ctx context.Context,
		id string,
		withFields ...string,
	) (*types.ResolverMetadata, error)
	GetIdByCatalogID(
		ctx context.Context,
		catalogID string,
	) (*types.ResolverMetadata, error)
	UpdateId(
		ctx context.Context,
		identity *types.ResolverMetadata,
	) (*types.ResolverMetadata, error)
	DeleteId(ctx context.Context, id string) error
}
