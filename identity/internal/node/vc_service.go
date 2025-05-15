// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"
	"fmt"

	"github.com/agntcy/identity/internal/core"
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
	Publish(ctx context.Context, credential *vctypes.EnvelopedCredential, proof *vctypes.Proof) error
}

type verifiableCredentialService struct {
	verificationService core.VerificationService
	idRepository        idcore.IdRepository
	issuerRepository    issuercore.Repository
	vcRepository        vccore.Repository
}

func NewVerifiableCredentialService(
	verificationService core.VerificationService,
	idRepository idcore.IdRepository,
	issuerRepository issuercore.Repository,
	vcRepository vccore.Repository,
) VerifiableCredentialService {
	return &verifiableCredentialService{
		verificationService: verificationService,
		idRepository:        idRepository,
		issuerRepository:    issuerRepository,
		vcRepository:        vcRepository,
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

	log.Debug("Verifying the ID proof and the issuer")

	if proof == nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_IDP_REQUIRED,
			"issuer without external IdP is not implemented",
			nil,
		)
	}

	iss, sub, err := s.verificationService.VerifyProof(ctx, proof)
	if err != nil {
		return errutil.ErrInfo(errtypes.ERROR_REASON_INVALID_PROOF, err.Error(), err)
	}

	storedIss, err := s.issuerRepository.GetIssuer(ctx, iss)
	if err != nil {
		return errutil.ErrInfo(errtypes.ERROR_REASON_INTERNAL, "unexpected error", err)
	} else if storedIss == nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_ISSUER_NOT_REGISTERED,
			fmt.Sprintf("the issuer %s is not registered", iss),
			err,
		)
	}

	// TODO: revisit this line after working on a better ID generation
	id := fmt.Sprintf("DUO-%s", sub)

	log.Debug("Resolving the ID into a ResolverMetadata")

	resolverMD, err := s.idRepository.ResolveID(ctx, id)
	if err != nil {
		return errutil.ErrInfo(errtypes.ERROR_REASON_INTERNAL, "unexpected error", err)
	} else if resolverMD == nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_RESOLVER_METADATA_NOT_FOUND,
			fmt.Sprintf("could not resolve the ID (%s) to a resolver metadata", id),
			err,
		)
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

	_, err = s.vcRepository.Create(ctx, validatedVC)
	if err != nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INTERNAL,
			"unable to store verifiable credential",
			err,
		)
	}

	return nil
}
