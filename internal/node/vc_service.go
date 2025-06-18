// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node

import (
	"context"
	"errors"
	"fmt"
	"strings"

	errcore "github.com/agntcy/identity/internal/core/errors"
	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	idcore "github.com/agntcy/identity/internal/core/id"
	issuerverification "github.com/agntcy/identity/internal/core/issuer/verification"
	vccore "github.com/agntcy/identity/internal/core/vc"
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
	GetVcs(
		ctx context.Context,
		resolverMetadataID string,
	) ([]*vctypes.EnvelopedCredential, error)

	// Parse and verify a Verifiable Credential
	Verify(
		ctx context.Context,
		credential *vctypes.EnvelopedCredential,
	) error
}

type verifiableCredentialService struct {
	idRepository idcore.IdRepository
	verifService issuerverification.Service
	vcRepository vccore.Repository
}

func NewVerifiableCredentialService(
	idRepository idcore.IdRepository,
	verifService issuerverification.Service,
	vcRepository vccore.Repository,
) VerifiableCredentialService {
	return &verifiableCredentialService{
		idRepository: idRepository,
		verifService: verifService,
		vcRepository: vcRepository,
	}
}

func (s *verifiableCredentialService) Publish(
	ctx context.Context,
	credential *vctypes.EnvelopedCredential,
	proof *vctypes.Proof,
) error {
	parsedVC, err := s.verifyEnvelopedCredential(ctx, credential)
	if err != nil {
		return err
	}

	id, ok := parsedVC.GetDID()
	if !ok {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_VERIFIABLE_CREDENTIAL,
			"unable to find the ID inside the CredentialSubject",
			nil,
		)
	}

	log.Debug("Validating the authentication proof")

	issuerVerification, err := s.verifService.VerifyExistingIssuer(ctx, proof)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(id, issuerVerification.Subject) {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_VERIFIABLE_CREDENTIAL,
			"the ID in the Verifiable Credential does not match the ID in the proof",
			nil,
		)
	}

	log.Debug("Storing the Verifiable Credential")

	_, err = s.vcRepository.Create(ctx, parsedVC, id)
	if err != nil {
		return errutil.ErrInfo(
			errtypes.ERROR_REASON_INTERNAL,
			"unable to store verifiable credential",
			err,
		)
	}

	return nil
}

func (s *verifiableCredentialService) GetVcs(
	ctx context.Context,
	resolverMetadataID string,
) ([]*vctypes.EnvelopedCredential, error) {
	log.Debug(
		"Retrieving well-known verifiable credentials for resolver metadata ID: ",
		resolverMetadataID,
	)

	vcs, err := s.vcRepository.GetByResolverMetadata(ctx, resolverMetadataID)
	if err != nil {
		if errors.Is(err, errcore.ErrResourceNotFound) {
			log.Debug(
				"No well-known verifiable credentials found for resolver metadata ID ",
				resolverMetadataID,
			)

			return []*vctypes.EnvelopedCredential{}, nil
		}

		return nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INTERNAL,
			"unable to retrieve well-known verifiable credentials",
			err,
		)
	}

	log.Debug(
		"Found well-known verifiable credentials for resolver metadata ID ",
		resolverMetadataID,
	)

	var envelopedCredentials []*vctypes.EnvelopedCredential

	for _, cred := range vcs {
		if cred.Proof == nil {
			log.Debug("Skipping credential with empty proof for ID: ", cred.ID)
			continue
		}

		if cred.Proof.Type == "" {
			log.Debug("Skipping credential with empty proof type for ID: ", cred.ID)
			continue
		}

		if cred.Proof.ProofValue == "" {
			log.Debug("Skipping credential with empty proof value for ID: ", cred.ID)
			continue
		}

		switch cred.Proof.Type {
		case "JWT":
			envelopedCredentials = append(envelopedCredentials, &vctypes.EnvelopedCredential{
				EnvelopeType: vctypes.CREDENTIAL_ENVELOPE_TYPE_JOSE,
				Value:        cred.Proof.ProofValue,
			})
		default:
			log.Debug(
				"Skipping credential with unsupported proof type: ",
				cred.Proof.Type,
				" for ID: ",
				cred.ID,
			)
		}
	}

	return envelopedCredentials, nil
}

// Verify an existing Verifiable Credential
func (s *verifiableCredentialService) Verify(
	ctx context.Context,
	credential *vctypes.EnvelopedCredential,
) error {
	_, err := s.verifyEnvelopedCredential(ctx, credential)
	return err
}

func (s *verifiableCredentialService) verifyEnvelopedCredential(
	ctx context.Context,
	credential *vctypes.EnvelopedCredential,
) (*vctypes.VerifiableCredential, error) {
	if credential.Value == "" {
		return nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_VALUE_FORMAT,
			"invalid credential envelope value",
			nil,
		)
	}

	parsedVC, err := vccore.ParseEnvelopedCredential(credential)
	if err != nil {
		return nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_VALUE_FORMAT,
			"invalid credential envelope value",
			err,
		)
	}

	id, ok := parsedVC.GetDID()
	if !ok {
		return nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_VERIFIABLE_CREDENTIAL,
			"unable to find the ID inside the CredentialSubject",
			nil,
		)
	}

	log.Debug("Resolving the ID into a ResolverMetadata")

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

	log.Debug("Validating the verifiable credential")

	err = vccore.VerifyEnvelopedCredential(credential, resolverMD.GetJwks())
	if err != nil {
		return nil, errutil.ErrInfo(
			errtypes.ERROR_REASON_INVALID_VERIFIABLE_CREDENTIAL,
			"unable to verify verifiable credential",
			err,
		)
	}

	return parsedVC, nil
}
