// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	"github.com/agntcy/identity/internal/issuer/issuer/types"
)

type IssuerRepository interface {
	AddIssuer(vaultId string, issuer *types.Issuer) (string, error)
	GetAllIssuers(vaultId string) ([]*types.Issuer, error)
	GetIssuer(vaultId, issuerId string) (*types.Issuer, error)
	RemoveIssuer(vaultId, issuerId string) error
}
