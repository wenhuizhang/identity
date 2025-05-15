// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

import (
	idtypes "github.com/agntcy/identity/internal/core/id/types"
)

// A Identity Issuer
type Issuer struct {
	// The organization of the issuer
	Organization string `json:"organization,omitempty" gorm:"not null;type:varchar(256);"`

	// The sub organization of the issuer
	SubOrganization string `json:"sub_organization,omitempty" gorm:"not null;type:varchar(256);"`

	// The common name of the issuer
	// Could be a FQDN or a FQDA
	CommonName string `json:"common_name,omitempty" gorm:"primaryKey"`

	// This field is optional
	// The keys of the issuer in JWK format
	// The public key is used to verify the signature of the different claims
	PublicKey *idtypes.Jwk `json:"public_key,omitempty" gorm:"embedded;embeddedPrefix:public_key_"`

	// This field is optional
	// The private key of the issuer in JWK format
	PrivateKey *idtypes.Jwk `json:"private_key,omitempty" gorm:"-"`
}
