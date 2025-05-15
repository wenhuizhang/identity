package postgres

import (
	"time"

	"github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/converters"
	"github.com/lib/pq"
)

type VerifiableCredential struct {
	ID                string `gorm:"primarykey"`
	CreatedAt         time.Time
	Context           pq.StringArray `gorm:"type:text[]"`
	Type              pq.StringArray `gorm:"type:text[]"`
	Issuer            string
	CredentialSubject string
	IssuanceDate      string
	ExpirationDate    string
	CredentialSchema  []*CredentialSchema `gorm:"foreignKey:VerifiableCredentialID"`
	Proof             *types.Proof        `gorm:"embedded;embeddedPrefix:proof_"`
}

func (vm *VerifiableCredential) ToCoreType() *types.VerifiableCredential {
	return &types.VerifiableCredential{
		Context:           vm.Context,
		Type:              vm.Type,
		Issuer:            vm.Issuer,
		CredentialSubject: vm.CredentialSubject,
		ID:                vm.ID,
		IssuanceDate:      vm.IssuanceDate,
		ExpirationDate:    vm.ExpirationDate,
		CredentialSchema: converters.ConvertSliceCallback(
			vm.CredentialSchema,
			func(c *CredentialSchema) *types.CredentialSchema {
				return c.ToCoreType()
			},
		),
		Proof: vm.Proof,
	}
}

type CredentialSchema struct {
	ID                     string `gorm:"primaryKey"`
	VerifiableCredentialID string `gorm:"primaryKey"`
	Type                   string
}

func (c *CredentialSchema) ToCoreType() *types.CredentialSchema {
	return &types.CredentialSchema{
		Type: c.Type,
		ID:   c.ID,
	}
}

func newVerifiableCredentialModel(src *types.VerifiableCredential) *VerifiableCredential {
	return &VerifiableCredential{
		ID:                src.ID,
		CreatedAt:         time.Now().UTC(),
		Context:           src.Context,
		Type:              src.Type,
		Issuer:            src.Issuer,
		CredentialSubject: src.CredentialSubject,
		IssuanceDate:      src.IssuanceDate,
		ExpirationDate:    src.ExpirationDate,
		CredentialSchema:  converters.ConvertSliceCallback(src.CredentialSchema, newCredentialSchemaModel),
		Proof:             src.Proof,
	}
}

func newCredentialSchemaModel(src *types.CredentialSchema) *CredentialSchema {
	return &CredentialSchema{
		ID:   src.ID,
		Type: src.Type,
	}
}
