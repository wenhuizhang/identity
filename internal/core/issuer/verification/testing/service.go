// Copyright 2025 AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"context"

	issuertypes "github.com/agntcy/identity/internal/core/issuer/types"
	"github.com/agntcy/identity/internal/core/issuer/verification"
	vctypes "github.com/agntcy/identity/internal/core/vc/types"
)

const (
	ValidProofIssuer string = "ISSUER"
	ValidProofSub    string = "SUBJECT"
)

type FakeVerifiedVerificationServiceStub struct{}

func NewFakeVerifiedVerificationServiceStub() verification.Service {
	return &FakeVerifiedVerificationServiceStub{}
}

func (f *FakeVerifiedVerificationServiceStub) Verify(
	ctx context.Context,
	issuer *issuertypes.Issuer,
	proof *vctypes.Proof,
) (*verification.Result, error) {
	return &verification.Result{
		Issuer:   issuer,
		Verified: true,
	}, nil
}

func (f *FakeVerifiedVerificationServiceStub) VerifyExistingIssuer(
	ctx context.Context,
	proof *vctypes.Proof,
) (*verification.Result, error) {
	panic("unimplemented")
}

type FakeUnverifiedVerificationServiceStub struct{}

func NewFakeUnverifiedVerificationServiceStub() verification.Service {
	return &FakeUnverifiedVerificationServiceStub{}
}

func (f *FakeUnverifiedVerificationServiceStub) Verify(
	ctx context.Context,
	issuer *issuertypes.Issuer,
	proof *vctypes.Proof,
) (*verification.Result, error) {
	return &verification.Result{
		Issuer:   issuer,
		Verified: false,
	}, nil
}

func (f *FakeUnverifiedVerificationServiceStub) VerifyExistingIssuer(
	ctx context.Context,
	proof *vctypes.Proof,
) (*verification.Result, error) {
	panic("unimplemented")
}
