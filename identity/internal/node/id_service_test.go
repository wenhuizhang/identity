// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package node_test

import (
	"context"
	"testing"

	errtypes "github.com/agntcy/identity/internal/core/errors/types"
	idtesting "github.com/agntcy/identity/internal/core/id/testing"
	idtypes "github.com/agntcy/identity/internal/core/id/types"
	issuertesting "github.com/agntcy/identity/internal/core/issuer/testing"
	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	coretesting "github.com/agntcy/identity/internal/core/testing"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
	"github.com/agntcy/identity/internal/node"
	"github.com/stretchr/testify/assert"
)

func TestGenerateID_Should_Not_Return_Errors(t *testing.T) {
	t.Parallel()

	verficationSrv := coretesting.NewFakeTruthyVerificationService()
	idRepo := idtesting.NewFakeIdRepository()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	sut := node.NewIdService(verficationSrv, idRepo, issuerRepo)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}
	_, _ = issuerRepo.CreateIssuer(context.Background(), issuer)

	_, err := sut.Generate(context.Background(), issuer, &vctypes.Proof{})

	assert.NoError(t, err)
}

func TestGenerateID_Should_Return_Idp_Required_Error(t *testing.T) {
	t.Parallel()

	sut := node.NewIdService(nil, nil, nil)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}

	_, err := sut.Generate(context.Background(), issuer, nil)

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_IDP_REQUIRED)
}

func TestGenerateID_Should_Return_Invalid_Proof_Error(t *testing.T) {
	t.Parallel()

	verficationSrv := coretesting.NewFalsyProofVerificationServiceStub()
	sut := node.NewIdService(verficationSrv, nil, nil)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}

	_, err := sut.Generate(context.Background(), issuer, &vctypes.Proof{})

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_INVALID_PROOF)
}

func TestGenerateID_Should_Return_Invalid_Issuer_Error(t *testing.T) {
	t.Parallel()

	verficationSrv := coretesting.NewFalsyCommonNameVerificationServiceStub()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	sut := node.NewIdService(verficationSrv, nil, issuerRepo)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}
	_, _ = issuerRepo.CreateIssuer(context.Background(), issuer)

	_, err := sut.Generate(context.Background(), issuer, &vctypes.Proof{})

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_INVALID_ISSUER)
}

func TestGenerateID_Should_Return_Unregistred_Issuer_Error(t *testing.T) {
	t.Parallel()

	verficationSrv := coretesting.NewFakeTruthyVerificationService()
	issuerRepo := issuertesting.NewFakeIssuerRepository()
	sut := node.NewIdService(verficationSrv, nil, issuerRepo)
	issuer := &issuertypes.Issuer{
		CommonName:   coretesting.ValidProofIssuer,
		Organization: "Some Org",
	}

	_, err := sut.Generate(context.Background(), issuer, &vctypes.Proof{})

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_ISSUER_NOT_REGISTERED)
}

func TestResolveID_Should_Return_Resolver_Metadata(t *testing.T) {
	t.Parallel()

	idRepo := idtesting.NewFakeIdRepository()
	sut := node.NewIdService(nil, idRepo, nil)
	md := &idtypes.ResolverMetadata{
		ID: "SOME_ID",
	}
	_, _ = idRepo.CreateID(context.Background(), md)

	_, err := sut.Resolve(context.Background(), md.ID)

	assert.NoError(t, err)
}

func TestResolveID_Should_Return_Resolver_Metadata_Not_Found_Error(t *testing.T) {
	t.Parallel()

	idRepo := idtesting.NewFakeIdRepository()
	sut := node.NewIdService(nil, idRepo, nil)

	_, err := sut.Resolve(context.Background(), "SOME_ID")

	assertErrorInfoReason(t, err, errtypes.ERROR_REASON_RESOLVER_METADATA_NOT_FOUND)
}

func assertErrorInfoReason(t *testing.T, err error, reason errtypes.ErrorReason) {
	t.Helper()

	var errInfo errtypes.ErrorInfo
	assert.ErrorAs(t, err, &errInfo)
	assert.Equal(t, reason, errInfo.Reason)
}
