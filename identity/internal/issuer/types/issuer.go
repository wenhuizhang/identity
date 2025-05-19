// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: appache

package types

type IssuerConfig struct {
	// The identity node configuration
	IdentityNodeConfig *IdentityNodeConfig `json:"identity_node_config,omitempty" gorm:"embedded"`
	// The identity provider configuration
	IdpConfig *IdpConfig `json:"idp_config,omitempty" gorm:"embedded"`
}
