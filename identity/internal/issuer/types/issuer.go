// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
)

type IssuerConfig struct {
	// The identity node configuration
	IdentityNodeConfig *IdentityNodeConfig `json:"identity_node_config,omitempty"`
	// The identity provider configuration
	IdpConfig *IdpConfig `json:"idp_config,omitempty"`
}

type Issuer struct {
	// The issuer ID
	Id string `json:"id,omitempty"`
	// The issuer
	Issuer *coreV1alpha.Issuer `json:"issuer,omitempty"`
}
