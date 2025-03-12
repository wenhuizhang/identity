package did

import (
	"context"

	"github.com/agntcy/pyramid/internal/core/did/types"
	"github.com/agntcy/pyramid/internal/core/pagination"
)

type DidRepository interface {
	GetDids(
		ctx context.Context,
		paginationFilter pagination.PaginationFilter,
		query *string,
	) (*pagination.Pageable[types.Did], error)
	CreateDid(
		ctx context.Context,
		pyramid *types.Did,
	) (*types.Did, error)
	GetDid(
		ctx context.Context,
		id string,
		withFields ...string,
	) (*types.Did, error)
	GetDidByCatalogID(
		ctx context.Context,
		catalogID string,
	) (*types.Did, error)
	UpdateDid(
		ctx context.Context,
		pyramid *types.Did,
	) (*types.Did, error)
	DeleteDid(ctx context.Context, id string) error
}
