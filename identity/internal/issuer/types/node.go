// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package types

type IdentityNodeConfig struct {
	// The address of the identity node
	IdentityNodeAddress string `json:"identity_node_address,omitempty" gorm:"not null;type:varchar(256);"`
}
