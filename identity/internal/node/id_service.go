// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"
	"fmt"

	"github.com/agntcy/identity/internal/core"
	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	idcore "github.com/agntcy/identity/internal/core/id"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
	issuercore "github.com/agntcy/identity/internal/core/issuer"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/log"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IdService interface {
	Generate(ctx context.Context, issuer *issuertypes.Issuer, proof *vctypes.Proof) (*idtypes.ResolverMetadata, error)
}

type idService struct {
	verificationService core.VerificationService
	idRepository        idcore.IdRepository
	issuerRepository    issuercore.Repository
}

func NewIdService(
	verificationService core.VerificationService,
	idRepository idcore.IdRepository,
	issuerRepository issuercore.Repository,
) IdService {
	return &idService{
		verificationService: verificationService,
		idRepository:        idRepository,
		issuerRepository:    issuerRepository,
	}
}

func (s *idService) Generate(
	ctx context.Context,
	issuer *issuertypes.Issuer,
	proof *vctypes.Proof,
) (*idtypes.ResolverMetadata, error) {
	if proof == nil {
		return nil, errutil.ErrInfo(errtypes.ERROR_REASON_IDP_REQUIRED, "issuer without external IdP is not implemented")
	}

	_, sub, err := s.verificationService.VerifyProof(ctx, proof)
	if err != nil {
		return nil, errutil.ErrInfo(errtypes.ERROR_REASON_INVALID_PROOF, err.Error())
	}

	iss, err := s.issuerRepository.GetIssuer(ctx, issuer.CommonName)
	if err != nil {
		// TODO: handle error
		return nil, errutil.ErrInfo(errtypes.ERROR_REASON_ISSUER_NOT_REGISTERED, fmt.Sprintf("the issuer %s is not registered", issuer.CommonName))
	}

	err = s.verificationService.VerifyCommonName(ctx, &iss.CommonName, proof)
	if err != nil {
		return nil, errutil.ErrInfo(errtypes.ERROR_REASON_INVALID_ISSUER, "failed to verify issuer common name")
	}

	// TODO: Check when it's Okta or Duo or other IdP
	id := fmt.Sprintf("DUO-%s", sub)
	keyID := fmt.Sprintf("%s#%s", id, uuid.NewString())

	resolverMetadata := &idtypes.ResolverMetadata{
		ID: id,
		VerificationMethod: []idtypes.VerificationMethod{
			{
				ID:           keyID,
				PublicKeyJwk: iss.PublicKey,
			},
		},
		AssertionMethod: []string{keyID},
		Service: []idtypes.Service{
			{
				ServiceEndpoint: []string{"The URL for duo/okta tenant"},
			},
		},
	}

	_, err = s.idRepository.CreateID(ctx, resolverMetadata)
	if err != nil {
		log.WithFields(logrus.Fields{log.ErrorField: err}).Error()
		return nil, errutil.ErrInfo(errtypes.ERROR_READON_INTERNAL, "unable to store the resolver metadata")
	}

	return resolverMetadata, nil
}
