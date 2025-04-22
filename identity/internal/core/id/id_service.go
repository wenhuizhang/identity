package id

import (
	"context"

	"github.com/agntcy/identity/internal/core/id/types"
)

type IdService interface {
	// Get a IdDocument by Id
	Get(ctx context.Context, id string) (*types.IdDocument, error)
}
