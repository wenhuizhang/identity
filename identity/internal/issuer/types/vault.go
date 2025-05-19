// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

type Vault struct {
	// The ID of the vault
	Id string `json:"id,omitempty" gorm:"primaryKey;type:varchar(256);"`
}
