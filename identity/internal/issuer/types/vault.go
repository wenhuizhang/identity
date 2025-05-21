// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/json"
	"fmt"
)

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
	FilePath string `json:"path,omitempty"`
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

// UnmarshalVault implements custom JSON unmarshaling for Vault
func (v *Vault) UnmarshalVault(data []byte) error {
	// Temporary struct to decode the JSON data
	type tempVault struct {
		Id     string          `json:"id,omitempty"`
		Type   VaultType       `json:"type,omitempty"`
		Config json.RawMessage `json:"config,omitempty"`
	}

	var temp tempVault
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Assign basic fields
	v.Id = temp.Id
	v.Type = temp.Type

	// Skip if no config is provided
	if len(temp.Config) == 0 {
		v.Config = nil
		return nil
	}

	// Create the appropriate VaultConfig based on the Type
	switch temp.Type {
	case VaultTypeTxt:
		var config VaultTxt
		if err := json.Unmarshal(temp.Config, &config); err != nil {
			return err
		}
		v.Config = &config

	case VaultType1Password:
		var config Vault1Password
		if err := json.Unmarshal(temp.Config, &config); err != nil {
			return err
		}
		v.Config = &config

	default:
		return fmt.Errorf("unknown vault type: %s", temp.Type)
	}

	return nil
}
