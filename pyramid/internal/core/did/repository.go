package id

import (
	"context"

	"github.com/agntcy/pyramid/internal/core/id/types"
	"github.com/agntcy/pyramid/internal/pkg/pagination"
)

type IdRepository interface {
	GetIds(
		ctx context.Context,
		paginationFilter pagination.PaginationFilter,
		query *string,
	) (*pagination.Pageable[types.Id], error)
	CreateId(
		ctx context.Context,
		pyramid *types.Id,
	) (*types.Id, error)
	GetId(
		ctx context.Context,
		id string,
		withFields ...string,
	) (*types.Id, error)
	GetIdByCatalogID(
		ctx context.Context,
		catalogID string,
	) (*types.Id, error)
	UpdateId(
		ctx context.Context,
		pyramid *types.Id,
	) (*types.Id, error)
	DeleteId(ctx context.Context, id string) error
}
