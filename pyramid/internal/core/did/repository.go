package id

import (
	"context"

	"github.com/agntcy/pyramid/internal/core/did/types"
	"github.com/agntcy/pyramid/internal/pkg/pagination"
)

type DidRepository interface {
	GetDids(
		ctx context.Context,
		paginationFilter pagination.PaginationFilter,
		query *string,
	) (*pagination.Pageable[types.DidDocument], error)
	CreateDid(
		ctx context.Context,
		pyramid *types.DidDocument,
	) (*types.DidDocument, error)
	GetDid(
		ctx context.Context,
		id string,
		withFields ...string,
	) (*types.DidDocument, error)
	GetDidByCatalogID(
		ctx context.Context,
		catalogID string,
	) (*types.DidDocument, error)
	UpdateDid(
		ctx context.Context,
		pyramid *types.DidDocument,
	) (*types.DidDocument, error)
	DeleteDid(ctx context.Context, id string) error
}
