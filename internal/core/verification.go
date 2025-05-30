// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"context"
	"fmt"
	"net/url"

	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/agntcy/identity/pkg/log"
)

// The VerificationService interface defines the core methods for
// common name verification
type VerificationService interface {
	VerifyCommonName(ctx context.Context, commonName *string, proof *vctypes.Proof) error
	VerifyProof(ctx context.Context, proof *vctypes.Proof) (string, string, error)
}

// The verificationService struct implements the VerificationService interface
type verificationService struct {
	oidcParser oidc.Parser
}

// NewVerificationService creates a new instance of the VerificationService
func NewVerificationService(oidcParser oidc.Parser) VerificationService {
	return &verificationService{
		oidcParser,
	}
}

// VerifyCommonName verifies the common name against the proof
// by checking if the common name is the same as the proof's issuer's hostname
func (v *verificationService) VerifyCommonName(
	ctx context.Context,
	commonName *string,
	proof *vctypes.Proof,
) error {
	// Verify the proof and get the subject and issuer
	issuer, _, err := v.VerifyProof(ctx, proof)
	if err != nil {
		return err
	}

	log.Debug("Verifying common name:", *commonName)

	// Extract the hostname from the issuer
	u, err := url.Parse(issuer)
	if err != nil {
		return err
	}

	log.Debug("Issuer hostname:", u.Hostname())

	// Verify common name is the same as the issuer's hostname
	if u.Hostname() != *commonName {
		return errutil.Err(nil, "common name does not match issuer")
	}

	log.Debug("Common name verified successfully")

	return nil
}

// VerifyProof verifies the proof and returns the subject and issuer
// based on the proof type
func (v *verificationService) VerifyProof(
	ctx context.Context,
	proof *vctypes.Proof,
) (string, string, error) {
	// Validate the proof
	if proof == nil {
		return "", "", errutil.Err(nil, "proof is empty")
	}

	log.Debug("Verifying proof of type", proof.Type)

	// Check the proof type
	if proof.IsJWT() {
		// Verify the JWT proof
		jwt, err := v.oidcParser.ParseJwt(ctx, &proof.ProofValue)
		if err != nil {
			return "", "", err
		}

		// Return the issuer and subject
		return jwt.Claims.Issuer, jwt.Claims.Subject, nil
	}

	return "", "", errutil.Err(nil, fmt.Sprintf("unsupported proof type '%s'", proof.Type))
}
