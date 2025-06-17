// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"
	"fmt"

	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	issuerverification "github.com/agntcy/identity/internal/core/issuer/verification"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/internal/pkg/oidc"
	"github.com/agntcy/identity/pkg/log"
)

// All IDP schemes supported by the ID generator.
// The ID generator creates IDs based on the proof and issuer information.
const (
	OktaScheme = "OKTA-"
	DuoScheme  = "DUO-"
	SelfScheme = "AGNTCY-"
)

type IDGenerator interface {
	GenerateFromProof(
		ctx context.Context,
		proof *vctypes.Proof,
	) (string, *issuertypes.Issuer, error)
}

type idGenerator struct {
	verifService issuerverification.Service
}

func NewIDGenerator(
	verifService issuerverification.Service,
) IDGenerator {
	return &idGenerator{
		verifService: verifService,
	}
}

func (g *idGenerator) GenerateFromProof(
	ctx context.Context,
	proof *vctypes.Proof,
) (string, *issuertypes.Issuer, error) {
	jwt, issuer, err := g.verifService.VerifyExistingIssuer(ctx, proof)
	if err != nil {
		return "", nil, err
	}

	var scheme string

	switch jwt.Provider {
	case oidc.OktaProviderName:
		scheme = OktaScheme
	case oidc.DuoProviderName:
		scheme = DuoScheme
	case oidc.SelfProviderName:
		scheme = SelfScheme
	default:
		return "", nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_UNKNOWN_IDP,
			"unknown JWT provider name",
			nil,
		)
	}

	log.Debug("Issuer is verified: ", issuer.Verified)
	log.Debug("JWT scheme: ", scheme)

	return fmt.Sprintf("%s%s", scheme, jwt.Claims.Subject), issuer, nil
}
