// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

type IssuerConfig struct {
	// The identity node configuration
	IdentityNodeConfig *IdentityNodeConfig `json:"identity_node_config,omitempty" gorm:"embedded"`
	// The identity provider configuration
	IdpConfig *IdpConfig `json:"idp_config,omitempty" gorm:"embedded"`
}
