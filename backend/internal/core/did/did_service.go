package did

import (
	"context"

	"github.com/agntcy/pyramid/internal/core/did/types"
)

type DidService interface {
	// Get a Did by id
	Get(ctx context.Context, id string) (*types.Did, error)
}
