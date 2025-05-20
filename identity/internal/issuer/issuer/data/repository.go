// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type IssuerRepository interface {
	RegisterIssuer(vaultId, identityNodeAddress string, idpConfig internalIssuerTypes.IdpConfig) (string, error)
	ListIssuerIds(vaultId string) ([]string, error)
	GetIssuer(vaultId, issuerId string) (*coreV1alpha.Issuer, error)
	ForgetIssuer(vaultId, issuerId string) error
}
