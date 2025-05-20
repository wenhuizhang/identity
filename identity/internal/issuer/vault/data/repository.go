// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package data

import (
	internalIssuerTypes "github.com/agntcy/identity/internal/issuer/types"
)

type VaultRepository interface {
	ConnectVault(vault *internalIssuerTypes.Vault) (*internalIssuerTypes.Vault, error)
	ListVaultIds() ([]string, error)
	GetVault(vaultId string) (*internalIssuerTypes.Vault, error)
	ForgetVault(vaultId string) error
}
