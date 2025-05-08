package id

import (
	"context"

	"github.com/agntcy/identity/internal/core/id/types"
)

type IdService interface {
	// Resolves an ID into a ResolverMetadata
	Resolve(ctx context.Context, id string) (*types.ResolverMetadata, error)
}
