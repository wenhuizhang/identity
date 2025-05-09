// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package issuer

import (
	"context"

	"github.com/agntcy/identity/internal/core/issuer/types"
)

type IssuerService interface {
	// Register a new Issuer
	Register(ctx context.Context, issuer *types.Issuer) (*string, error)
}
