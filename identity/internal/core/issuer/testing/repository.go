// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"

	issuercore "github.com/agntcy/identity/internal/core/issuer"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
)

type FakeIssuerRepository struct {
	store map[string]*issuertypes.Issuer
}

func NewFakeIssuerRepository() issuercore.Repository {
	return &FakeIssuerRepository{
		store: make(map[string]*issuertypes.Issuer),
	}
}

func (r *FakeIssuerRepository) CreateIssuer(
	ctx context.Context,
	issuer *issuertypes.Issuer,
) (*issuertypes.Issuer, error) {
	r.store[issuer.CommonName] = issuer
	return issuer, nil
}

func (r *FakeIssuerRepository) GetIssuer(
	ctx context.Context,
	commonName string,
) (*issuertypes.Issuer, error) {
	if issuer, ok := r.store[commonName]; ok {
		return issuer, nil
	}

	return nil, errutil.Err(
		nil, "issuer not found",
	)
}
