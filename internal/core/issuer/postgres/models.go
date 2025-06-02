// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	id "github.com/agntcy/identity/internal/core/id/postgres"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
	"github.com/agntcy/identity/internal/core/issuer/types"
)

type Issuer struct {
	CommonName       string                 `gorm:"primaryKey"`
	Organization     string                 `gorm:"not null;type:varchar(256);"`
	SubOrganization  string                 `gorm:"not null;type:varchar(256);"`
	PublicKey        *idtypes.Jwk           `gorm:"embedded;embeddedPrefix:public_key_"`
	ResolverMetadata []*id.ResolverMetadata `gorm:"foreignKey:IssuerCommonName"`
}

func (i *Issuer) ToCoreType() *types.Issuer {
	return &types.Issuer{
		Organization:    i.Organization,
		SubOrganization: i.SubOrganization,
		CommonName:      i.CommonName,
		PublicKey:       i.PublicKey,
	}
}

func newIssuerModel(src *types.Issuer) *Issuer {
	return &Issuer{
		CommonName:      src.CommonName,
		Organization:    src.Organization,
		SubOrganization: src.SubOrganization,
		PublicKey:       src.PublicKey,
	}
}
