// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	errcore "github.com/agntcy/identity/internal/core/errors"
	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	issuercore "github.com/agntcy/identity/internal/core/issuer"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/agntcy/identity/pkg/log"
)

type IDGenerator interface {
	GenerateFromProof(
		ctx context.Context,
		proof *vctypes.Proof,
	) (string, *issuertypes.Issuer, error)
}

type idGenerator struct {
	oidcParser       oidc.Parser
	issuerRepository issuercore.Repository
}

func NewIDGenerator(oidcParser oidc.Parser, issuerRepository issuercore.Repository) IDGenerator {
	return &idGenerator{
		oidcParser:       oidcParser,
		issuerRepository: issuerRepository,
	}
}

func (g *idGenerator) GenerateFromProof(
	ctx context.Context,
	proof *vctypes.Proof,
) (string, *issuertypes.Issuer, error) {
	if proof == nil {
		return "", nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_IDP_REQUIRED,
			"issuer without external IdP is not implemented",
			nil,
		)
	}

	log.Debug("Verifying the proof ", proof.ProofValue)

	if proof.IsJWT() {
		jwt, err := g.oidcParser.ParseJwt(ctx, &proof.ProofValue)
		if err != nil {
			return "", nil, errutil.ErrInfo(errtypes.ERROR_REASON_INVALID_PROOF, err.Error(), err)
		}

		// Extract the hostname from the issuer
		u, err := url.Parse(jwt.Claims.Issuer)
		if err != nil {
			return "", nil, errutil.ErrInfo(errtypes.ERROR_REASON_INVALID_PROOF, err.Error(), err)
		}

		issuer, err := g.getIssuer(ctx, u.Hostname())
		if err != nil {
			return "", nil, err
		}

		var scheme string

		switch jwt.Provider {
		case oidc.OktaProviderName:
			scheme = "OKTA"
		case oidc.DuoProviderName:
			scheme = "DUO"
		default:
			return "", nil, errutil.ErrInfo(errtypes.ERROR_REASON_UNKNOWN_IDP, "unknown JWT provider name", nil)
		}

		return fmt.Sprintf("%s-%s", scheme, jwt.Claims.Subject), issuer, nil
	}

	return "", nil, errutil.ErrInfo(
		errtypes.ERROR_REASON_UNSUPPORTED_PROOF,
		fmt.Sprintf("unsupported proof type: %s", proof.Type),
		nil,
	)
}

func (g *idGenerator) getIssuer(ctx context.Context, commonName string) (*issuertypes.Issuer, error) {
	issuer, err := g.issuerRepository.GetIssuer(ctx, commonName)
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
