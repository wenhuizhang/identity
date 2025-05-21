// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

type MetadataConfig struct {
	// The identity provider configuration
	IdpConfig *IdpConfig `json:"idp_config,omitempty"`
}
