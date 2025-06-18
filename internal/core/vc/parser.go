// Copyright 2025 Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package vc

import (
	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
	"github.com/agntcy/identity/internal/core/vc/jose"
	"github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/pkg/errutil"
	"github.com/agntcy/identity/pkg/log"
)

func ParseEnvelopedCredential(cred *types.EnvelopedCredential) (*types.VerifiableCredential, error) {
	var parsed *types.VerifiableCredential

	switch cred.EnvelopeType {
	case types.CREDENTIAL_ENVELOPE_TYPE_EMBEDDED_PROOF:
		return nil, credentialEnvelopeTypeNotImplErr()
	case types.CREDENTIAL_ENVELOPE_TYPE_JOSE:
		log.Debug("Parsing the JOSE Verifiable Credential")

		parsedVC, err := jose.Parse(cred)
		if err != nil {
			return nil, err
		}

		parsed = parsedVC
	default:
		return nil, invalidCredentialEnvelopeTypeErr()
	}

	return parsed, nil
}

func VerifyEnvelopedCredential(cred *types.EnvelopedCredential, jwks *idtypes.Jwks) error {
	switch cred.EnvelopeType {
	case types.CREDENTIAL_ENVELOPE_TYPE_EMBEDDED_PROOF:
		return credentialEnvelopeTypeNotImplErr()
	case types.CREDENTIAL_ENVELOPE_TYPE_JOSE:
		log.Debug("Verifying the JOSE Verifiable Credential")

		return jose.Verify(jwks, cred)
	default:
		return invalidCredentialEnvelopeTypeErr()
	}
}

func credentialEnvelopeTypeNotImplErr() error {
	return errutil.ErrInfo(
		errtypes.ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_TYPE,
		"credential envelope type not implemented yet",
		nil,
	)
}

func invalidCredentialEnvelopeTypeErr() error {
	return errutil.ErrInfo(
		errtypes.ERROR_REASON_INVALID_CREDENTIAL_ENVELOPE_TYPE,
		"invalid credential envelope type",
		nil,
	)
}
