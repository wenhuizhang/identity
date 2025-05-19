// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"

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
) (*vctypes.VerifiableCredential, error) {
	r.store[credential.ID] = credential
	return credential, nil
}
