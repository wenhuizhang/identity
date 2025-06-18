// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"
	"errors"

	errcore "github.com/agntcy/identity/internal/core/errors"
	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	issuercore "github.com/agntcy/identity/internal/core/issuer"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	"github.com/agntcy/identity/internal/core/issuer/verification"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/joseutil"
	"github.com/agntcy/identity/pkg/jwk"
	"github.com/agntcy/identity/pkg/oidc"
)

// The IssuerService interface defines the Node methods for Issuers
type IssuerService interface {
	// Register a new Issuer
	// In case of external IdPs provide a proof of ownership
	Register(ctx context.Context, issuer *issuertypes.Issuer, proof *vctypes.Proof) error

	// Find the issuer by common name
	// Return the public keys of the Issuer
	GetJwks(ctx context.Context, commonName string) (*jwk.Jwks, error)
}

// The issuerService struct implements the IssuerService interface
type issuerService struct {
	issuerRepository   issuercore.Repository
	verficationService verification.Service
}

// NewIssuerService creates a new instance of the IssuerService
func NewIssuerService(
	issuerRepository issuercore.Repository,
	verficationService verification.Service,
) IssuerService {
	return &issuerService{
		issuerRepository,
		verficationService,
	}
}

// Register a new Issuers
// In case of external IdPs provide a proof of ownership
func (i *issuerService) Register(
	ctx context.Context,
	issuer *issuertypes.Issuer,
	proof *vctypes.Proof,
) error {
	if issuer == nil || issuer.ValidateCommonName() != nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_ISSUER,
			"issuer is empty or has invalid common name",
			nil,
		)
	}

	err := joseutil.ValidatePubKey(issuer.PublicKey)
	if err != nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_ISSUER,
			"issuer has invalid public key",
			nil,
		)
	}

	// Verify the issuer's common name
	// Validate the proof exists
	// If the proof is self signed, the issuer will not be verified
	verifRes, err := i.verficationService.Verify(
		ctx,
		issuer,
		proof,
	)
	if err != nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_ISSUER,
			"failed to verify common name",
			err,
		)
	}

	// Check if issuer already exists
	existingIssuer, _ := i.issuerRepository.GetIssuer(
		ctx,
		issuer.CommonName,
	)
	if existingIssuer != nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_ISSUER,
			"issuer already exists",
			nil,
		)
	}

	// Set the verified status of the issuer
	// When we implement the verification process we need to change this code
	issuer.Verified = verifRes.Verified

	if verifRes.Provider == oidc.SelfProviderName {
		issuer.AuthType = issuertypes.ISSUER_AUTH_TYPE_SELF
	} else {
		issuer.AuthType = issuertypes.ISSUER_AUTH_TYPE_IDP
	}

	// Save the issuer in the database
	_, repositoryErr := i.issuerRepository.CreateIssuer(
		ctx,
		issuer,
	)
	if repositoryErr != nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INTERNAL,
			"unexpected error",
			repositoryErr,
		)
	}

	return nil
}

// GetJwks returns the public keys of the Issuers
// The common name is used to find the Issuers
func (i *issuerService) GetJwks(
	ctx context.Context,
	commonName string,
) (*jwk.Jwks, error) {
	// Validate the common name
	if commonName == "" {
		return nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_ISSUER,
			"issuer common name is empty",
			nil,
		)
	}

	// Find the issuer by common name
	issuer, err := i.issuerRepository.GetIssuer(ctx, commonName)
	if err != nil {
		if errors.Is(err, errcore.ErrResourceNotFound) {
			return nil, errutil.ErrInfo(
				errtypes.ERROR_REASON_ISSUER_NOT_REGISTERED,
				"the issuer is not registered",
				nil,
			)
		}

		return nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INTERNAL,
			"unexpected error",
			err,
		)
	}

	// Return the public keys of the Issuer
	return &jwk.Jwks{
		Keys: []*jwk.Jwk{
			issuer.PublicKey,
		},
	}, nil
}
