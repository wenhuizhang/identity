// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package verification

import (
	"context"
	"fmt"

	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/internal/pkg/httputil"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/agntcy/identity/pkg/log"
)

// The VerificationService interface defines the core methods for
// common name verification
type Service interface {
	Verify(ctx context.Context, issuer *issuertypes.Issuer, proof *vctypes.Proof) (bool, error)
}

// The verificationService struct implements the VerificationService interface
type service struct {
	oidcParser oidc.Parser
}

// NewVerificationService creates a new instance of the VerificationService
func NewService(oidcParser oidc.Parser) Service {
	return &service{
		oidcParser,
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

	log.Debug("Verifying common name:", issuer.CommonName)

	// Verify the proof and get the subject and issuer
	err := v.verifyProof(ctx, issuer, proof)
	if err != nil {
		return false, err
	}

	log.Debug("Common name verified successfully")

	return true, nil
}

// VerifyProof verifies the proof for the issuer by checking the proof type
func (v *service) verifyProof(
	ctx context.Context,
	issuer *issuertypes.Issuer,
	proof *vctypes.Proof,
) error {
	// Validate the proof
	if proof == nil {
		return errutil.Err(nil, "proof is empty")
	}

	log.Debug("Verifying proof of type: ", proof.Type)

	// Check the proof type
	if proof.IsJWT() {
		// Verify the JWT proof
		jwt, err := v.oidcParser.ParseJwt(ctx, &proof.ProofValue, issuer.PublicKey.Jwks().String())
		if err != nil {
			return err
		}

		// Verify common name is the same as the issuer's hostname
		if httputil.Hostname(jwt.Claims.Issuer) != issuer.CommonName {
			return errutil.Err(nil, "common name does not match issuer")
		}

		return nil
	}

	return errutil.Err(nil, fmt.Sprintf("unsupported proof type '%s'", proof.Type))
}
