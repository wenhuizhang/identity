// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type VaultRepository interface {
	AddVault(vault *internalIssuerTypes.Vault) (*internalIssuerTypes.Vault, error)
	GetAllVaults() ([]*internalIssuerTypes.Vault, error)
	GetVault(vaultId string) (*internalIssuerTypes.Vault, error)
	RemoveVault(vaultId string) error
}
