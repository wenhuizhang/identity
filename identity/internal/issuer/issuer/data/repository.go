// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type IssuerRepository interface {
	AddIssuer(
		vaultId, identityNodeAddress string, idpConfig internalIssuerTypes.IdpConfig, issuer *coreV1alpha.Issuer,
	) (string, error)
	GetAllIssuers(vaultId string) ([]*internalIssuerTypes.Issuer, error)
	GetIssuer(vaultId, issuerId string) (*coreV1alpha.Issuer, error)
	RemoveIssuer(vaultId, issuerId string) error
}
