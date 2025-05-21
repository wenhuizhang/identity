// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

import (
	coreV1alpha "github.com/agntcy/identity/api/agntcy/identity/core/v1alpha1"
)

type Badge struct {
	// The badge ID
	Id string `json:"id,omitempty"`
	// The verifiable credential
	EnvelopedCredential *coreV1alpha.EnvelopedCredential `json:"issuer,omitempty"`
}
