package vc

import (
	"context"

	"github.com/agntcy/identity/internal/core/vc/types"
)

type Repository interface {
	Create(
		ctx context.Context,
		credential *types.VerifiableCredential,
	) (*types.VerifiableCredential, error)
}
