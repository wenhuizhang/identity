// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type IssuerRepository interface {
	AddIssuer(vaultId string, issuer *internalIssuerTypes.Issuer) (string, error)
	GetAllIssuers(vaultId string) ([]*internalIssuerTypes.Issuer, error)
	GetIssuer(vaultId, issuerId string) (*internalIssuerTypes.Issuer, error)
	RemoveIssuer(vaultId, issuerId string) error
}
