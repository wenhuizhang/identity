// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"

	errcore "github.com/agntcy/identity/internal/core/errors"
	vccore "github.com/agntcy/identity/internal/core/vc"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
)

type FakeVCRepository struct {
	store map[string]*vctypes.VerifiableCredential
}

func NewFakeVCRepository() vccore.Repository {
	return &FakeVCRepository{
		store: make(map[string]*vctypes.VerifiableCredential),
	}
}

func (r *FakeVCRepository) Create(
	ctx context.Context,
	credential *vctypes.VerifiableCredential,
	resolverMetadataID string,
) (*vctypes.VerifiableCredential, error) {
	r.store[credential.ID] = credential
	return credential, nil
}

func (r *FakeVCRepository) GetWellKnown(
	ctx context.Context,
	resolverMetadataID string,
) (*[]*vctypes.EnvelopedCredential, error) {
	var credentials []*vctypes.EnvelopedCredential
	return &credentials, nil
}

func (r *FakeVCRepository) GetByResolverMetadata(
	ctx context.Context,
	resolverMetadataID string,
) ([]*vctypes.VerifiableCredential, error) {
	result := make([]*vctypes.VerifiableCredential, 0)

	for _, vc := range r.store {
		if vc.CredentialSubject["id"] == resolverMetadataID {
			result = append(result, vc)
		}
	}

	return result, nil
}

func (r *FakeVCRepository) GetByID(
	ctx context.Context,
	id string,
) (*vctypes.VerifiableCredential, error) {
	if credential, ok := r.store[id]; ok {
		return credential, nil
	}

	return nil, errcore.ErrResourceNotFound
}

func (r *FakeVCRepository) Update(
	ctx context.Context,
	credential *vctypes.VerifiableCredential,
	resolverMetadataID string,
) (*vctypes.VerifiableCredential, error) {
	r.store[credential.ID] = credential
	return credential, nil
}
