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
	issuercore "github.com/agntcy/identity/internal/core/issuer"
	vccore "github.com/agntcy/identity/internal/core/vc"
	"github.com/agntcy/identity/internal/core/vc/jose"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/log"
)

type VerifiableCredentialService interface {
	Publish(
		ctx context.Context,
		credential *vctypes.EnvelopedCredential,
		proof *vctypes.Proof,
	) error

	// Find the vcs by resolver metadata ID
	GetWellKnown(
		ctx context.Context,
		resolverMetadataID string,
	) (*[]*vctypes.EnvelopedCredential, error)
}

type verifiableCredentialService struct {
	verificationService core.VerificationService
	idRepository        idcore.IdRepository
	issuerRepository    issuercore.Repository
	vcRepository        vccore.Repository
	idGenerator         IDGenerator
}

func NewVerifiableCredentialService(
	verificationService core.VerificationService,
	idRepository idcore.IdRepository,
	issuerRepository issuercore.Repository,
	vcRepository vccore.Repository,
	idGenerator IDGenerator,
) VerifiableCredentialService {
	return &verifiableCredentialService{
		verificationService: verificationService,
		idRepository:        idRepository,
		issuerRepository:    issuerRepository,
		vcRepository:        vcRepository,
		idGenerator:         idGenerator,
	}
}

func (s *verifiableCredentialService) Publish(
	ctx context.Context,
	credential *vctypes.EnvelopedCredential,
	proof *vctypes.Proof,
) error {
	if credential.Value == "" {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_VALUE_FORMAT,
			"invalid credential envelope value",
			nil,
		)
	}

	id, _, err := s.idGenerator.GenerateFromProof(ctx, proof)
	if err != nil {
		return err
	}

	log.Debug("Resolving the ID into a ResolverMetadata")

	resolverMD, err := s.idRepository.ResolveID(ctx, id)
	if err != nil {
		if errors.Is(err, errcore.ErrResourceNotFound) {
			return errutil.ErrInfo(
				errtypes.ERROR_REASON_RESOLVER_METADATA_NOT_FOUND,
				fmt.Sprintf("could not resolve the ID (%s) to a resolver metadata", id),
				err,
			)
		}

		return errutil.ErrInfo(errtypes.ERROR_REASON_INTERNAL, "unexpected error", err)
	}

	var validatedVC *vctypes.VerifiableCredential

	log.Debug("Validating the verifiable credential")

	switch credential.EnvelopeType {
	case vctypes.CREDENTIAL_ENVELOPE_TYPE_EMBEDDED_PROOF:
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_TYPE,
			"credential envelope type not implemented yet",
			err,
		)
	case vctypes.CREDENTIAL_ENVELOPE_TYPE_JOSE:
		log.Debug("Verifying and parsing the JOSE Verifiable Credential")

		parsedVC, err := jose.Verify(resolverMD.GetJwks(), credential)
		if err != nil {
			return err
		}

		validatedVC = parsedVC
	default:
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_TYPE,
			"invalid credential envelope type",
			nil,
		)
	}

	log.Debug("Storing the Verifiable Credential")

	_, err = s.vcRepository.Create(ctx, validatedVC, id)
	if err != nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INTERNAL,
			"unable to store verifiable credential",
			err,
		)
	}

	return nil
}

func (s *verifiableCredentialService) GetWellKnown(
	ctx context.Context,
	resolverMetadataID string,
) (*[]*vctypes.EnvelopedCredential, error) {
	vcs, err := s.vcRepository.GetWellKnown(ctx, resolverMetadataID)
	if err != nil {
		return nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INTERNAL,
			"unable to retrieve well-known verifiable credentials",
			err,
		)
	}

	if vcs == nil || len(*vcs) == 0 {
		log.Debug("No well-known verifiable credentials found for resolver metadata ID ", resolverMetadataID)
		// Return an empty slice instead of nil to avoid nil pointer dereference
		emptySlice := []*vctypes.EnvelopedCredential{}

		return &emptySlice, nil
	}

	log.Debug("Found well-known verifiable credentials for resolver metadata ID ", resolverMetadataID)
	// Return the found verifiable credentials

	return vcs, nil
}
