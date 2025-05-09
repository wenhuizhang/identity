// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package id

import (
	"context"

	"github.com/agntcy/identity/internal/core/id/types"
)

type IdService interface {
	// Resolves an ID into a ResolverMetadata
	Resolve(ctx context.Context, id string) (*types.ResolverMetadata, error)
}
