// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node_test

import (
	"context"
	"errors"
	"testing"

	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	idtesting "github.com/agntcy/identity/internal/core/id/testing"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
	issuertesting "github.com/agntcy/identity/internal/core/issuer/testing"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	coretesting "github.com/agntcy/identity/internal/core/testing"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/node"
	"github.com/agntcy/identity/internal/pkg/oidc"
	oidctesting "github.com/agntcy/identity/internal/pkg/oidc/testing"
	"github.com/stretchr/testify/assert"
)

func TestGenerateID_Should_Not_Return_Errors(t *testing.T) {
	t.Parallel()

	verficationSrv := coretesting.NewFakeTruthyVerificationService()
	idRepo := idtesting.NewFakeIdRepository()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	jwt := &oidc.ParsedJWT{
		Provider: oidc.DuoProviderName,
		Claims: &oidc.Claims{
			Issuer:  "http://" + coretesting.ValidProofIssuer,
			Subject: "test",
		},
	}
	idGen := node.NewIDGenerator(
		oidctesting.NewFakeParser(jwt, nil),
		issuerRepo,
	)
	sut := node.NewIdService(verficationSrv, idRepo, issuerRepo, idGen)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}
	_, _ = issuerRepo.CreateIssuer(context.Background(), issuer)

	md, err := sut.Generate(context.Background(), issuer, &vctypes.Proof{Type: "JWT"})

	assert.NoError(t, err)
	assert.Equal(t, "DUO-test", md.ID)
}

func TestGenerateID_Should_Return_Idp_Required_Error(t *testing.T) {
	t.Parallel()

	oidcParser := oidctesting.NewFakeParser(&oidc.ParsedJWT{}, nil)
	idGen := node.NewIDGenerator(oidcParser, nil)
	sut := node.NewIdService(nil, nil, nil, idGen)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}

	_, err := sut.Generate(context.Background(), issuer, nil)

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_IDP_REQUIRED)
}

func TestGenerateID_Should_Return_Invalid_Proof_Error(t *testing.T) {
	t.Parallel()

	idGen := node.NewIDGenerator(oidctesting.NewFakeParser(nil, errors.New("")), nil)
	verficationSrv := coretesting.NewFalsyProofVerificationServiceStub()
	sut := node.NewIdService(verficationSrv, nil, nil, idGen)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}

	_, err := sut.Generate(context.Background(), issuer, &vctypes.Proof{Type: "JWT"})

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_INVALID_PROOF)
}

func TestGenerateID_Should_Return_Invalid_Issuer_Error(t *testing.T) {
	t.Parallel()

	verficationSrv := coretesting.NewFalsyCommonNameVerificationServiceStub()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	jwt := &oidc.ParsedJWT{
		Provider: oidc.DuoProviderName,
		Claims: &oidc.Claims{
			Issuer:  "http://" + coretesting.ValidProofIssuer,
			Subject: "test",
		},
	}
	idGen := node.NewIDGenerator(
		oidctesting.NewFakeParser(jwt, nil),
		issuerRepo,
	)
	idRepo := idtesting.NewFakeIdRepository()
	sut := node.NewIdService(verficationSrv, idRepo, issuerRepo, idGen)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}
	_, _ = issuerRepo.CreateIssuer(context.Background(), issuer)
	invalidIssuer := &issuertypes.Issuer{
		CommonName: "INVALID",
	}

	_, err := sut.Generate(context.Background(), invalidIssuer, &vctypes.Proof{Type: "JWT"})

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_INVALID_ISSUER)
}

func TestGenerateID_Should_Return_Unregistred_Issuer_Error(t *testing.T) {
	t.Parallel()

	verficationSrv := coretesting.NewFakeTruthyVerificationService()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	jwt := &oidc.ParsedJWT{
		Provider: oidc.DuoProviderName,
		Claims:   &oidc.Claims{Subject: "test"},
	}
	idGen := node.NewIDGenerator(
		oidctesting.NewFakeParser(jwt, nil),
		issuerRepo,
	)
	sut := node.NewIdService(verficationSrv, nil, issuerRepo, idGen)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}

	_, err := sut.Generate(context.Background(), issuer, &vctypes.Proof{Type: "JWT"})

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_ISSUER_NOT_REGISTERED)
}

func TestGenerateID_Should_Return_ID_Already_Exists_Error(t *testing.T) {
	t.Parallel()

	idRepo := idtesting.NewFakeIdRepository()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	claims := &oidc.Claims{
		Issuer:  "http://" + coretesting.ValidProofIssuer,
		Subject: "test",
	}
	jwt := &oidc.ParsedJWT{
		Provider: oidc.DuoProviderName,
		Claims:   claims,
	}
	idGen := node.NewIDGenerator(
		oidctesting.NewFakeParser(jwt, nil),
		issuerRepo,
	)
	sut := node.NewIdService(nil, idRepo, nil, idGen)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}
	_, _ = issuerRepo.CreateIssuer(context.Background(), issuer)
	_, _ = idRepo.CreateID(context.Background(), &idtypes.ResolverMetadata{ID: "DUO-" + claims.Subject}, issuer)

	_, err := sut.Generate(context.Background(), nil, &vctypes.Proof{Type: "JWT"})

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_ID_ALREADY_REGISTERED)
}

func TestResolveID_Should_Return_Resolver_Metadata(t *testing.T) {
	t.Parallel()

	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}

	idRepo := idtesting.NewFakeIdRepository()
	sut := node.NewIdService(nil, idRepo, nil, nil)
	md := &idtypes.ResolverMetadata{
		ID: "SOME_ID",
	}
	_, _ = idRepo.CreateID(context.Background(), md, issuer)

	_, err := sut.Resolve(context.Background(), md.ID)

	assert.NoError(t, err)
}

func TestResolveID_Should_Return_Resolver_Metadata_Not_Found_Error(t *testing.T) {
	t.Parallel()

	idRepo := idtesting.NewFakeIdRepository()
	sut := node.NewIdService(nil, idRepo, nil, nil)

	_, err := sut.Resolve(context.Background(), "SOME_ID")

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_RESOLVER_METADATA_NOT_FOUND)
}

func assertErrorInfoReason(t *testing.T, err error, reason errtypes.ErrorReason) {
	t.Helper()

	var errInfo errtypes.ErrorInfo
	assert.ErrorAs(t, err, &errInfo)
	assert.Equal(t, reason, errInfo.Reason)
}
