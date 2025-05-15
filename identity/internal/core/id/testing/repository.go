// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"

	idcore "github.com/agntcy/identity/internal/core/id"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
)

type FakeIdRepository struct {
	store map[string]*idtypes.ResolverMetadata
}

func NewFakeIdRepository() idcore.IdRepository {
	return &FakeIdRepository{
		store: make(map[string]*idtypes.ResolverMetadata),
	}
}

func (r *FakeIdRepository) CreateID(
	ctx context.Context,
	metadata *idtypes.ResolverMetadata,
) (*idtypes.ResolverMetadata, error) {
	r.store[metadata.ID] = metadata
	return metadata, nil
}

func (r *FakeIdRepository) ResolveID(ctx context.Context, id string) (*idtypes.ResolverMetadata, error) {
	if md, ok := r.store[id]; ok {
		return md, nil
	}

	return nil, nil
}
