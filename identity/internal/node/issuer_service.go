// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"

	"github.com/agntcy/identity/internal/core"
	"github.com/agntcy/identity/internal/core/issuer"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/internal/pkg/grpcutil"
)

// The IssuerService interface defines the Node methods for Issuers
type IssuerService interface {
	// Register a new Issuer
	// In case of external IdPs provide a proof of ownership
	Register(ctx context.Context, issuer *issuertypes.Issuer, proof *vctypes.Proof) (*string, error)
}

// The issuerService struct implements the IssuerService interface
type issuerService struct {
	issuerRepository   issuer.Repository
	verficationService core.VerificationService
}

// NewIssuerService creates a new instance of the IssuerService
func NewIssuerService(
	issuerRepository issuer.Repository,
	verficationService core.VerificationService,
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
) (*string, error) {
	// Validate the issuer
	if issuer == nil || issuer.CommonName == "" {
		return nil, grpcutil.BadRequestError(
			errutil.Err(nil, "issuer is empty or has invalid common name"),
		)
	}

	// Validate the public key

	// Verify the issuer's common name
	// Validate the proof exists
	if proof == nil {
		// In case of external IdPs, the proof is nil
		// This service should return an actionable URI
		// to the caller to finalize the registration
		// This is currently not supported
		return nil, grpcutil.UnimplementedError(
			errutil.Err(nil, "issuer without external IdP is not implemented"),
		)
	} else {
		verificationErr := i.verficationService.VerifyCommonName(
			ctx,
			&issuer.CommonName,
			proof,
		)

		if verificationErr != nil {
			return nil, grpcutil.BadRequestError(
				errutil.Err(verificationErr, "failed to verify common name"),
			)
		}
	}

	// Save the issuer in the database
	_, repositoryErr := i.issuerRepository.CreateIssuer(
		ctx,
		issuer,
	)
	if repositoryErr != nil {
		return nil, grpcutil.InternalError(
			errutil.Err(repositoryErr, ""),
		)
	}

	//nolint:nilnil // Ignore linting for nil return, means no action uri
	return nil, nil
}
