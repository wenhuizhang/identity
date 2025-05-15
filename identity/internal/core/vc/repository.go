// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

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
