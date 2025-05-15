package postgres

import (
	"context"
	"errors"

	issuercore "github.com/agntcy/identity/internal/core/issuer"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/db"
	"gorm.io/gorm"
)

type repository struct {
	dbContext db.Context
}

// NewIssuerRepository creates a new instance of the IssuerRepository
func NewRepository(dbContext db.Context) issuercore.Repository {
	return &repository{
		dbContext,
	}
}

// CreateIssuer creates a new Issuer
func (r *repository) CreateIssuer(
	ctx context.Context,
	issuer *issuertypes.Issuer,
) (*issuertypes.Issuer, error) {
	model := newIssuerModel(issuer)

	// Create the issuer
	inserted := r.dbContext.Client().Create(model)
	if inserted.Error != nil {
		return nil, errutil.Err(
			inserted.Error, "there was an error creating the issuer",
		)
	}

	return issuer, nil
}

func (r *repository) GetIssuer(ctx context.Context, commonName string) (*issuertypes.Issuer, error) {
	var isser Issuer
	result := r.dbContext.Client().First(&isser, commonName)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errutil.Err(
			result.Error, "there was an error fetching the issuer",
		)
	}

	return isser.ToCoreType(), nil
}
