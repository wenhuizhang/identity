// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package types

type IdpConfig struct {
	// The client id of the identity provider
	ClientId string `json:"client_id,omitempty" gorm:"not null;type:varchar(256);"`
	// The client secret of the identity provider
	ClientSecret string `json:"client_secret,omitempty" gorm:"not null;type:varchar(256);"`
	// The issuer url of the identity provider
	IssuerUrl string `json:"issuer_url,omitempty" gorm:"not null;type:varchar(256);"`
}
