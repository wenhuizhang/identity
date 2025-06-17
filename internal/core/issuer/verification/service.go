// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verification

import (
	"context"
	"errors"
	"fmt"

	errcore "github.com/agntcy/identity/internal/core/errors"
	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	issuercore "github.com/agntcy/identity/internal/core/issuer"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/agntcy/identity/pkg/log"
)

// The VerificationService interface defines the core methods for
// common name verification
type Service interface {
	Verify(
		ctx context.Context,
		issuer *issuertypes.Issuer,
		proof *vctypes.Proof,
	) (bool, error)
	VerifyExistingIssuer(
		ctx context.Context,
		proof *vctypes.Proof,
	) (*oidc.ParsedJWT, *issuertypes.Issuer, error)
}

// The verificationService struct implements the VerificationService interface
type service struct {
	oidcParser oidc.Parser
	repository issuercore.Repository
}

// NewVerificationService creates a new instance of the VerificationService
func NewService(oidcParser oidc.Parser, repository issuercore.Repository) Service {
	return &service{
		oidcParser,
		repository,
	}
}

// Verify verifies the issuer's common name against the proof
// In case the proof is self provided, the issuer will be unverified
func (v *service) Verify(
	ctx context.Context,
	issuer *issuertypes.Issuer,
	proof *vctypes.Proof,
) (bool, error) {
	// The proof is required to verify the issuer's common name
	if proof == nil {
		return false, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_PROOF,
			"a proof is required to verify the issuer's common name",
			nil,
		)
	}

	// Verify the proof
	// If the proof is self provided, the issuer will be unverified
	verified, err := v.verifyProof(ctx, issuer, proof)
	if err != nil {
		return false, err
	}

	return verified, nil
}

// VerifyProof verifies the proof for the issuer by checking the proof type
func (v *service) verifyProof(
	ctx context.Context,
	issuer *issuertypes.Issuer,
	proof *vctypes.Proof,
) (bool, error) {
	// Validate the proof
	if proof == nil {
		return false, errutil.Err(nil, "proof is empty")
	}

	log.Debug("Verifying proof: ", proof.ProofValue, " of type: ", proof.Type)

	// Check the proof type
	if !proof.IsJWT() {
		return false, errutil.Err(nil, fmt.Sprintf("unsupported proof type '%s'", proof.Type))
	}

	// Parse JWT to extract the common name and issuer information
	jwt, err := v.oidcParser.ParseAndVerifyJwt(
		ctx,
		&proof.ProofValue,
		issuer.PublicKey.Jwks().String(),
	)
	if err != nil {
		return false, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_PROOF,
			"failed to parse and verify JWT",
			err,
		)
	}

	log.Debug("Verifying common name:", issuer.CommonName)

	// Verify common name is the same as the issuer's hostname
	if jwt.CommonName != issuer.CommonName {
		return false, errutil.Err(nil, "common name does not match issuer")
	}

	log.Debug("Common name verified successfully")

	return jwt.Verified, nil
}

func (v *service) VerifyExistingIssuer(
	ctx context.Context,
	proof *vctypes.Proof,
) (*oidc.ParsedJWT, *issuertypes.Issuer, error) {
	// Validate the proof
	if proof == nil {
		return nil, nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_PROOF,
			"proof is empty",
			nil,
		)
	}

	log.Debug("Verifying proof: ", proof.ProofValue, " of type: ", proof.Type)

	if !proof.IsJWT() {
		return nil, nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_UNSUPPORTED_PROOF,
			fmt.Sprintf("unsupported proof type: %s", proof.Type),
			nil,
		)
	}

	// Parse JWT to extract the common name and issuer information
	parsedJWT, err := v.oidcParser.ParseJwt(ctx, &proof.ProofValue)
	if err != nil {
		return nil, nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_PROOF,
			err.Error(),
			err,
		)
	}

	issuer, err := v.getIssuer(ctx, parsedJWT.CommonName)
	if err != nil {
		return nil, nil, err
	}

	// Verify the JWT signature
	err = v.oidcParser.VerifyJwt(ctx, parsedJWT, issuer.PublicKey.Jwks().String())
	if err != nil {
		return nil, nil, errutil.ErrInfo(errtypes.ERROR_REASON_INVALID_PROOF, err.Error(), err)
	}

	// If the issuer is verified
	// we require a valid proof from the IdP
	if issuer.Verified && parsedJWT.Provider == oidc.SelfProviderName {
		return nil, nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_IDP_REQUIRED,
			"the issuer is verified so the proof must be from an IdP",
			nil,
		)
	}

	return parsedJWT, issuer, nil
}

func (v *service) getIssuer(
	ctx context.Context,
	commonName string,
) (*issuertypes.Issuer, error) {
	issuer, err := v.repository.GetIssuer(ctx, commonName)
	if err != nil {
		if errors.Is(err, errcore.ErrResourceNotFound) {
			return nil, errutil.ErrInfo(
				errtypes.ERROR_REASON_ISSUER_NOT_REGISTERED,
				fmt.Sprintf("the issuer %s is not registered", commonName),
				err,
			)
		}

		return nil, errutil.ErrInfo(errtypes.ERROR_REASON_INTERNAL, "unexpected error", err)
	}

	return issuer, nil
}
