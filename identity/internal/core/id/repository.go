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
	) (*pagination.Pageable[types.IdDocument], error)
	CreateId(
		ctx context.Context,
		identity *types.IdDocument,
	) (*types.IdDocument, error)
	GetId(
		ctx context.Context,
		id string,
		withFields ...string,
	) (*types.IdDocument, error)
	GetIdByCatalogID(
		ctx context.Context,
		catalogID string,
	) (*types.IdDocument, error)
	UpdateId(
		ctx context.Context,
		identity *types.IdDocument,
	) (*types.IdDocument, error)
	DeleteId(ctx context.Context, id string) error
}
