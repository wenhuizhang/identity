// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"
	"errors"
	"fmt"

	"github.com/agntcy/identity/internal/core"
	errcore "github.com/agntcy/identity/internal/core/errors"
	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	idcore "github.com/agntcy/identity/internal/core/id"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
	issuercore "github.com/agntcy/identity/internal/core/issuer"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/log"
	"github.com/google/uuid"
)

type IdService interface {
	Generate(
		ctx context.Context,
		issuer *issuertypes.Issuer,
		proof *vctypes.Proof,
	) (*idtypes.ResolverMetadata, error)
	Resolve(
		ctx context.Context,
		id string,
	) (*idtypes.ResolverMetadata, error)
}

type idService struct {
	verificationService core.VerificationService
	idRepository        idcore.IdRepository
	issuerRepository    issuercore.Repository
	idGenerator         IDGenerator
}

func NewIdService(
	verificationService core.VerificationService,
	idRepository idcore.IdRepository,
	issuerRepository issuercore.Repository,
	idGenerator IDGenerator,
) IdService {
	return &idService{
		verificationService: verificationService,
		idRepository:        idRepository,
		issuerRepository:    issuerRepository,
		idGenerator:         idGenerator,
	}
}

func (s *idService) Generate(
	ctx context.Context,
	issuer *issuertypes.Issuer,
	proof *vctypes.Proof,
) (*idtypes.ResolverMetadata, error) {
	id, storedIss, err := s.idGenerator.GenerateFromProof(ctx, proof)
	if err != nil {
		return nil, err
	}

	log.Debug("ID generated ", id)

	err = s.verifyIssuer(issuer, storedIss)
	if err != nil {
		return nil, err
	}

	log.Debug("Generating a ResolverMetadata")

	keyID := fmt.Sprintf("%s#%s", id, uuid.NewString())

	resolverMetadata := &idtypes.ResolverMetadata{
		ID: id,
		VerificationMethod: []*idtypes.VerificationMethod{
			{
				ID:           keyID,
				PublicKeyJwk: storedIss.PublicKey, // Should we compare issuer.PK with iss.PK?
			},
		},
		AssertionMethod: []string{keyID},
		Service: []*idtypes.Service{
			{
				// This works for now, but not when we add support for non IdP backed Issuers
				ServiceEndpoint: []string{fmt.Sprintf("https://%s", storedIss.CommonName)},
			},
		},
	}

	log.Debug("Storing the ResolverMetadata")

	_, err = s.idRepository.CreateID(ctx, resolverMetadata)
	if err != nil {
		return nil, errutil.ErrInfo(errtypes.ERROR_REASON_INTERNAL, "unable to store the resolver metadata", err)
	}

	return resolverMetadata, nil
}

func (s *idService) verifyIssuer(
	input *issuertypes.Issuer,
	existing *issuertypes.Issuer,
) error {
	if input == nil || existing == nil || input.CommonName != existing.CommonName {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_ISSUER,
			"failed to verify issuer common name",
			nil,
		)
	}

	return nil
}

func (s *idService) Resolve(ctx context.Context, id string) (*idtypes.ResolverMetadata, error) {
	resolverMD, err := s.idRepository.ResolveID(ctx, id)
	if err != nil {
		if errors.Is(err, errcore.ErrResourceNotFound) {
			return nil, errutil.ErrInfo(
				errtypes.ERROR_REASON_RESOLVER_METADATA_NOT_FOUND,
				fmt.Sprintf("could not resolve the ID (%s) to a resolver metadata", id),
				err,
			)
		}

		return nil, errutil.ErrInfo(errtypes.ERROR_REASON_INTERNAL, "unexpected error", err)
	}

	return resolverMD, nil
}
