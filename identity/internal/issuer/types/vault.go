// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

// VaultType represents the type of vault
type VaultType string

const (
	VaultTypeTxt       VaultType = "txt"
	VaultType1Password VaultType = "1password"
)

// VaultConfig is an interface that all vault implementations must satisfy
type VaultConfig interface {
	GetVaultType() VaultType
}

type Vault struct {
	// The ID of the vault
	Id string `json:"id,omitempty"`
	// The type of the vault
	Type VaultType `json:"type,omitempty"`
	// The vault implementation
	Config VaultConfig `json:"config,omitempty"`
}

type VaultTxt struct {
	// The text file vault path
	Path string `json:"path,omitempty"`
}

// GetVaultType returns the type of this vault implementation
func (v *VaultTxt) GetVaultType() VaultType {
	return VaultTypeTxt
}

type Vault1Password struct {
	// The 1Password service account token
	ServiceAccountToken string `json:"serviceAccountToken,omitempty"`
	// The 1Password vault ID
	VaultID string `json:"vaultId,omitempty"`
	// The 1Password item ID
	ItemID string `json:"itemId,omitempty"`
}

// GetVaultType returns the type of this vault implementation
func (v *Vault1Password) GetVaultType() VaultType {
	return VaultType1Password
}
