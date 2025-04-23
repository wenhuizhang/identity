package id

import (
	"context"

	"github.com/agntcy/identity/internal/core/id/types"
)

type IdService interface {
	// Get a ResolverMetadata by Id
	Get(ctx context.Context, id string) (*types.ResolverMetadata, error)
}
