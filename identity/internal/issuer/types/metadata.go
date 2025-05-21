// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
)

type Metadata struct {
	// The metadata ID
	Id string `json:"id,omitempty"`
	// The metadata
	ResolverMetadata *coreV1alpha.ResolverMetadata `json:"metadata,omitempty"`
	// The identity provider configuration
	IdpConfig *IdpConfig `json:"idp_config,omitempty"`
}
